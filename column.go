package ddlmaker

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kayac/ddl-maker/dialect"
)

// column is mapping struct field value.
type column struct {
	name     string
	typeName string
	tag      string
	dialect  dialect.Dialect
}

func newColumn(name, typeName, tag string, d dialect.Dialect) column {
	return column{
		name:     name,
		typeName: typeName,
		tag:      tag,
		dialect:  d,
	}
}

func (c column) size() (uint64, error) {
	specs := c.specs()
	if specs["size"] == "" {
		return 0, nil
	}

	return strconv.ParseUint(specs["size"], 10, 64)
}

func (c column) specs() map[string]string {
	elems := strings.Split(c.tag, ",")
	specs := make(map[string]string, len(elems))
	for _, elem := range elems {
		ss := strings.Split(elem, "=")
		switch len(ss) {
		case 1:
			specs[ss[0]] = ""
		case 2:
			specs[ss[0]] = ss[1]
		}
	}

	return specs
}

func (c column) attribute() string {
	var attributes []string
	specs := c.specs()

	if _, ok := specs["null"]; ok {
		attributes = append(attributes, "NULL")
	} else {
		attributes = append(attributes, "NOT NULL")
	}

	if defaultVal, ok := specs["default"]; ok {
		attributes = append(attributes, "DEFAULT")
		attributes = append(attributes, defaultVal)
	}

	if _, ok := specs["auto"]; ok {
		attributes = append(attributes, c.dialect.AutoIncrement())
	}

	return strings.Join(attributes, " ")
}

func (c column) Name() string {
	return c.name
}

// ToSQL is convert struct value to sql.
func (c column) ToSQL() string {
	var columnType string
	specs := c.specs()

	if typeName, ok := specs["type"]; ok {
		columnType = typeName
	} else {
		columnType = c.typeName
	}

	name := c.dialect.Quote(c.name)
	size, err := c.size()
	if err != nil {
		log.Fatalf("error size parse error %v", err)
	}

	sql := c.dialect.ToSQL(columnType, size)
	attribute := c.attribute()

	return fmt.Sprintf("%s %s %s", name, sql, attribute)
}
