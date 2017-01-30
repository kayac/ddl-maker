package ddlmaker

// Config set user environment
type Config struct {
	OutFilePath string
	DB          DBConfig
}

// DBConfig set user db environment
type DBConfig struct {
	Driver  string
	Engine  string
	Charset string
}
