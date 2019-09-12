package mysql

import (
	"fmt"
	"log"
	"strings"
)

const (
	defaultVarcharSize   = 191
	defaultVarbinarySize = 767
	autoIncrement        = "AUTO_INCREMENT"
)

// MySQL XXX
type MySQL struct {
	Engine  string
	Charset string
}

// Index XXX
type Index struct {
	columns []string
	name    string
}

// UniqueIndex XXX
type UniqueIndex struct {
	columns []string
	name    string
}

// PrimaryKey XXX
type PrimaryKey struct {
	columns []string
}

// HeaderTemplate XXX
func (mysql MySQL) HeaderTemplate() string {
	return `SET foreign_key_checks=0;
`
}

// FooterTemplate XXX
func (mysql MySQL) FooterTemplate() string {
	return `SET foreign_key_checks=1;`
}

// TableTemplate XXX
func (mysql MySQL) TableTemplate() string {
	return `
DROP TABLE IF EXISTS {{ .Name }};

CREATE TABLE {{ .Name }} (
    {{ range .Columns }}{{ .ToSQL }},
    {{ end }}{{ range .Indexes.Sort  }}{{ .ToSQL }},
    {{end}}{{ .PrimaryKey.ToSQL }}
) ENGINE={{ .Dialect.Engine }} DEFAULT CHARACTER SET {{ .Dialect.Charset }};

`
}

// ToSQL convert mysql sql string from typeName and size
func (mysql MySQL) ToSQL(typeName string, size uint64) string {
	var columns []string

	switch typeName {
	case "int8", "*int8":
		return "TINYINT"
	case "int16", "*int16":
		return "SMALLINT"
	case "int32", "*int32":
		return "INTEGER"
	case "int64", "*int64", "sql.NullInt64":
		return "BIGINT"
	case "uint8", "*uint8":
		return "TINYINT unsigned"
	case "uint16", "*uint16":
		return "SMALLINT unsigned"
	case "uint32", "*uint32":
		return "INTEGER unsigned"
	case "uint64", "*uint64":
		return "BIGINT unsigned"
	case "float32", "*float32":
		return "FLOAT"
	case "float64", "*float64", "sql.NullFloat64":
		return "DOUBLE"
	case "string", "*string", "sql.NullString":
		return varchar(size)
	case "[]uint8", "sql.RawBytes":
		return varbinary(size)
	case "bool", "*bool", "sql.NullBool":
		return "TINYINT(1)"
	case "tinytext":
		return "TINYTEXT"
	case "text":
		return "TEXT"
	case "mediumtext":
		return "MEDIUMTEXT"
	case "longtext":
		return "LONGTEXT"
	case "tinyblob":
		return "TINYBLOB"
	case "blob":
		return "BLOB"
	case "mediumblob":
		return "MEDIUMBLOB"
	case "longblob":
		return "LONGBLOB"
	case "time":
		return "TIME"
	case "time.Time", "*time.Time":
		return "DATETIME"
	case "mysql.NullTime", "sql.NullTime":
		// sql.NullTime is supported from Go 1.13
		// mysql.NullTime is an alias of sql.NullTime now.
		// https://godoc.org/github.com/go-sql-driver/mysql#NullTime
		// https://github.com/go-sql-driver/mysql/blob/b2c03bcae3d4bf9658e63ad4b9b32581d810fb20/nulltime_go113.go#L31
		return "DATETIME"
	case "json.RawMessage", "*json.RawMessage":
		return "JSON"
	default:
		log.Fatalf("%s is not match.", typeName)
	}

	if size != 0 {
		columns = append(columns, fmt.Sprintf("(%d)", size))
	}
	return strings.Join(columns, "")
}

// Quote XXX
func (mysql MySQL) Quote(s string) string {
	return quote(s)
}

// AutoIncrement XXX
func (mysql MySQL) AutoIncrement() string {
	return autoIncrement
}

// Name XXX
func (i Index) Name() string {
	return i.name
}

// Columns XXX
func (i Index) Columns() []string {
	return i.columns
}

// ToSQL return index sql string
func (i Index) ToSQL() string {
	var columnsStr []string

	for _, c := range i.columns {
		columnsStr = append(columnsStr, quote(c))
	}

	return fmt.Sprintf("INDEX %s (%s)", quote(i.name), strings.Join(columnsStr, ", "))
}

// Name XXX
func (ui UniqueIndex) Name() string {
	return ui.name
}

// Columns XXX
func (ui UniqueIndex) Columns() []string {
	return ui.columns
}

// ToSQL return unique index sql string
func (ui UniqueIndex) ToSQL() string {
	var columnsStr []string
	for _, c := range ui.columns {
		columnsStr = append(columnsStr, quote(c))
	}

	return fmt.Sprintf("UNIQUE %s (%s)", quote(ui.name), strings.Join(columnsStr, ", "))
}

// Columns XXX
func (pk PrimaryKey) Columns() []string {
	return pk.columns
}

// ToSQL return primary key sql string
func (pk PrimaryKey) ToSQL() string {
	var columnsStr []string
	for _, c := range pk.columns {
		columnsStr = append(columnsStr, quote(c))
	}

	return fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(columnsStr, ", "))
}

// AddIndex XXX
func AddIndex(idxName string, columns ...string) Index {
	return Index{
		name:    idxName,
		columns: columns,
	}
}

// AddUniqueIndex XXX
func AddUniqueIndex(idxName string, columns ...string) UniqueIndex {
	return UniqueIndex{
		name:    idxName,
		columns: columns,
	}
}

// AddPrimaryKey XXX
func AddPrimaryKey(columns ...string) PrimaryKey {
	return PrimaryKey{
		columns: columns,
	}
}

func varchar(size uint64) string {
	if size == 0 {
		return fmt.Sprintf("VARCHAR(%d)", defaultVarcharSize)
	}

	return fmt.Sprintf("VARCHAR(%d)", size)
}

func varbinary(size uint64) string {
	if size == 0 {
		return fmt.Sprintf("VARBINARY(%d)", defaultVarbinarySize)
	}

	return fmt.Sprintf("VARBINARY(%d)", size)
}

func quote(s string) string {
	return fmt.Sprintf("`%s`", s)
}
