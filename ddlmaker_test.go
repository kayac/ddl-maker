package ddlmaker

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/kayac/ddl-maker/dialect/mysql"
)

type Test1 struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t1 Test1) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

type Test2 struct {
	ID        uint64
	Test1ID   uint64
	Comment   sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t2 Test2) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id", "created_at")
}

func TestNewMaker(t *testing.T) {
	conf := Config{}
	_, err := NewMaker(conf)
	if err == nil {
		t.Fatal("Not set driver name", err)
	}

	conf = Config{
		DB: DBConfig{Driver: "dummy"},
	}
	_, err = NewMaker(conf)
	if err == nil {
		t.Fatal("Set unsupport driver name", err)
	}

	conf = Config{
		DB: DBConfig{Driver: "mysql"},
	}
	_, err = NewMaker(conf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddStruct(t *testing.T) {
	dm, err := NewMaker(Config{
		DB: DBConfig{Driver: "mysql"},
	})
	if err != nil {
		t.Fatal(err)
	}
	dm.AddStruct(Test1{}, Test2{})
	if len(dm.Structs) != 2 {
		t.Fatal("[error] add stuct")
	}

	err = dm.AddStruct(Test1{})
	if err != nil {
		t.Fatal("[error] add duplicate struct")
	}
}

func TestGenerate(t *testing.T) {
	m := mysql.MySQL{}
	generatedDDL := fmt.Sprintf(`%s
DROP TABLE IF EXISTS %s;

CREATE TABLE %s (
    %s BIGINT unsigned NOT NULL,
    %s VARCHAR(191) NOT NULL,
    %s DATETIME NOT NULL,
    %s DATETIME NOT NULL,
    PRIMARY KEY (%s)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

%s`, m.HeaderTemplate(), m.Quote("test1"), m.Quote("test1"), m.Quote("id"), m.Quote("name"), m.Quote("created_at"), m.Quote("updated_at"), m.Quote("id"), m.FooterTemplate())

	tmpFile, err := ioutil.TempFile("", "create_ddl_")

	if err != nil {
		t.Fatal("error create tmp file", err)
	}
	defer os.Remove(tmpFile.Name())

	dm, err := NewMaker(Config{
		OutFilePath: tmpFile.Name(),
		DB: DBConfig{
			Driver:  "mysql",
			Engine:  "InnoDB",
			Charset: "utf8mb4",
		},
	})
	if err != nil {
		t.Fatal("error new maker", err)
	}

	err = dm.AddStruct(Test1{})
	if err != nil {
		t.Fatal("error add struct", err)
	}

	err = dm.Generate()
	if err != nil {
		t.Fatal("error generate ddl", err)
	}

	b, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal("error read file", err)
	}

	if string(b) != generatedDDL {
		t.Fatalf("generatedDDL: %s \n checkDDLL: %s \n", string(b), generatedDDL)
	}
}
