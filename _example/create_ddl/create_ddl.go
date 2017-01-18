package main

import (
	"flag"
	"log"

	"github.com/kayac/ddl-maker"
	ex "github.com/kayac/ddl-maker/_example"
)

func main() {
	var (
		driver      string
		engine      string
		charset     string
		outFilePath string
	)
	flag.StringVar(&driver, "d", "", "set driver")
	flag.StringVar(&driver, "driver", "", "set driver")
	flag.StringVar(&outFilePath, "o", "./sql/master.sql", "set ddl output file path")
	flag.StringVar(&outFilePath, "outfile", "./sql/master.sql", "set ddl output file path")
	flag.StringVar(&engine, "e", "InnoDB", "set driver engine")
	flag.StringVar(&engine, "engine", "InnoDB", "set driver engine")
	flag.StringVar(&charset, "c", "utf8mb4", "set driver charset")
	flag.StringVar(&charset, "charset", "utf8mb4", "set driver charset")
	flag.Parse()

	if driver == "" {
		log.Println("Please set driver name. -d or -driver")
		return
	}
	if outFilePath == "" {
		log.Println("Please set outFilePath. -o or -outfile")
		return
	}

	conf := ddlmaker.Config{
		DB: ddlmaker.DBConfig{
			Driver:  driver,
			Engine:  engine,
			Charset: charset,
		},
		OutFilePath: outFilePath,
	}

	dm, err := ddlmaker.NewMaker(conf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	structs := []interface{}{
		&ex.User{},
		ex.Entry{},
		ex.PlayerComment{},
		ex.Bookmark{},
	}

	dm.AddStruct(structs...)

	err = dm.Generate()
	if err != nil {
		log.Println(err.Error())
		return
	}
}
