package ddlmaker

import (
	"reflect"
	"testing"

	"github.com/kayac/ddl-maker/dialect/mysql"
)

func TestSize(t *testing.T) {
	c := column{
		name: "dummy",
	}

	if size, err := c.size(); size != 0 || err != nil {
		t.Fatal("parse size error")
	}

	c = column{
		name: "dummy",
		tag:  "size=10",
	}

	if size, err := c.size(); size != 10 || err != nil {
		t.Fatal("parse size error")
	}
}

func TestSpecs(t *testing.T) {
	c := column{
		name: "name",
		tag:  "size=10,pk,default=jon",
	}

	specs := map[string]string{
		"size":    "10",
		"pk":      "",
		"default": "jon",
	}

	if !reflect.DeepEqual(c.specs(), specs) {
		t.Fatalf("parse tag error. result: %q", c.specs())
	}
}

func TestAttribute(t *testing.T) {
	c := column{dialect: mysql.MySQL{}}

	if c.attribute() != "NOT NULL" {
		t.Fatalf("error column attribute. result:%s", c.attribute())
	}

	c.tag = "null"
	if c.attribute() != "NULL" {
		t.Fatalf("error column attribute. result:%s", c.attribute())
	}

	c.tag = "default=0"
	if c.attribute() != "NOT NULL DEFAULT 0" {
		t.Fatalf("error column attribute. result:%s", c.attribute())
	}

	c.tag = "auto"
	if c.attribute() != "NOT NULL AUTO_INCREMENT" {
		t.Fatalf("error column attribute. result:%s", c.attribute())
	}
}

func TestToSQL(t *testing.T) {
	c := column{
		typeName: "int64",
		name:     "id",
		dialect:  mysql.MySQL{},
	}

	if c.ToSQL() != "`id` BIGINT NOT NULL" {
		t.Fatalf("error ToSQL. result: %s", c.ToSQL())
	}

	c.typeName = "uint64"
	if c.ToSQL() != "`id` BIGINT unsigned NOT NULL" {
		t.Fatalf("error ToSQL. result: %s", c.ToSQL())
	}

	c = column{
		typeName: "string",
		name:     "description",
		tag:      "size=20,null",
		dialect:  mysql.MySQL{},
	}

	if c.ToSQL() != "`description` VARCHAR(20) NULL" {
		t.Fatalf("error ToSQL. result: %s", c.ToSQL())
	}

	c = column{
		typeName: "string",
		name:     "comment",
		tag:      "null,type=text",
		dialect:  mysql.MySQL{},
	}

	if c.ToSQL() != "`comment` TEXT NULL" {
		t.Fatalf("error ToSQL. result: %s", c.ToSQL())
	}
}
