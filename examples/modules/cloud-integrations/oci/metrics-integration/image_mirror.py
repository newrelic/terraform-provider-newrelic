#!/usr/bin/env python3
"""
Copies the New Relic function image from a public registry (Docker Hub by default) into a
private OCIR repository in the customer's tenancy, which the function then runs from.

Resource Manager jobs have no Docker daemon, so this talks to both registries over the
Docker Registry HTTP API v2 and uses only the Python standard library.

Modes:
  resolve  Terraform external data source. Reads {"source_image", "platform"} from stdin and
           prints {"digest", "tag"}: the digest of the platform-specific manifest and the tag
           to push it under. Read-only.
  copy     Called from a local-exec provisioner. Copies SOURCE_DIGEST (with its config and
           layers) to DEST_REGISTRY/DEST_REPOSITORY:DEST_TAG. Settings come from environment
           variables so the registry password never appears on the command line.

The manifest bytes are pushed unchanged, so the digest in OCIR equals the source digest and
Terraform knows it at plan time.
"""
import base64
import hashlib
import json
import os
import re
import socket
import sys
import tempfile
import time
import urllib.error
import urllib.parse
import urllib.request

OCI_INDEX = "application/vnd.oci.image.index.v1+json"
OCI_MANIFEST = "application/vnd.oci.image.manifest.v1+json"
DOCKER_LIST = "application/vnd.docker.distribution.manifest.list.v2+json"
DOCKER_MANIFEST = "application/vnd.docker.distribution.manifest.v2+json"
MANIFEST_ACCEPT = ",".join([OCI_INDEX, DOCKER_LIST, OCI_MANIFEST, DOCKER_MANIFEST])

MAX_ATTEMPTS = 6
TIMEOUT_SECONDS = 120
CHUNK_SIZE = 1024 * 1024
# A newly created OCI auth token can take a few minutes before the registry accepts it.
AUTH_PROPAGATION_SECONDS = int(os.environ.get("AUTH_PROPAGATION_SECONDS", "300"))


class RegistryError(Exception):
    pass


def log(message):
    print("[image_mirror] " + message, file=sys.stderr, flush=True)


def parse_image(ref):
    """Splits an image reference into (registry host, repository, tag, digest)."""
    digest = None
    if "@" in ref:
        ref, digest = ref.split("@", 1)
    first, _, rest = ref.partition("/")
    if rest and ("." in first or ":" in first or first == "localhost"):
        registry, path = first, rest
    else:
        registry, path = "docker.io", ref
    tag = None
    if ":" in path.rsplit("/", 1)[-1]:
        path, tag = path.rsplit(":", 1)
    if registry == "docker.io":
        registry = "registry-1.docker.io"
        if "/" not in path:
            path = "library/" + path
    if not tag and not digest:
        tag = "latest"
    return registry, path, tag, digest


class _NoRedirect(urllib.request.HTTPRedirectHandler):
    # Redirects are followed by hand so the registry's Authorization header is not
    # forwarded to the blob storage host (S3/CDN reject requests with two auth schemes).
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


_opener = urllib.request.build_opener(_NoRedirect)


def _open(method, url, headers, body=None):
    """Sends one request, following redirects. body is bytes, a file path, or None."""
    for _ in range(5):
        upload = open(body, "rb") if isinstance(body, str) else None
        try:
            request = urllib.request.Request(url, data=upload or body, headers=headers, method=method)
            return _opener.open(request, timeout=TIMEOUT_SECONDS)
        except urllib.error.HTTPError as e:
            if e.code not in (301, 302, 303, 307, 308):
                raise
            location = urllib.parse.urljoin(url, e.headers["Location"])
            if urllib.parse.urlparse(location).netloc != urllib.parse.urlparse(url).netloc:
                headers = {k: v for k, v in headers.items() if k.lower() != "authorization"}
            url = location
            if e.code == 303:
                method, body = "GET", None
        finally:
            if upload:
                upload.close()
    raise RegistryError("too many redirects for " + url)


