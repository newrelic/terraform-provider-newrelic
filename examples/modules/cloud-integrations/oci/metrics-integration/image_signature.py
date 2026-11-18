"""
Checks that a function image was signed by New Relic before image_mirror.py copies it.

New Relic signs every image it publishes with its OCI signing service. The signature is a cosign
"simple signing" signature stored next to the image under the tag sha256-<hex>.sig: a manifest
whose layers each hold a payload naming the image digest, with an Ed25519 signature over that
payload in the layer annotation. This verifies it against the public keys New Relic publishes for
that signing key (a JWKS at publickeys.newrelic.com), read when the stack runs so a rotated key is
picked up without a new stack release. Public Sigstore services are not used.

Standard library only, like image_mirror.py: Resource Manager jobs cannot install packages, so
Ed25519 verification is implemented here from RFC 8032.
"""
import base64
import hashlib
import json
import os
import socket
import tempfile
import time
import urllib.error
import urllib.request

SIMPLE_SIGNING = "application/vnd.dev.cosign.simplesigning.v1+json"
SIGNATURE_ANNOTATION = "dev.cosignproject.cosign/signature"
SIGNATURE_TYPE = "cosign container image signature"
# Signature payloads and key sets are a few hundred bytes; anything larger is not one.
MAX_PAYLOAD_BYTES = 64 * 1024
MAX_JWKS_BYTES = 64 * 1024
KEYS_ATTEMPTS = 3
KEYS_TIMEOUT_SECONDS = 30


class SignatureError(Exception):
    pass


# --- Ed25519 verification (RFC 8032 section 5.1.7) -------------------------------------------
# Only verification: it handles public data, so constant-time arithmetic is not needed.
# Inverses use Fermat's little theorem (pow(x, P - 2, P)) so this runs on any Python 3.

