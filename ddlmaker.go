package ddlmaker

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"text/template"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/pkg/errors"
)

const (
	// TAGPREFIX is struct tag field prefix
	TAGPREFIX = "ddl"
)

// DDLMaker XXX
type DDLMaker struct {
	config  Config
	Dialect dialect.Dialect
	Structs []interface{}
	Tables  []dialect.Table
}

// NewMaker return DDLMaker
func NewMaker(conf Config) (*DDLMaker, error) {
	dialect, err := dialect.NewDialect(conf.DB.Driver, conf.DB.Engine, conf.DB.Charset)
	if err != nil {
		return nil, errors.Wrap(err, "error NewDialect()")
	}

	return &DDLMaker{
		config:  conf,
		Dialect: dialect,
	}, nil
}

// AddStruct XXX
func (dm *DDLMaker) AddStruct(ss ...interface{}) error {
	pkgs := make(map[string]bool, 0)

	for _, s := range ss {
		if s == nil {
			return fmt.Errorf("nil is not supported")
		}

		val := reflect.Indirect(reflect.ValueOf(s))
		rt := val.Type()

		structName := fmt.Sprintf("%s.%s", rt.PkgPath(), rt.Name())
		if pkgs[structName] {
			return fmt.Errorf("%s is already added", structName)
		}

		dm.Structs = append(dm.Structs, s)
		pkgs[structName] = true
	}

	return nil
}

// Generate ddl file
func (dm *DDLMaker) Generate() error {
	log.Printf("start generate %s \n", dm.config.OutFilePath)
	err := dm.parse()
	if err != nil {
		return errors.Wrap(err, "error parse")
	}

	file, err := os.Create(dm.config.OutFilePath)
	if err != nil {
		return errors.Wrap(err, "error create ddl file")
	}
	defer file.Close()

	err = dm.generate(file)
	if err != nil {
		return errors.Wrap(err, "error generate")
	}

	log.Printf("done generate %s \n", dm.config.OutFilePath)

	return nil
}

func (dm *DDLMaker) generate(w io.Writer) error {
	header, err := template.New("header").Parse(dm.Dialect.HeaderTemplate())
	if err != nil {
		return errors.Wrap(err, "error parse header template")
	}

	footer, err := template.New("footer").Parse(dm.Dialect.FooterTemplate())
	if err != nil {
		return errors.Wrap(err, "error parse header footer")
	}

	tmpl, err := template.New("ddl").Parse(dm.Dialect.TableTemplate())
	if err != nil {
		return errors.Wrap(err, "error parse template")
	}

	header.Execute(w, nil)
	for _, table := range dm.Tables {
		err := tmpl.Execute(w, table)
		if err != nil {
			return errors.Wrap(err, "template execute error")
		}
	}
	footer.Execute(w, nil)

	return nil
}
