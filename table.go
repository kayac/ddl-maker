package ddlmaker

import (
	"github.com/kayac/ddl-maker/dialect"
)

// Table is mapping struct info
type table struct {
	name        string
	primaryKey  dialect.PrimaryKey
	foreignKeys dialect.ForeignKeys
	columns     []dialect.Column
	indexes     dialect.Indexes
	dialect     dialect.Dialect
}

func newTable(name string, pk dialect.PrimaryKey, fks dialect.ForeignKeys, columns []dialect.Column, indexes dialect.Indexes, d dialect.Dialect) table {
	return table{
		name:        name,
		primaryKey:  pk,
		foreignKeys: fks,
		columns:     columns,
		indexes:     indexes,
		dialect:     d,
	}
}

func (t table) Name() string {
	return t.dialect.Quote(t.name)
}

func (t table) PrimaryKey() dialect.PrimaryKey {
	return t.primaryKey
}

func (t table) ForeignKeys() dialect.ForeignKeys {
	return t.foreignKeys
}

func (t table) Columns() []dialect.Column {
	return t.columns
}

func (t table) Indexes() dialect.Indexes {
	return t.indexes
}

func (t table) Dialect() dialect.Dialect {
	return t.dialect
}
