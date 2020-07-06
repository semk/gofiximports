# gofiximports
Utility to fix/replace import paths in Go files recursively & pretty-print like gofmt.

## Installation

```
$ go get github.com/semk/gofiximports
```

## Usage

```
Usage of gofiximports:
  -dir string
    	directory where the replacements are to be done (default "./")
  -from string
    	import statement to be replaced
  -indent int
    	all code is indented at least by this much
  -rawformat
    	do not use a tabwriter; if set, -usespaces is ignored
  -recursive
    	peform the replacecements recursively on the directory (to be used with -dir) (default true)
  -sourcepos
    	emit //line directives to preserve original source positions
  -stdin
    	read the file names from stdin where the replacements are to be done (overrides -dir)
  -tabindent
    	use tabs for indentation independent of -usespaces (default true)
  -tabwidth int
    	tab width (default 8)
  -to string
    	replacement import statement
  -usespaces
    	use spaces instead of tabs for alignment (default true)
```

## Examples

The following example replaces all imports of `"library/module"` to `"repository/library/module"`
recursively inside the `awesome_go_project` directory. The command only modifies `.go` files.

```
$ gofiximports -dir awesome_go_project -from "library/module" -to "repository/library/module"
```

NOTE: The above command will replace all the imports starting with `"library/module"`.
`"library/module/x"` will be changed to `"repository/library/module/x"` as well.

Before fixing imports:
```go
package main
import logger "fmt"
func main() {
    logger.Println("hello world")
}
```

```
$ gofiximports -dir test -from "fmt" -to "log"
2020/07/06 15:07:51 Modified import paths for 1 files in path "test"
```

After fixing imports:
```go
package main

import logger "log"

func main() {
	logger.Println("hello world")
}

```

You can also pass the list of files to `stdin` using pipe.

```
$ ls test/test.go | gofiximports -stdin -from "fmt" -to "log" 
```
