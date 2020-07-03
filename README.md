# gofiximports
Utility to fix/replace import paths in Go files recursively & format like gofmt

## Examples

The following example replaces all imports of `"library/module"` to `"repository/library/module"`
recursively inside the `awesome_go_project` directory. The command only modifies `.go` files.

```
$ gofiximports -dir awesome_go_project -from "library/module" -to "repository/library/module"
```

NOTE: The above command will replace all the imports starting with `"library/module"`.
`"library/module/x"` will be changed to `"repository/library/module/x"` as well.
