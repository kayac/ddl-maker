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
		{"mysql.NullTime", 0, "DATETIME"}, // https://godoc.org/github.com/go-sql-driver/mysql#NullTime
		{"date", 0, "DATE"},
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