def _backoff(attempt, retry_after=None):
    if retry_after and retry_after.isdigit():
        delay = min(int(retry_after), 120)
    else:
        delay = min(2 ** attempt * 2, 60)
    time.sleep(delay)


class Registry:
    def __init__(self, host, repository, actions, username="", password=""):
        self.host = host
        self.repository = repository
        self.actions = actions
        self.username = username
        self.password = password
        self.authorization = None

    def url(self, path):
        if path.startswith("https://") or path.startswith("http://"):
            return path
        scheme = "http" if self.host.startswith("localhost") or self.host.startswith("127.0.0.1") else "https"
        return "%s://%s/v2/%s%s" % (scheme, self.host, self.repository, path)

    def request(self, method, path, headers=None, body=None, ok=(200,), retry=True):
        """Returns the open response. Handles auth challenges, rate limits and 5xx."""
        url = self.url(path)
        refreshed = False
        attempt = 0
        while True:
            all_headers = dict(headers or {})
            if self.authorization:
                all_headers["Authorization"] = self.authorization
            try:
                response = _open(method, url, all_headers, body)
                if response.getcode() in ok:
                    return response
                raise RegistryError("%s %s returned %s" % (method, url, response.getcode()))
            except urllib.error.HTTPError as e:
                if e.code in ok:
                    return e
                if e.code == 401 and not refreshed:
                    self._authenticate(e.headers.get("WWW-Authenticate", ""))
                    refreshed = True
                    continue
                if retry and (e.code == 429 or e.code >= 500) and attempt < MAX_ATTEMPTS:
                    attempt += 1
                    log("%s %s returned %s, retrying (%d/%d)" % (method, url, e.code, attempt, MAX_ATTEMPTS))
                    _backoff(attempt, e.headers.get("Retry-After"))
                    continue
                detail = e.read()[:500].decode("utf-8", "replace")
                raise RegistryError("%s %s returned %s: %s" % (method, url, e.code, detail))
            except (urllib.error.URLError, ConnectionError, socket.timeout) as e:
                if retry and attempt < MAX_ATTEMPTS:
                    attempt += 1
                    log("%s %s failed (%s), retrying (%d/%d)" % (method, url, e, attempt, MAX_ATTEMPTS))
                    _backoff(attempt)
                    continue
                raise RegistryError("%s %s failed: %s" % (method, url, e))

    def _authenticate(self, challenge):
        basic = None
        if self.username:
            basic = "Basic " + base64.b64encode(
                ("%s:%s" % (self.username, self.password)).encode()).decode()
        if challenge.lower().startswith("basic"):
            if not basic:
                raise RegistryError("%s requires credentials" % self.host)
            self.authorization = basic
            return
        params = dict(re.findall(r'(\w+)="([^"]*)"', challenge))
        if "realm" not in params:
            raise RegistryError("unsupported auth challenge from %s: %r" % (self.host, challenge))
        query = {"scope": "repository:%s:%s" % (self.repository, self.actions)}
        if params.get("service"):
            query["service"] = params["service"]
        token_url = params["realm"] + "?" + urllib.parse.urlencode(query)
        headers = {"Authorization": basic} if basic else {}

        deadline = time.time() + (AUTH_PROPAGATION_SECONDS if basic else 0)
        attempt = 0
        while True:
            try:
                with _open("GET", token_url, headers) as response:
                    payload = json.loads(response.read())
                token = payload.get("token") or payload.get("access_token")
                if not token:
                    raise RegistryError("no token in auth response from %s" % self.host)
                self.authorization = "Bearer " + token
                return
            except urllib.error.HTTPError as e:
                if e.code in (401, 403) and time.time() < deadline:
                    log("%s rejected the credentials for %s; a new auth token can take a few "
                        "minutes to become active, retrying in 15s" % (self.host, self.username))
                    time.sleep(15)
                    continue
                if (e.code == 429 or e.code >= 500) and attempt < MAX_ATTEMPTS:
                    attempt += 1
                    _backoff(attempt, e.headers.get("Retry-After"))
                    continue
                raise RegistryError("authentication to %s failed with %s: %s" % (
                    self.host, e.code, e.read()[:300].decode("utf-8", "replace")))

    def get_manifest(self, reference):
        with self.request("GET", "/manifests/" + reference, headers={"Accept": MANIFEST_ACCEPT}) as response:
            body = response.read()
            media_type = (response.headers.get("Content-Type") or "").split(";")[0].strip()
        if media_type not in (OCI_INDEX, DOCKER_LIST, OCI_MANIFEST, DOCKER_MANIFEST):
            media_type = json.loads(body).get("mediaType", media_type)
        return media_type, body

    def manifest_digest(self, reference):
        """Digest the registry has for reference, or None if it does not exist."""
        with self.request("HEAD", "/manifests/" + reference, headers={"Accept": MANIFEST_ACCEPT},
                          ok=(200, 404)) as response:
            if response.getcode() == 404:
                return None
            return response.headers.get("Docker-Content-Digest")

    def blob_exists(self, digest):
        with self.request("HEAD", "/blobs/" + digest, ok=(200, 404)) as response:
            return response.getcode() == 200

    def download_blob(self, digest, path):
        for attempt in range(1, MAX_ATTEMPTS + 1):
            sha = hashlib.sha256()
            with self.request("GET", "/blobs/" + digest) as response, open(path, "wb") as out:
                while True:
                    chunk = response.read(CHUNK_SIZE)
                    if not chunk:
                        break
                    sha.update(chunk)
                    out.write(chunk)
            if "sha256:" + sha.hexdigest() == digest:
                return
            log("blob %s failed digest verification, downloading again (%d/%d)" % (digest, attempt, MAX_ATTEMPTS))
        raise RegistryError("blob %s does not match its digest" % digest)

    def upload_blob(self, digest, path):
        size = os.path.getsize(path)
        for attempt in range(1, MAX_ATTEMPTS + 1):
            try:
                with self.request("POST", "/blobs/uploads/", headers={"Content-Length": "0"},
                                  body=b"", ok=(202,)) as response:
                    location = self._location(response)
                # One PATCH with the whole blob, then commit: the sequence `docker push` uses.
                with self.request("PATCH", location, headers={
                        "Content-Type": "application/octet-stream", "Content-Length": str(size)},
                        body=path, ok=(202,), retry=False) as response:
                    location = self._location(response)
                separator = "&" if "?" in location else "?"
                commit_url = location + separator + "digest=" + urllib.parse.quote(digest)
                with self.request("PUT", commit_url, headers={"Content-Length": "0"},
                                  body=b"", ok=(201,), retry=False):
                    return
            except RegistryError as e:
                if attempt == MAX_ATTEMPTS:
                    raise
                log("upload of %s failed (%s), starting a new upload (%d/%d)" % (digest, e, attempt, MAX_ATTEMPTS))
                _backoff(attempt)

    def put_manifest(self, tag, media_type, body):
        with self.request("PUT", "/manifests/" + tag, headers={"Content-Type": media_type},
                          body=body, ok=(200, 201)) as response:
            return response.headers.get("Docker-Content-Digest")

    def _location(self, response):
        location = response.headers.get("Location")
        if not location:
            raise RegistryError("%s did not return an upload location" % self.host)
        return urllib.parse.urljoin(self.url("/"), location)


