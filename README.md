# go-template

Empty project layout for making Golang apps following the [project-layout](https://github.com/golang-standards/project-layout) standard

## Auto-versioning

The `Makefile` will automatically pull the version from the latest `git tag` and pass that through to the linker.  To use this feature, do the following:

### Add a 'Version' to your main package

```
package main

import "fmt"

var (
  // Version is your app version (updated by Makefile, don't forget to TAG YOUR RELEASE)
  Version = "undefined"
)

func main() {
  fmt.Printf("Example App version: %s\n", Version)
}
```

### Create a tag before you build your release

For example, to make a version 0.0.1:

```
git tag v0.0.1
```

### Example Version Strings

```
# No Tags (latest sha):
Example App version: g1de6b99

# Clean tag:
Example App version: v0.0.3

# Latest tag: v0.0.3
# One commit has passed since that tag
# SHA of current commit
Example App version: v0.0.3-1-g1de6b99

# Local changes, uncommited
Example App version: v0.0.3-1-g1de6b99-dirty
```

### Notes

* If you have NO commits, make will fail... Solve this with an initial commit in the repo: `git commit -m 'Initial commit'`
* If you do not create a tag, you will get the sha as the version
* If you have uncommitted changes, your version will end with `-dirty` (i.e. `v1.2.3-dirty`)
