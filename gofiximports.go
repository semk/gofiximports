/******************************************************************************
* This utility is used to replace the import paths of .go files in a directory
* recursively.
*
*	gofiximports -dir <go-files-dir> -from <old-import> -to <new-import>
*
******************************************************************************/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func replaceImportsInDir(dir, from, to string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	modifiedFilesCount := 0
	for _, pkg := range pkgs {
		for filePath, file := range pkg.Files {
			reWriteSpec := make(map[string]string)
			for _, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, "\"")
				if strings.HasPrefix(importPath, from) {
					reWriteSpec[importPath] = to + importPath[len(from):]
				}
			}

			var rewrote bool
			for newFrom, newTo := range reWriteSpec {
				mod := astutil.RewriteImport(fset, file, newFrom, newTo)
				rewrote = rewrote || mod
			}

			if rewrote {
				// Mimic gofmt and rewrite the formatted file
				printMode := printer.TabIndent | printer.UseSpaces
				printConfig := &printer.Config{Mode: printMode, Tabwidth: 4}
				var outputBuffer bytes.Buffer

				err := printConfig.Fprint(&outputBuffer, fset, file)
				if err != nil {
					return err
				}

				err = ioutil.WriteFile(filePath, outputBuffer.Bytes(), os.ModePerm)
				if err != nil {
					return err
				}
				modifiedFilesCount++
			}
		}
	}

	if modifiedFilesCount > 0 {
		fmt.Printf("Modified import paths for %d files in path \"%s\"\n", modifiedFilesCount, dir)
	}

	return nil

}

func replaceImportsInDirRecursive(dir, from, to string) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				err = replaceImportsInDir(path, from, to)
			}
			return err
		})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	dir := flag.String("dir", "./", "Directory path where the replacements are to be done")
	from := flag.String("from", "", "Import path to be matched")
	to := flag.String("to", "", "Import path to be changed")
	flag.Parse()

	err := replaceImportsInDirRecursive(*dir, *from, *to)
	if err != nil {
		log.Fatal(err)
	}
}
