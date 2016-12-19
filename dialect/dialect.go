package dialect

import (
	"fmt"
	"sort"

	"github.com/kayac/ddl-maker/dialect/mysql"
)

// Dialect XXX
type Dialect interface {
	HeaderTemplate() string
	FooterTemplate() string
	TableTemplate() string
	ToSQL(typeNmae string, size uint64) string
	Quote(string) string
	AutoIncrement() string
}

// Table XXX
type Table interface {
	Name() string
	PrimaryKey() PrimaryKey
	Indexes() Indexes
	Columns() []Column
	Dialect() Dialect
}

// Column XXX
type Column interface {
	Name() string
	ToSQL() string
}

// PrimaryKey XXX
type PrimaryKey interface {
	Columns() []string
	ToSQL() string
}

// Indexes XXX
type Indexes []Index

// Index XXX
type Index interface {
	Name() string
	Columns() []string
	ToSQL() string
}

// Sort is sort index value by alphabets
func (indexes Indexes) Sort() Indexes {
	indexMap := make(map[string]Index, 0)
	var indexStr []string
	var sortIndexes []Index

	for _, index := range indexes {
		indexStr = append(indexStr, index.ToSQL())
		indexMap[index.ToSQL()] = index
	}

	sort.Strings(indexStr)
	for _, key := range indexStr {
		sortIndexes = append(sortIndexes, indexMap[key])
	}

	return sortIndexes
}

// NewDialect return Dialect
func NewDialect(driver, engine, charset string) (Dialect, error) {
	var d Dialect

	switch driver {
	case "mysql":
		d = &mysql.MySQL{
			Engine:  engine,
			Charset: charset,
		}
	default:
		return d, fmt.Errorf("No such driver: %s", driver)
	}

	return d, nil
}
