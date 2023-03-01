package main

import (
	"go/ast"
	"reflect"
	"testing"
)

func TestQueryBuilder_InsertQuery(t *testing.T) {
	type fields struct {
		TableName string
		Fields    []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "t1", fields: fields{TableName: "users", Fields: []string{"id", "name", "invitation_code"}}, want: "INSERT INTO users (id, name, invitation_code) VALUES (:id, :name, :invitation_code)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := QueryBuilder{
				TableName: tt.fields.TableName,
				Fields:    tt.fields.Fields,
			}
			if got := q.InsertQuery(); got != tt.want {
				t.Errorf("InsertQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "t1", args: args{s: "cat"}, want: "cats"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pluralize(tt.args.s); got != tt.want {
				t.Errorf("Pluralize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "t1", args: args{s: "helloWorld"}, want: "hello_world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeCase(tt.args.s); got != tt.want {
				t.Errorf("SnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFields(t *testing.T) {
	type args struct {
		s *ast.StructType
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "t1", args: args{s: &ast.StructType{Fields: &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: "ID"}}, Type: &ast.Ident{Name: "int"}}, {Names: []*ast.Ident{{Name: "Name"}}, Type: &ast.Ident{Name: "string"}}}}}}, want: []string{"ID", "Name"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFields(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractStructs(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want []QueryBuilder
	}{
		{name: "t1", args: args{filePath: "./example.go"}, want: []QueryBuilder{{TableName: "users", Fields: []string{"id", "name", "invitation_code"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractStructs(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractStructs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeTableName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "t1", args: args{name: "User"}, want: "users"},
		{name: "t1", args: args{name: "PostCount"}, want: "post_counts"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeTableName(tt.args.name); got != tt.want {
				t.Errorf("normalizeTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
