package go_dutch

import (
	"code.google.com/p/gcfg"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GoDutch struct {
		ListenAddress   string
		ListenPort      int
		ChecksDirectory string
	}

	Ncsa struct {
		Enabled       bool
		ServerAddress string
		ServerPort    int
	}
}

// LoadConfig: Loads a INI formatted configuration file that will contain
// primary Go-Dutch settings. And it returns a struct typed Config.
func LoadConfig(path string) Config {
	// transforming informed path on the full absolute representation
	config_file_path, _ := filepath.Abs(path)
	// absolute path must exists on file-system
	if _, err := os.Stat(config_file_path); os.IsNotExist(err) {
		panic(fmt.Sprintf(
			"Can't find config file at: '%s'", config_file_path))
	}

	// loading INI file using gcfg, which by design add validation where it
	// must find the values described on the struct Config, and if entries
	// some are missing it will also trigger the panic
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, config_file_path)
	if err != nil {
		error_msg := fmt.Sprintf(
			"Can't load Go-Dutch configuration from: '%s'", err,
		)
		panic(error_msg)
	}

	return cfg
}

/* EOF */
