package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"unicode"
)

type QueryBuilder struct {
	TableName string
	Fields    []string
}

func (q QueryBuilder) InsertQuery() string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (:%s)", strings.ToLower(q.TableName), strings.Join(q.Fields, ", "), strings.Join(q.Fields, ", :"))
}

func main() {
	filePath := "./example.go"
	// Parse the file's AST
	structs := ExtractStructs(filePath)
	// Generate MySQL insert queries for each struct
	for _, s := range structs {
		fmt.Println(s.InsertQuery())
	}
}

func ExtractStructs(filePath string) []QueryBuilder {
	fset := token.NewFileSet()
	parsedAST, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Find all struct declarations
	var structs []QueryBuilder
	ast.Inspect(parsedAST, func(n ast.Node) bool {
		// try to convert n to ast.TypeSpec
		if typeSpec, isTypeSpec := n.(*ast.TypeSpec); isTypeSpec {
			s, isStructType := typeSpec.Type.(*ast.StructType)
			// check if conversion was successful
			if !isStructType {
				return true
			}
			// get the struct name
			structName := typeSpec.Name.Name
			// get Fields helper function
			fields := getFields(s)
			structs = append(structs, QueryBuilder{TableName: normalizeTableName(structName), Fields: fields})
		}
		return true
	})
	return structs
}

func getFields(s *ast.StructType) []string {
	fields := make([]string, len(s.Fields.List))
	for i, field := range s.Fields.List {
		if len(field.Names) == 0 {
			continue
		}
		fields[i] = SnakeCase(field.Names[0].Name)
	}
	return fields
}

// normalize the struct name to lowercase, pluralize it and apply snakeCase
// for example, User -> users, ReviewPost -> review_posts
func normalizeTableName(name string) string {
	return Pluralize(SnakeCase(name))
}

// Pluralize the string
func Pluralize(s string) string {
	if strings.HasSuffix(s, "y") {
		return strings.TrimSuffix(s, "y") + "ies"
	}
	return s + "s"
}

// SnakeCase converts UpperCamelCase to snake_case
func SnakeCase(s string) string {
	var str strings.Builder
	var prev rune
	for i, r := range s {
		// check if we should insert a underscore
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(prev) {
			str.WriteRune('_')
		}
		// lower case all characters
		str.WriteRune(unicode.ToLower(r))
		prev = r
	}
	return str.String()
}
