// Program enumflags generates methods for specified types to provide human readable
// representations of all cgo generated enum values mapping to both the golang type and the cgo type name
//
// Enum (non-bitset) types generate a method to implement the fmt.Formatter
// interface along with a map of type to name and a map of type to C identifier.
// The suffix name is the name of the Go value. If the type name is the prefix of
// the Go value, the type name is removed.
//
// Flag types (aka bitsets) are similar but also have methods for Add(), Del() and Has().
// A flag type (like flags in C) can be a single bit value or contain multiple values.
// If there is only a single value then IsSingleFlag() will return false.
//
// fmt.Formatter
//
// Both types support the fmt.Formatter interface with multiple verbs. "%s" will
// print the name derived from the enum value. "%c" will print the name derived
// from the imported C identifier. '%d' and '%b' will treat the value as an int64
// and format accordingly. For a Flag type, if the flag has multiple values the
// %s and %c verbs will print all the names (either suffix name or C name) of the
// values that are set joined by a pipe.
package main

import "C"
import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"slices"
	"strings"
	"text/template"
	"unicode"
)

var (
	enumTypes     = flag.String("enum_types", "", "comma separate list of enum types to generate")
	flagTypes     = flag.String("flag_types", "", "comma separate list of bitflag types to generate")
	output        = flag.String("output", "", "filename to use to output. If not specified the filename_generated.go is used.")
	printContents = flag.Bool("print", false, "If set will output the contents of the file generated.")
)

type EnumData struct {
	FileName      string
	TypeName      string
	SuffixMapName string
	SuffixEntries []Entry
	CgoMapName    string
	CgoEntries    []Entry
}

type Entry struct {
	Enum   string
	String string
}

//go:embed *.tmpl
var templates embed.FS

func main() {
	flag.Parse()
	pkg := os.Getenv("GOPACKAGE") // go generate sets GOPACKAGE env var
	if pkg == "" {
		log.Fatal("GOPACKAGE environment variable not set. Run with 'go generate'.")
	}
	if *enumTypes == "" && *flagTypes == "" {
		log.Fatal("-enum_types or -flag_types is required")
	}
	var allTypes, flagTypeNames []string
	if *enumTypes != "" {
		allTypes = append(allTypes, strings.Split(*enumTypes, ",")...)
	}
	if *flagTypes != "" {
		flagTypeNames = strings.Split(*flagTypes, ",")
		allTypes = append(allTypes, flagTypeNames...)
	}

	typesByFilename := make(map[string][]EnumData)
	for _, typeName := range allTypes {
		typeName = strings.TrimSpace(typeName)
		if typeName == "" {
			continue
		}
		typeData := buildEnumData(pkg, typeName)

		if typeData.FileName == "" {
			log.Fatalf("Found no enums for type: %s", typeName)
		}

		outFile := *output
		if *output == "" {
			outFile = addGenerated(typeData.FileName)
		}
		typesByFilename[outFile] = append(typesByFilename[outFile], typeData)
	}

	if len(typesByFilename) == 0 {
		log.Fatalf("No types found in list %s\n", allTypes)
	}
	t, err := template.New("enumflags").ParseFS(templates, "*.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	for filename, enums := range typesByFilename {
		var buf bytes.Buffer
		templateData := struct {
			PackageName string
			HasBitSet   bool
		}{
			PackageName: pkg,
			HasBitSet:   len(flagTypeNames) > 0,
		}
		if err := t.ExecuteTemplate(&buf, "header.tmpl", templateData); err != nil {
			log.Fatalf("failed to execute template: %v", err)
		}
		for _, typeData := range enums {
			if slices.Contains(flagTypeNames, typeData.TypeName) {
				if err := t.ExecuteTemplate(&buf, "bitset_formatter.tmpl", typeData); err != nil {
					log.Fatalf("failed to execute template: %v", err)
				}
			} else {
				if err := t.ExecuteTemplate(&buf, "enum_formatter.tmpl", typeData); err != nil {
					log.Fatalf("failed to execute template: %v", err)
				}
			}
			if err := t.ExecuteTemplate(&buf, "name_maps.tmpl", typeData); err != nil {
				log.Fatalf("failed to execute template: %v", err)
			}
		}
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Printf(string(buf.Bytes()))
			log.Fatalf("failed to format generated code: %v", err)
		}
		if err := os.WriteFile(filename, formatted, 0644); err != nil {
			log.Fatalf("failed to write output file: %v", err)
		}
		log.Printf("Generated %s for type(s) %s", filename, strings.Join(allTypes, ","))
		if *printContents {
			fmt.Printf(string(formatted))
		}
	}
}

func addGenerated(filename string) string {
	if strings.HasSuffix(filename, ".go") {
		filename = strings.TrimSuffix(filename, ".go")
	}
	return filename + "_generated.go"
}

