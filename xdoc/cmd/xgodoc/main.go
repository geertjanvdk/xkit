// Copyright (c) 2021, Geert JM Vanderkelen

package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/geertjanvdk/xkit/xansi"
	"github.com/geertjanvdk/xkit/xdoc"
	"github.com/geertjanvdk/xkit/xpath"
)

var (
	flagPackageDir string
	flagOut        string
	flagForceWrite bool
	flagNoHR       bool
)

var writeTo io.Writer

func write(a ...interface{}) {
	_, _ = fmt.Fprint(writeTo, a...)
}

func writef(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(writeTo, format, a...)
}

func writeln(a ...interface{}) {
	_, _ = fmt.Fprintln(writeTo, a...)
}

func printErrorf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, xansi.Render{xansi.Red}.Sprintf("Error: "+format+"\n", a...)+xansi.Reset())
}

func printSuccessf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, xansi.Render{xansi.Blue}.Sprintf(format+"\n", a...)+xansi.Reset())
}

func printDebugf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, xansi.Render{xansi.BrightBlack}.Sprintf(format+"\n", a...)+xansi.Reset())
}

var _ = printDebugf // fool linters: function is used

func main() {
	flag.StringVar(&flagPackageDir, "package-dir", ".",
		"Path of package to generate docs for (can be relative)")
	flag.StringVar(&flagOut, "output", "",
		"File where output is stored; use '-' for STDOUT (default <packageName>.md)")
	flag.BoolVar(&flagForceWrite, "force", false,
		"Whether we overwrite existing files or not")
	flag.BoolVar(&flagNoHR, "no-hr", false,
		"Whether to show horizontal rules or not")

	flag.Parse()

	doc := xdoc.New(flagPackageDir)

	if flagOut == "" {
		flagOut = doc.PackageName + ".md"
	}

	if flagOut != "-" {
		if !path.IsAbs(flagOut) {
			wd, err := os.Getwd()
			if err != nil {
				printErrorf("failed getting current working directory (%s)", err)
				os.Exit(1)
			}
			flagOut = path.Join(wd, flagOut)
		}

		if !flagForceWrite && xpath.IsRegularFile(flagOut) {
			printErrorf("file %s exists (use -force to overwrite)", flagOut)
			os.Exit(1)
		}

		var err error
		writeTo, err = os.OpenFile(flagOut, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			printErrorf("output file cannot be used (%s)", err)
			os.Exit(1)
		}
		defer func() {
			if err := writeTo.(io.Closer).Close(); err != nil {
				printErrorf("output file cannot be closed (%s)", err)
				os.Exit(1)
			}
		}()
		printSuccessf("Writing to file %s", flagOut)
	} else {
		writeTo = os.Stdout
	}

	mdTitlef(1, "Package %s", doc.PackageName)

	mdTitle(2, "Overview")

	write(strings.TrimSpace(doc.PackageDoc) + "\n")

	mdHorizontalRule()

	/*
	 * Index
	 */

	mdTitle(2, "Index")

	if len(doc.Constants) > 0 {
		n := "Constants"
		writef("* [%s](#%s)\n", n, gitHubHeaderID(n))
	}

	if len(doc.Variables) > 0 {
		n := "Variables"
		writef("* [%s](#%s)\n", n, gitHubHeaderID(n))
	}

	if len(doc.Functions) > 0 {
		for _, pf := range doc.Functions.SortedKeys() {
			f := doc.Functions[pf]
			writef("* [func %s](#%s)\n", doc.Functions[pf].Signature(), gitHubHeaderID("func "+f.Signature()))
		}
	}

	if len(doc.Structs) > 0 {
		for _, t := range doc.Structs.SortedKeys() {
			c := doc.Structs[t]
			writef("* [type %s](#%s)\n", c.Name, gitHubHeaderID(c.Name))
			for _, sf := range c.Constructors.SortedKeys() {
				f := c.Constructors[sf]
				writef("  - [func %s](#%s)\n", f.Signature(), gitHubHeaderID("func "+f.Signature()))
			}
			for _, sf := range c.Methods.SortedKeys() {
				f := c.Methods[sf]
				writef("  - [func %s](#%s)\n", f.Signature(), gitHubHeaderID("func "+f.Signature()))
			}
		}
		writeln()
	}

	/*
	 * Content
	 */

	if len(doc.Constants) > 0 {
		mdHorizontalRule()
		mdTitle(2, "Constants")

		for _, f := range doc.Constants.SortedKeys() {
			mdTitlef(3, "const %s", doc.Constants[f].Name)
			writef("    %s\n\n", doc.Constants[f].Signature())
			if doc.Constants[f].Doc != "" {
				writef("%s\n\n", doc.Constants[f].Doc)
			}
		}
	}

	if len(doc.Variables) > 0 {
		mdHorizontalRule()
		mdTitle(2, "Variables")

		for _, f := range doc.Variables.SortedKeys() {
			mdTitlef(3, "var %s", doc.Variables[f].Name)
			writef("    %s\n\n", doc.Variables[f].Signature())
			if doc.Variables[f].Doc != "" {
				writef("%s\n\n", doc.Variables[f].Doc)
			}
		}
	}

	if len(doc.Functions) > 0 {
		mdHorizontalRule()
		mdTitle(2, "Functions")

		for _, df := range doc.Functions.SortedKeys() {
			f := doc.Functions[df]
			mdTitlef(3, "func %s",
				f.Signature(),
			)
			writef("    %s\n\n", f.Signature())
			if f.Doc != "" {
				writef("%s\n\n", f.Doc)
			}
		}
	}

	if len(doc.Structs) > 0 {
		mdHorizontalRule()
		mdTitle(2, "Types")

		for _, t := range doc.Structs.SortedKeys() {
			c := doc.Structs[t]
			mdTitlef(3, "type %s", c.Name)

			writef("    type %s struct {\n", c.Name)
			for _, f := range c.Fields {
				writef("         %s\n", f)
			}
			writef("    }\n")

			if c.Doc != "" {
				writef("%s\n\n", c.Doc)
			}
			for _, pf := range c.Constructors.SortedKeys() {
				f := c.Constructors[pf]
				mdTitlef(4, "func %s",
					f.Signature(),
				)
				writef("    %s\n\n", f.Signature())
				if f.Doc != "" {
					writef("%s\n\n", f.Doc)
				}
			}
			for _, pf := range c.Methods.SortedKeys() {
				f := c.Methods[pf]
				mdTitlef(4, "func %s",
					f.Signature(),
				)
				writef("    %s\n\n", f.Signature())
				if f.Doc != "" {
					writef("%s\n\n", f.Doc)
				}
			}
		}
	}

	printSuccessf("Successfully written to %s", flagOut)
}

