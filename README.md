# gofiximports
Utility to replace import paths in Go files recursively with formatting similar to gofmt

## Examples

The following example replaces all imports of `"library/module"` to `"repository/library/module"`
recursively inside the `awesome_go_project` directory. The command only modifies `.go` files.

```
gofiximports -dir awesome_go_project -from "library/module" -to "repository/library/module"
```