_P = 2 ** 255 - 19
_L = 2 ** 252 + 27742317777372353535851937790883648493
_D = -121665 * pow(121666, _P - 2, _P) % _P
_SQRT_M1 = pow(2, (_P - 1) // 4, _P)


def _add(a, b):
    """Adds two points in extended coordinates (x, y, z, t)."""
    x1, y1, z1, t1 = a
    x2, y2, z2, t2 = b
    a = (y1 - x1) * (y2 - x2) % _P
    b = (y1 + x1) * (y2 + x2) % _P
    c = 2 * t1 * t2 * _D % _P
    d = 2 * z1 * z2 % _P
    e, f, g, h = b - a, d - c, d + c, b + a
    return (e * f % _P, g * h % _P, f * g % _P, e * h % _P)


def _multiply(scalar, point):
    result = (0, 1, 1, 0)
    while scalar:
        if scalar & 1:
            result = _add(result, point)
        point = _add(point, point)
        scalar >>= 1
    return result


def _equal(a, b):
    return ((a[0] * b[2] - b[0] * a[2]) % _P == 0
            and (a[1] * b[2] - b[1] * a[2]) % _P == 0)


def _decode_point(data):
    """Decodes a 32-byte point, or returns None if it is not a canonical encoding of one."""
    if len(data) != 32:
        return None
    y = int.from_bytes(data, "little")
    sign = y >> 255
    y &= (1 << 255) - 1
    if y >= _P:
        return None
    x2 = (y * y - 1) * pow(_D * y * y + 1, _P - 2, _P) % _P
    x = pow(x2, (_P + 3) // 8, _P)
    if (x * x - x2) % _P:
        x = x * _SQRT_M1 % _P
    if (x * x - x2) % _P:
        return None
    if x == 0 and sign:
        return None
    if x & 1 != sign:
        x = _P - x
    return (x, y, 1, x * y % _P)


_BASE = _decode_point((4 * pow(5, _P - 2, _P) % _P).to_bytes(32, "little"))


def ed25519_verify(public_key, message, signature):
    """True if signature is a valid Ed25519 signature of message under public_key (32 bytes)."""
    if len(signature) != 64:
        return False
    a = _decode_point(public_key)
    r = _decode_point(signature[:32])
    s = int.from_bytes(signature[32:], "little")
    # s >= L would make the signature malleable: reject it as RFC 8032 requires.
    if a is None or r is None or s >= _L:
        return False
    k = int.from_bytes(hashlib.sha512(signature[:32] + public_key + message).digest(), "little") % _L
    return _equal(_multiply(s, _BASE), _add(r, _multiply(k, a)))


# --- Keys ---------------------------------------------------------------------------------------

def _b64url(value):
    return base64.urlsafe_b64decode(value + "=" * (-len(value) % 4))


def parse_keys(data, source):
    """Returns [(kid, public key bytes)] for the Ed25519 keys in a JWKS document."""
    if len(data) > MAX_JWKS_BYTES:
        raise SignatureError("%s is larger than a key set" % source)
    try:
        jwks = json.loads(data)
    except ValueError:
        raise SignatureError("%s is not a JSON key set" % source)
    keys = []
    for key in jwks.get("keys", []):
        if key.get("kty") != "OKP" or key.get("crv") != "Ed25519" or not key.get("x"):
            continue
        x = _b64url(key["x"])
        if len(x) == 32:
            keys.append((key.get("kid") or "?", x))
    if not keys:
        raise SignatureError("%s has no New Relic signing keys" % source)
    return keys


def fetch_keys(url):
    """Reads the JWKS New Relic publishes at url and returns its Ed25519 keys."""
    if not url.startswith("https://"):
        raise SignatureError("signing keys must come over HTTPS, not %s" % url)
    request = urllib.request.Request(url, headers={
        "Accept": "application/json", "User-Agent": "newrelic-oci-image-mirror"})
    for attempt in range(1, KEYS_ATTEMPTS + 1):
        try:
            with urllib.request.urlopen(request, timeout=KEYS_TIMEOUT_SECONDS) as response:
                return parse_keys(response.read(MAX_JWKS_BYTES + 1), url)
        except urllib.error.HTTPError as e:
            if e.code < 500 or attempt == KEYS_ATTEMPTS:
                raise SignatureError("%s returned HTTP %s" % (url, e.code))
        except (urllib.error.URLError, ConnectionError, socket.timeout) as e:
            if attempt == KEYS_ATTEMPTS:
                raise SignatureError("could not reach %s: %s" % (url, getattr(e, "reason", e)))
        time.sleep(2 ** attempt)


# --- cosign signatures ----------------------------------------------------------------------

def _repository(reference):
    """Normalizes an image reference to registry/repository, without tag or digest."""
    reference = reference.split("@", 1)[0]
    name, _, last = reference.rpartition("/")
    if ":" in last:
        last = last.rsplit(":", 1)[0]
    reference = name + "/" + last if name else last
    first, _, rest = reference.partition("/")
    if not rest or not ("." in first or ":" in first or first == "localhost"):
        first, rest = "docker.io", reference
    if first in ("index.docker.io", "registry-1.docker.io"):
        first = "docker.io"
    if first == "docker.io" and "/" not in rest:
        rest = "library/" + rest
    return first + "/" + rest


def _payload(source, layer, workdir):
    """The layer's payload bytes, checked against the layer digest."""
    digest = layer.get("digest", "")
    if layer.get("size", 0) > MAX_PAYLOAD_BYTES:
        raise SignatureError("signature payload %s is %d bytes" % (digest, layer["size"]))
    if layer.get("data"):
        data = base64.b64decode(layer["data"])
        if "sha256:" + hashlib.sha256(data).hexdigest() != digest:
            raise SignatureError("inline signature payload does not match %s" % digest)
        return data
    path = os.path.join(workdir, digest.replace(":", "_"))
    source.download_blob(digest, path)
    with open(path, "rb") as f:
        return f.read()


def _check_layer(source, layer, digest, identity, keys, workdir):
    """Returns the kid of the key that signed this layer for digest; raises SignatureError if not."""
    if layer.get("mediaType") != SIMPLE_SIGNING:
        raise SignatureError("not a cosign signature layer")
    encoded = (layer.get("annotations") or {}).get(SIGNATURE_ANNOTATION)
    if not encoded:
        raise SignatureError("no signature annotation")
    signature = base64.b64decode(encoded)
    payload = _payload(source, layer, workdir)
    # The payload is only read once the signature over it checks out.
    kid = next((kid for kid, key in keys if ed25519_verify(key, payload, signature)), None)
    if kid is None:
        raise SignatureError("signature does not verify against any New Relic key")
    critical = json.loads(payload).get("critical") or {}
    if critical.get("type") != SIGNATURE_TYPE:
        raise SignatureError("signed payload has type %r" % critical.get("type"))
    signed_digest = (critical.get("image") or {}).get("docker-manifest-digest")
    if signed_digest != digest:
        raise SignatureError("signature is for %s" % signed_digest)
    # Bound to the repository New Relic published to, not the tag: the payload names the version
    # tag, but stable and latest point at the same digest.
    signed_identity = (critical.get("identity") or {}).get("docker-reference", "")
    if _repository(signed_identity) != _repository(identity):
        raise SignatureError("signature is for %s" % signed_identity)
    return kid


def verify(source, digest, identity, keys):
    """
    Checks that the image `digest` in registry `source` (an image_mirror.Registry) carries a
    signature by one of `keys` for `identity`, the repository New Relic publishes the image to.
    Returns the kid of the signing key; raises SignatureError otherwise.
    """
    if not digest.startswith("sha256:"):
        raise SignatureError("cannot look up a signature for %s" % digest)
    signature_tag = digest.replace(":", "-") + ".sig"
    if source.manifest_digest(signature_tag) is None:
        raise SignatureError("no signature found (%s:%s)" % (source.repository, signature_tag))
    _, body = source.get_manifest(signature_tag)
    reasons = []
    with tempfile.TemporaryDirectory() as workdir:
        # Signing the same digest again adds a layer, so any one valid layer is enough.
        for layer in json.loads(body).get("layers", []):
            try:
                return _check_layer(source, layer, digest, identity, keys, workdir)
            except (SignatureError, ValueError, KeyError) as e:
                reasons.append(str(e))
    raise SignatureError("; ".join(sorted(set(reasons))) or "signature manifest has no layers")
