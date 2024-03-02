//go:build ignore
// +build ignore

package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const tpl = `// Code generated. DO NOT EDIT!

package {{ .Package }}

// {{ .Map }} is generated form a fileset
var {{ .Map }} = map[string]string{
{{ range $constant, $content := .Files }}` + "\t" + `"{{ $constant }}": ` + "`{{ $content }}`" + `,
{{ end }}}

// {{ .Map }}Checksums is generated from a fileset and contains files checksums
var {{ .Map }}Checksums = map[string]string{
{{ range $constant, $content := .Checksums }}` + "\t" + `"{{ $constant }}": "{{ $content }}",
{{ end }}}
`

var bundleTpl = template.Must(template.New("").Parse(tpl))

type Bundle struct {
	Package    string
	Map        string
	ImportPath string
	Files      map[string]string
	Checksums  map[string]string
}

func (b *Bundle) Write(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bundleTpl.Execute(f, b)
}

func NewBundle(pkg, mapName, importPath string) *Bundle {
	return &Bundle{
		Package:    pkg,
		Map:        mapName,
		ImportPath: importPath,
		Files:      make(map[string]string),
		Checksums:  make(map[string]string),
	}
}

func readFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return data
}

func checksum(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func basename(filename string) string {
	return path.Base(filename)
}

func stripExtension(filename string) string {
	filename = strings.TrimSuffix(filename, path.Ext(filename))
	return strings.Replace(filename, " ", "_", -1)
}

func glob(pattern string) []string {
	files, _ := filepath.Glob(pattern)
	for i := range files {
		if strings.Contains(files[i], "\\") {
			files[i] = strings.Replace(files[i], "\\", "/", -1)
		}
	}
	return files
}

func concat(files []string) string {
	var b strings.Builder
	for _, file := range files {
		b.Write(readFile(file))
	}
	return b.String()
}

func generateBundle(bundleFile, pkg, mapName string, srcFiles []string) {
	log.Printf("Generating %s ...\n", bundleFile)
	bundle := NewBundle(pkg, mapName, pkg)

	for _, srcFile := range srcFiles {
		data := readFile(srcFile)
		filename := stripExtension(basename(srcFile))

		bundle.Files[filename] = string(data)
		bundle.Checksums[filename] = checksum(data)
	}

	bundle.Write(bundleFile)
	log.Printf("Generating %s [done]\n", bundleFile)
}

func main() {
	generateBundle(
		"autogen/db/postgres/db_sql_migration.go",
		"postgres",
		"DatabaseSQLMigration",
		glob("internal/db/postgres/sql/*.sql"),
	)
}
