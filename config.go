package ddlmaker

// Config set user enviroment
type Config struct {
	OutFilePath string
	DB          DBConfig
}

// DBConfig set user db enviroment
type DBConfig struct {
	Driver  string
	Engine  string
	Charset string
}