func mdTitle(level int, title string) {
	if level < 1 || level > 6 {
		panic("level between 1 and 6")
	}

	switch level {
	case 1:
		writef("%s\n%s\n\n", title, strings.Repeat("=", len(title)))
	case 2:
		writef("%s\n%s\n\n", title, strings.Repeat("-", len(title)))
	default:
		writef("%s %s\n\n", strings.Repeat("#", level), title)
	}
}

func mdTitlef(level int, title string, a ...interface{}) {
	mdTitle(level, fmt.Sprintf(title, a...))
}

func mdHorizontalRule() {
	if flagNoHR {
		writeln()
		return
	}
	writeln("\n- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
}

var reGitHubHeaderBanned = regexp.MustCompile(`[^\P{P}-]`)
var reGitHubUpperEncoded = regexp.MustCompile(`%%{[a-f0-9]{2}`)
var reWhiteSpaces = regexp.MustCompile(`\s+`)

var headerIndex = make(map[string]int)

func gitHubHeaderID(h string) string {
	result := h

	result = reWhiteSpaces.ReplaceAllString(result, "-")
	result = reGitHubHeaderBanned.ReplaceAllString(result, "")
	result = strings.ToLower(result)

	result = url.PathEscape(result)
	result = reGitHubUpperEncoded.ReplaceAllStringFunc(result, strings.ToUpper)

	if _, have := headerIndex[result]; have {
		headerIndex[result]++
		result = fmt.Sprintf("%s-%d", result, headerIndex[result])
	} else {
		headerIndex[result] = 0
	}

	return result
}
