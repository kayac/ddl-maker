package mysql

import (
	"fmt"
	"testing"
)

func TestToSQL(t *testing.T) {
	m := MySQL{}

	testcases := []struct {
		typeName string
		size     uint64
		output   string
	}{
		{"bool", 0, "TINYINT(1)"},
		{"*bool", 0, "TINYINT(1)"},
		{"sql.NullBool", 0, "TINYINT(1)"},
		{"int8", 0, "TINYINT"},
		{"*int8", 0, "TINYINT"},
		{"int16", 0, "SMALLINT"},
		{"*int16", 0, "SMALLINT"},
		{"int32", 0, "INTEGER"},
		{"*int32", 0, "INTEGER"},
		{"sql.NullInt32", 0, "INTEGER"},
		{"int64", 0, "BIGINT"},
		{"*int64", 0, "BIGINT"},
		{"sql.NullInt64", 0, "BIGINT"},
		{"uint8", 0, "TINYINT unsigned"},
		{"*uint8", 0, "TINYINT unsigned"},
		{"uint16", 0, "SMALLINT unsigned"},
		{"*uint16", 0, "SMALLINT unsigned"},
		{"uint32", 0, "INTEGER unsigned"},
		{"*uint32", 0, "INTEGER unsigned"},
		{"uint64", 0, "BIGINT unsigned"},
		{"*uint64", 0, "BIGINT unsigned"},
		{"float32", 0, "FLOAT"},
		{"*float32", 0, "FLOAT"},
		{"float64", 0, "DOUBLE"},
		{"*float64", 0, "DOUBLE"},
		{"sql.NullFloat64", 0, "DOUBLE"},
		{"string", 0, fmt.Sprintf("VARCHAR(%d)", defaultVarcharSize)},
		{"*string", 0, fmt.Sprintf("VARCHAR(%d)", defaultVarcharSize)},
		{"sql.NullString", 0, fmt.Sprintf("VARCHAR(%d)", defaultVarcharSize)},
		{"string", 10, "VARCHAR(10)"},
		{"*string", 10, "VARCHAR(10)"},
		{"sql.NullString", 10, "VARCHAR(10)"},
		{"[]uint8", 10, "VARBINARY(10)"},
		{"sql.RawBytes", 10, "VARBINARY(10)"},
		{"tinytext", 0, "TINYTEXT"},
		{"text", 0, "TEXT"},
		{"mediumtext", 0, "MEDIUMTEXT"},
		{"longtext", 0, "LONGTEXT"},
		{"tinyblob", 0, "TINYBLOB"},
		{"blob", 0, "BLOB"},
		{"mediumblob", 0, "MEDIUMBLOB"},
		{"longblob", 0, "LONGBLOB"},
		{"time", 0, "TIME"},
		{"time.Time", 0, "DATETIME"},
		{"time.Time", 6, "DATETIME(6)"},
		{"mysql.NullTime", 0, "DATETIME"}, // https://godoc.org/github.com/go-sql-driver/mysql#NullTime
		{"sql.NullTime", 0, "DATETIME"},   // from Go 1.13
		{"date", 0, "DATE"},
		{"geometry", 0, "GEOMETRY"},
		{"json.RawMessage", 0, "JSON"},
	}

	for _, tc := range testcases {
		if m.ToSQL(tc.typeName, tc.size) != tc.output {
			t.Fatalf("error %s to sql %s. but result %s", tc.typeName, tc.output, m.ToSQL(tc.typeName, tc.size))
		}
	}
}

func TestQuote(t *testing.T) {
	column := "id"

	if quote(column) != "`id`" {
		t.Fatalf("error %s quote. result:%s ", column, quote(column))
	}
}

func TestAuotIncrement(t *testing.T) {
	m := MySQL{}
	if m.AutoIncrement() != autoIncrement {
		t.Fatalf("error auto increament: %s. result:%s", autoIncrement, m.AutoIncrement())
	}
}

func TestAddIndex(t *testing.T) {
	index := AddIndex("player_id_idx", "player_id")
	if index.ToSQL() != "INDEX `player_id_idx` (`player_id`)" {
		t.Fatal("[error] parse player_id_idx. ", index.ToSQL())
	}

	index = AddIndex("player_entry_id_idx", "player_id", "entry_id")
	if index.ToSQL() != "INDEX `player_entry_id_idx` (`player_id`, `entry_id`)" {
		t.Fatal("[error] parse player_entry_id_idx", index.ToSQL())
	}
}