def sha256_digest(data):
    return "sha256:" + hashlib.sha256(data).hexdigest()


def resolve(source_image, platform):
    """Returns (digest, tag) of the manifest for platform behind source_image."""
    host, repository, tag, digest = parse_image(source_image)
    source = Registry(host, repository, "pull",
                      os.environ.get("SOURCE_USERNAME", ""), os.environ.get("SOURCE_PASSWORD", ""))
    media_type, body = source.get_manifest(digest or tag)

    if media_type in (OCI_INDEX, DOCKER_LIST):
        wanted_os, wanted_arch = platform.split("/")[:2]
        matches = [m for m in json.loads(body).get("manifests", [])
                   if m.get("platform", {}).get("os") == wanted_os
                   and m.get("platform", {}).get("architecture") == wanted_arch]
        if not matches:
            raise RegistryError("%s has no %s image" % (source_image, platform))
        media_type, body = source.get_manifest(matches[0]["digest"])

    manifest_digest = sha256_digest(body)
    # A digest-pinned source has no tag; derive a stable one from the digest.
    dest_tag = tag or "sha256-" + manifest_digest.split(":", 1)[1][:12]
    return manifest_digest, dest_tag


def copy(env):
    source_image = env["SOURCE_IMAGE"]
    digest = env["SOURCE_DIGEST"]
    platform = env.get("PLATFORM", "linux/amd64")
    host, repository, _, _ = parse_image(source_image)
    source = Registry(host, repository, "pull",
                      env.get("SOURCE_USERNAME", ""), env.get("SOURCE_PASSWORD", ""))
    dest = Registry(env["DEST_REGISTRY"], env["DEST_REPOSITORY"], "pull,push",
                    env.get("DEST_USERNAME", ""), env.get("DEST_PASSWORD", ""))
    dest_tag = env["DEST_TAG"]
    dest_ref = "%s/%s:%s" % (dest.host, dest.repository, dest_tag)

    if dest.manifest_digest(dest_tag) == digest:
        log("%s already points at %s, nothing to copy" % (dest_ref, digest))
        print(digest)
        return

    log("copying %s (%s) to %s" % (source_image, digest, dest_ref))
    media_type, body = source.get_manifest(digest)
    if sha256_digest(body) != digest:
        raise RegistryError("manifest from %s does not match %s" % (source_image, digest))
    if media_type in (OCI_INDEX, DOCKER_LIST):
        raise RegistryError("%s is a multi-platform index; expected a single %s manifest" % (digest, platform))
    manifest = json.loads(body)
    blobs = [manifest["config"]] + manifest.get("layers", [])

    with tempfile.TemporaryDirectory() as workdir:
        for blob in blobs:
            blob_digest = blob["digest"]
            if dest.blob_exists(blob_digest):
                log("blob %s already present" % blob_digest)
                continue
            path = os.path.join(workdir, blob_digest.replace(":", "_"))
            log("copying blob %s (%.1f MB)" % (blob_digest, blob.get("size", 0) / 1e6))
            source.download_blob(blob_digest, path)
            if blob is manifest["config"]:
                with open(path, "rb") as f:
                    config = json.load(f)
                image_platform = "%s/%s" % (config.get("os"), config.get("architecture"))
                if image_platform != platform:
                    raise RegistryError("%s is built for %s, the function needs %s" % (
                        source_image, image_platform, platform))
            dest.upload_blob(blob_digest, path)
            os.remove(path)

    pushed_digest = dest.put_manifest(dest_tag, media_type, body)
    if pushed_digest and pushed_digest != digest:
        raise RegistryError("%s reported digest %s, expected %s" % (dest_ref, pushed_digest, digest))
    log("pushed %s@%s" % (dest_ref, digest))
    print(digest)


def main():
    mode = sys.argv[1] if len(sys.argv) > 1 else ""
    try:
        if mode == "resolve":
            query = json.load(sys.stdin)
            digest, tag = resolve(query["source_image"], query.get("platform") or "linux/amd64")
            print(json.dumps({"digest": digest, "tag": tag}))
        elif mode == "copy":
            copy(os.environ)
        else:
            print("usage: image_mirror.py resolve|copy", file=sys.stderr)
            sys.exit(2)
    except (RegistryError, KeyError, ValueError) as e:
        log("error: %s" % e)
        sys.exit(1)


if __name__ == "__main__":
    main()
