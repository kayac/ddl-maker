package dialect

import (
	"testing"

	"github.com/kayac/ddl-maker/dialect/mysql"
)

func TestNewDialect(t *testing.T) {
	_, err := NewDialect("", "", "")
	if err == nil {
		t.Fatal("error not set driver")
	}

	_, err = NewDialect("mysql", "", "")
	if err != nil {
		t.Fatalf("error new dialect:%s error", "mysql")
	}
}

func TestSort(t *testing.T) {
	var indexes Indexes

	idx1 := mysql.AddUniqueIndex("fuga_idx", "fuga")
	indexes = append(indexes, idx1)
	idx2 := mysql.AddIndex("hoge_idx", "hoge")
	indexes = append(indexes, idx2)

	idxes := indexes.Sort()
	if len(idxes) != 2 {
		t.Fatal("error sort Indexes", idxes)
	}
	if idxes[0].ToSQL() != idx2.ToSQL() {
		t.Fatal("error sort index", idxes[0].ToSQL())
	}

}
