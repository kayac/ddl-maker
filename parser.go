package ddlmaker

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/serenize/snaker"
)

// Table is for type assertion
type Table interface {
	Table() string
}

// PrimaryKey is for type assertion
type PrimaryKey interface {
	PrimaryKey() dialect.PrimaryKey
}

// Index is for type assertion
type Index interface {
	Indexes() dialect.Indexes
}

func (dm *DDLMaker) parse() {
	for _, s := range dm.Structs {
		val := reflect.Indirect(reflect.ValueOf(s))
		rt := val.Type()

		var columns []dialect.Column
		for i := 0; i < rt.NumField(); i++ {
			rtField := rt.Field(i)
			column := parseField(rtField, dm.Dialect)
			columns = append(columns, column)
		}

		table := parseTable(s, columns, dm.Dialect)
		dm.Tables = append(dm.Tables, table)
	}
}

func parseField(field reflect.StructField, d dialect.Dialect) dialect.Column {
	tagStr := strings.Replace(field.Tag.Get(TAGPREFIX), " ", "", -1)

	var typeName string
	if field.Type.PkgPath() != "" {
		// ex) time.Time
		pkgName := field.Type.PkgPath()
		if strings.Contains(pkgName, "/") {
			pkgs := strings.Split(pkgName, "/")
			pkgName = pkgs[len(pkgs)-1]
		}
		typeName = fmt.Sprintf("%s.%s", pkgName, field.Type.Name())
	} else if field.Type.Kind().String() == "ptr" {
		// pointer type
		typeName = fmt.Sprintf("*%s", field.Type.Elem())
	} else {
		typeName = field.Type.Name()
	}

	return newColumn(snaker.CamelToSnake(field.Name), typeName, tagStr, d)
}

func parseTable(s interface{}, columns []dialect.Column, d dialect.Dialect) dialect.Table {
	var tableName string
	var primaryKey dialect.PrimaryKey
	var indexes dialect.Indexes

	if v, ok := s.(Table); ok {
		tableName = snaker.CamelToSnake(v.Table())
	} else {
		val := reflect.Indirect(reflect.ValueOf(s))
		tableName = snaker.CamelToSnake(val.Type().Name())
	}
	if v, ok := s.(PrimaryKey); ok {
		primaryKey = v.PrimaryKey()
	} else {
		panic(`you must implement PrimaryKey interface: ` + tableName)
	}
	if v, ok := s.(Index); ok {
		indexes = v.Indexes()
	}

	return newTable(tableName, primaryKey, columns, indexes, d)
}