func TestAddUniqIndex(t *testing.T) {
	uniqIndex := AddUniqueIndex("player_id_idx", "player_id")
	if uniqIndex.ToSQL() != "UNIQUE `player_id_idx` (`player_id`)" {
		t.Fatal("[error] parse unique player_id_idx", uniqIndex.ToSQL())
	}

	uniqIndex = AddUniqueIndex("player_entry_id_idx", "player_id", "entry_id")
	if uniqIndex.ToSQL() != "UNIQUE `player_entry_id_idx` (`player_id`, `entry_id`)" {
		t.Fatal("[error] parse unique player_entry_id_idx", uniqIndex.ToSQL())
	}
}

func TestAddFullTextIndex(t *testing.T) {
	fullTextIndex := AddFullTextIndex("full_text_idx", "content")
	if fullTextIndex.ToSQL() != "FULLTEXT `full_text_idx` (`content`)" {
		t.Fatal("[error] parse full_text_idx", fullTextIndex.ToSQL())
	}

	fullTextIndex = AddFullTextIndex("full_text_idx", "content", "title")
	if fullTextIndex.ToSQL() != "FULLTEXT `full_text_idx` (`content`, `title`)" {
		t.Fatal("[error] parse full_text_idx", fullTextIndex.ToSQL())
	}

	fullTextIndex = AddFullTextIndex("full_text_idx", "content").WithParser("ngram")
	if fullTextIndex.ToSQL() != "FULLTEXT `full_text_idx` (`content`) WITH PARSER `ngram`" {
		t.Fatal("[error] parse full_text_idx", fullTextIndex.ToSQL())
	}
}

func TestAddAddSpatialIndex(t *testing.T) {
	spatialIndex := AddSpatialIndex("geometry_idx", "g")
	if spatialIndex.ToSQL() != "SPATIAL KEY `geometry_idx` (`g`)" {
		t.Fatal("[error] parse geometry_idx", spatialIndex.ToSQL())
	}

	spatialIndex = AddSpatialIndex("geometry_idx", "g", "g1")
	if spatialIndex.ToSQL() != "SPATIAL KEY `geometry_idx` (`g`, `g1`)" {
		t.Fatal("[error] parse geometry_idx", spatialIndex.ToSQL())
	}
}

func TestAddPrimaryKey(t *testing.T) {
	pk := AddPrimaryKey("id")
	if pk.ToSQL() != "PRIMARY KEY (`id`)" {
		t.Fatal("[error] parse primary key", pk.ToSQL())
	}

	pk = AddPrimaryKey("id", "created_at")
	if pk.ToSQL() != "PRIMARY KEY (`id`, `created_at`)" {
		t.Fatal("[error] parse primary key", pk.ToSQL())
	}

	pk = AddPrimaryKey("created_at", "id")
	if pk.ToSQL() != "PRIMARY KEY (`created_at`, `id`)" {
		t.Fatal("[error] parse primary key", pk.ToSQL())
	}
}

func TestAddForeignKey(t *testing.T) {
	foreignColumns := []string{"player_id"}
	referenceColumns := []string{"id"}
	fk := AddForeignKey(foreignColumns, referenceColumns, "player")
	if fk.ToSQL() != "FOREIGN KEY (`player_id`) REFERENCES `player` (`id`)" {
		t.Fatal("[error] parse foreign key", fk.ToSQL())
	}

	foreignColumns = []string{"product_category", "product_id"}
	referenceColumns = []string{"category", "id"}
	fk = AddForeignKey(foreignColumns, referenceColumns, "product")
	if fk.ToSQL() != "FOREIGN KEY (`product_category`, `product_id`) REFERENCES `product` (`category`, `id`)" {
		t.Fatal("[error] parse foreign key", fk.ToSQL())
	}

	foreignColumns = []string{"product_category", "product_id"}
	referenceColumns = []string{"category", "id"}
	fk = AddForeignKey(foreignColumns, referenceColumns, "product", WithUpdateForeignKeyOption(ForeignKeyOptionNoAction), WithDeleteForeignKeyOption(ForeignKeyOptionNoAction))
	if fk.ToSQL() != "FOREIGN KEY (`product_category`, `product_id`) REFERENCES `product` (`category`, `id`)" {
		t.Fatal("[error] parse foreign key", fk.ToSQL())
	}

	foreignColumns = []string{"product_category", "product_id"}
	referenceColumns = []string{"category", "id"}
	fk = AddForeignKey(foreignColumns, referenceColumns, "product", WithUpdateForeignKeyOption(ForeignKeyOptionCascade))
	if fk.ToSQL() != "FOREIGN KEY (`product_category`, `product_id`) REFERENCES `product` (`category`, `id`) ON UPDATE CASCADE" {
		t.Fatal("[error] parse foreign key", fk.ToSQL())
	}
}
