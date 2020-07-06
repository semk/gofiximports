/******************************************************************************
* This utility is used to replace the import paths of .go files in a directory
* recursively.
*
*	gofiximports -dir <go-files-dir> -from <old-import> -to <new-import>
*
*	Copyright (c) 2020 Sreejith Kesavan <sreejithemk@gmail.com>
*
******************************************************************************/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"go/ast"
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

func replaceImportsInFileAST(fset *token.FileSet, file *ast.File, filePath, from, to string, printMode printer.Mode, tabWidth, indent int) (bool, error) {
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
		// Use the provided formatting and rewrite the formatted file
		printConfig := &printer.Config{Mode: printMode, Tabwidth: tabWidth, Indent: indent}
		var outputBuffer bytes.Buffer

		err := printConfig.Fprint(&outputBuffer, fset, file)
		if err != nil {
			return false, err
		}

		err = ioutil.WriteFile(filePath, outputBuffer.Bytes(), os.ModePerm)
		if err != nil {
			return false, err
		}
	}
	return rewrote, nil
}

func replaceImportsInFile(filePath, from, to string, printMode printer.Mode, tabWidth, indent int) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)

	rewrote, err := replaceImportsInFileAST(fset, file, filePath, from, to, printMode, tabWidth, indent)
	if err != nil {
		return err
	}

	if rewrote {
		log.Printf("Modified import paths in file \"%s\"\n", filePath)
	}

	return nil
}

func replaceImportsInFilesFromStdin(from, to string, printMode printer.Mode, tabWidth, indent int) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := replaceImportsInFile(scanner.Text(), from, to, printMode, tabWidth, indent)
		if err != nil {
			return err
		}
	}

	return nil
}

func replaceImportsInDir(dir, from, to string, printMode printer.Mode, tabWidth, indent int) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	modifiedFilesCount := 0
	for _, pkg := range pkgs {
		for filePath, file := range pkg.Files {
			rewrote, err := replaceImportsInFileAST(fset, file, filePath, from, to, printMode, tabWidth, indent)
			if err != nil {
				return err
			}

			if rewrote {
				modifiedFilesCount++
			}
		}
	}

	if modifiedFilesCount > 0 {
		log.Printf("Modified import paths for %d files in path \"%s\"\n", modifiedFilesCount, dir)
	}

	return nil
}

func replaceImportsInDirRecursive(dir, from, to string, printMode printer.Mode, tabWidth, indent int) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				err = replaceImportsInDir(path, from, to, printMode, tabWidth, indent)
			}
			return err
		})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Common options
	dir := flag.String("dir", "./", "directory where the replacements are to be done")
	recursive := flag.Bool("recursive", true, "peform the replacecements recursively on the directory (to be used with -dir)")
	stdin := flag.Bool("stdin", false, "read the file names from stdin where the replacements are to be done (overrides -dir)")
	from := flag.String("from", "", "import statement to be replaced")
	to := flag.String("to", "", "replacement import statement")

	// Code formatting options
	rawFormat := flag.Bool("rawformat", false, "do not use a tabwriter; if set, -usespaces is ignored")
	tabIndent := flag.Bool("tabindent", true, "use tabs for indentation independent of -usespaces")
	useSpaces := flag.Bool("usespaces", true, "use spaces instead of tabs for alignment")
	sourcePos := flag.Bool("sourcepos", false, "emit //line directives to preserve original source positions")
	tabWidth := flag.Int("tabwidth", 8, "tab width")
	indent := flag.Int("indent", 0, "all code is indented at least by this much")

	// Parse commandline arguments
	flag.Parse()

	// Source code printer options
	var printMode printer.Mode
	if *rawFormat {
		printMode = printMode | printer.RawFormat
	}
	if *tabIndent {
		printMode = printMode | printer.TabIndent
	}
	if *useSpaces {
		printMode = printMode | printer.UseSpaces
	}
	if *sourcePos {
		printMode = printMode | printer.SourcePos
	}

	var err error
	if *stdin {
		err = replaceImportsInFilesFromStdin(*from, *to, printMode, *tabWidth, *indent)
	} else if *recursive {
		err = replaceImportsInDirRecursive(*dir, *from, *to, printMode, *tabWidth, *indent)
	} else {
		err = replaceImportsInDir(*dir, *from, *to, printMode, *tabWidth, *indent)
	}

	if err != nil {
		log.Fatal(err)
	}
}
