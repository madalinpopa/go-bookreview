package app

import "flag"

// Config represents the configuration settings for application.
type Config struct {

	// Addr specifies the address for the service, which may include host and port information.
	Addr string

	// dsn specifies the database connection string used to configure the database connection for the application.
	Dsn string

	// port specifies the numerical network port used for the service to listen on or connect to.
	Port int

	// UploadDir specifies the directory path where uploaded files are stored for the application.
	UploadDir string
}

// NewConfig initializes and returns a pointer to a config struct populated with default CLI flags and values.
func NewConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.Addr, "addr", "localhost", "address to listen on")
	flag.IntVar(&config.Port, "port", 4000, "port to listen on")
	flag.StringVar(&config.Dsn, "dsn", "db.sqlite", "database connection string")
	flag.StringVar(&config.UploadDir, "upload-dir", "uploads", "directory for uploaded files")
	flag.Parse()

	return config
}