// buildEnumData looks at the specified package for a
func buildEnumData(pkg, typeName string) EnumData {
	var typeData = EnumData{
		TypeName:      typeName,
		SuffixMapName: mapNameIdentifier(typeName, "NameMap"),
		CgoMapName:    mapNameIdentifier(typeName, "CgoNameMap"),
	}

	findTargetASTFile(pkg, func(filename string, fileAST *ast.File) bool {
		// for each `const FooTypeValue = FooType(C.FOO_VALUE)`
		// enumName will be "FooTypeValue" and arg will be "C.FOO_VALUE"
		visitEnumTypes(fileAST, typeName, func(enumName, cgoName string) bool {
			if typeData.FileName == "" {
				typeData.FileName = filename
			}
			stringName := enumName
			// Remove the prefix of the type if present
			// "FooTypeValue" => "Value"
			if strings.HasPrefix(enumName, typeName) {
				stringName = strings.TrimPrefix(enumName, typeName)
			}
			typeData.SuffixEntries = append(typeData.SuffixEntries, Entry{enumName, stringName})
			if cgoName != "" {
				typeData.CgoEntries = append(typeData.CgoEntries, Entry{enumName, cgoName})
			}
			return true
		})
		return true
	})
	return typeData
}

// visitEnumTypes walks the AST looking for const declarations of the specified type name.
// For any declaration like `const FooTypeValue = FooType(C.FOO_VALUE)`, f is called with the
// identifier of the constant (ie `FooTypeValue`) and the list of arguments (ie `C.FOO_VALUE`)
// passed to the type call. If the value is not imported from C the cName will be empty.
func visitEnumTypes(file *ast.File, typeName string, yield func(goName string, cName string) bool) {
	var inIota bool
	visitConstantDecls(file, func(constName string, constType ast.Expr, constValue ast.Expr) bool {
		// Handle `FooTypeTow FooType = iota`
		if constType != nil {
			if typeIdent, ok := constType.(*ast.Ident); ok && typeIdent.Name == typeName {
				if valueIdent, ok := constValue.(*ast.Ident); ok && valueIdent.Name == "iota" {
					inIota = true
					return yield(constName, "")
				}
			}
		}
		// Handle `IotaTypeTwo` with no type or value specified
		if inIota && constType == nil && constValue == nil {
			return yield(constName, "")
		}
		if inIota && constType != nil {
			inIota = false
		}

		if call, ok := constValue.(*ast.CallExpr); ok {
			// CallExpr: `FooThing("value-one")`
			if ident, ok := call.Fun.(*ast.Ident); ok {
				// Make sure FooThing matches typeName
				if ident.Name == typeName && len(call.Args) == 1 {

					// If the value expression is like "C.FOO_VALUE", extract "FOO_VALUE"
					if sel, ok := call.Args[0].(*ast.SelectorExpr); ok {
						if id, ok := sel.X.(*ast.Ident); ok && id.Name == "C" {
							if !yield(constName, sel.Sel.Name) {
								return false
							}
						}
					} else {
						if !yield(constName, "") {
							return false
						}
					}

				}
			}
		}
		return true
	})
}

// visitConstantDecls will call yield with every const declaration in the specified package.
//
// Examples:
//
//	const CarTypeOne = CarType(C.CAR_ENUM_ONE)
//	const FooFlagOne FooFlag = "foo flag"
//	const FooFlagOne = FooFlag("foo flag")
//	const IotaTypeOne IotaType = iota
//	const IotaTypeTwo
func visitConstantDecls(file *ast.File, yield func(name string, typeName ast.Expr, value ast.Expr) bool) {
	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		// Look for package level const declarations
		//
		// const ( FooThing = FooType("value-one") )
		if !ok || genDecl.Tok != token.CONST {
			return true // Not a const declaration, continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			// valueSpec: `FooThingValueOne = FooThing("value-one")`
			// var isIota bool
			for i, name := range valueSpec.Names {
				var value ast.Expr
				if len(valueSpec.Values) > i {
					value = valueSpec.Values[i]
				}
				if !yield(name.Name, valueSpec.Type, value) {
					return false
				}
			}
		}
		return true
	})
}

// findTargetASTFile calls yield with each ast.File for the specified package.
func findTargetASTFile(pkg string, yield func(string, *ast.File) bool) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("failed to parse directory: %v", err)
	}
	for _, p := range pkgs {
		if p.Name == pkg {
			for name, file := range p.Files {
				if !yield(name, file) {
					return
				}
			}
		}
	}
}

func mapNameIdentifier(typeName, suffix string) string {
	name := typeName + suffix
	r := []rune(name) // Convert string to slice of runes to handle Unicode characters
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
