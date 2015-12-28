package godutch

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//
// Holds INI configuration file contents mapped into Config data struture.
//
type Config struct {
	GoDutch    GoDutchConfig
	NRPE       NRPEConfig
	NCSA       NCSAConfig
	CollectD   CollectDConfig
	Containers map[string]*ContainerConfig
}

//
// Auxiliary types to compose configuration, and to be imported and used
// further on the services implementation.
//

type GoDutchConfig struct {
	UseUnixSockets bool   `ini:"use_unix_sockets"`
	ContainersDir  string `ini:"containers_dir"`
	SocketsDir     string `ini:"sockets_dir"`
	TCPPortsRange  string `ini:"tcp_ports_range"`
}

type NRPEConfig struct {
	Enabled   bool   `ini:"enabled"`
	SSL       bool   `ini:"ssl"`
	Interface string `ini:"interface"`
	Port      int    `ini:"port"`
}

type NCSAConfig struct {
	Enabled bool   `ini:"enabled"`
	Address string `ini:"address"`
	Port    int    `ini:"port"`
}

type CollectDConfig struct {
	Enabled bool   `ini:"enabled"`
	Address string `ini:"address"`
	Port    int    `ini:"port"`
}

type ContainerConfig struct {
	Name    string   `ini:"name"`
	Enabled bool     `ini:"enabled"`
	Command []string `ini:"command"`
}

// Instantiate a new Config type, by loading informed configuration file and
// following Container's configuration directory, which might contain more
// files.
func NewConfig(configPath string) (*Config, error) {
	var err error
	var cfg *Config
	var cfgPathAbs string

	// definining absolute location for informed file
	cfgPathAbs, _ = filepath.Abs(configPath)

	// checking if informed config file indeed exists
	if _, err = exists(configPath); err != nil {
		log.Fatalln("Can't find configuration file at:", cfgPathAbs)
		return nil, err
	}

	// transforming INI file into local Config object
	if cfg, err = parseConfigINI(cfgPathAbs); err != nil {
		return nil, err
	}

	// verifying if socket directory exists
	if _, err = exists(cfg.GoDutch.SocketsDir); err != nil {
		log.Fatalln("Can't find directory for 'socket_dir': ", err)
		return nil, err
	}

	// loading containers configuration
	if err = cfg.globContainersConfig(filepath.Dir(cfgPathAbs)); err != nil {
		log.Fatalln("Error during containers load:", err)
		return nil, err
	}

	return cfg, nil
}

// Identifies absolute path for containers' directory and glob for INI files in
// there, composing a list of INI files on that directory.
func (cfg *Config) globContainersConfig(baseDir string) error {
	var err error
	var containersDir string
	var glob string
	var containers []string

	containersDir, _ = filepath.Abs(
		filepath.Join(baseDir, cfg.GoDutch.ContainersDir))

	log.Println("Containers directory at:", containersDir)
	if _, err = exists(containersDir); err != nil {
		log.Fatalln("Containers directory does not exist at:", containersDir)
		return err
	}

	// listing files on containers' config directory
	glob = fmt.Sprintf("%s/*.ini", containersDir)
	if containers, err = filepath.Glob(glob); err != nil {
		log.Fatalln("Errors on directory glob:", err)
		return err
	}

	// loading container configuration files
	log.Println("Containers config files:", containers)
	if len(containers) > 0 {
		cfg.loadContainersConfig(containers)
	}

	return nil
}

// Receives a list of Container INI configuration files and load them into
// "Containers" section of primary Config instance.
func (cfg *Config) loadContainersConfig(filePaths []string) error {
	var err error
	var iniCfg *ini.File
	var containerCfg *ContainerConfig
	var containerName string
	var filePath string

	// allocating memory for containers configuration
	cfg.Containers = make(map[string]*ContainerConfig)

	for _, filePath = range filePaths {
		// parsing container INI file
		if iniCfg, err = ini.Load(filePath); err != nil {
			log.Fatalln("Config file '", filePath, "' error:", err)
			return err
		}

		containerCfg = new(ContainerConfig)

		// by default a Container configuration file will hold by the section
		// also named "Container"
		if err = iniCfg.Section("Container").MapTo(containerCfg); err != nil {
			log.Fatalln("Error mapping Container INI:", err)
			return err
		}

		if len(containerCfg.Name) > 1 {
			if containerName, err = sanitizeName(containerCfg.Name); err != nil {
				log.Fatalln("Error on setting up container containerName:", err)
				return err
			}
			log.Println("Container Name:", containerName)
			cfg.Containers[containerName] = containerCfg
		} else {
			log.Println("[SKIP] No container name found on:", filePath)
		}
	}

	return nil
}

// Returns a sanitized name based on input raw input string. By a sanitized name
// it means only alpha-numeric cachacters, all lower.
func sanitizeName(rawName string) (string, error) {
	var err error
	var reg *regexp.Regexp
	var safe string

	if reg, err = regexp.Compile("[^A-Za-z0-9]+"); err != nil {
		log.Fatal("Error on compiling regexp:", err)
		return "", err
	}

	safe = reg.ReplaceAllString(rawName, "")
	safe = strings.ToLower(strings.Trim(safe, ""))

	return safe, nil
}

// Check if a file (or directory) exists.
func exists(path string) (bool, error) {
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

// Loads the primary configuration file and maps it into a Config type.
func parseConfigINI(cfgPathAbs string) (*Config, error) {
	var err error
	var cfg *Config = new(Config)
	var iniCfg *ini.File

	// loading the INI file contents into local struct
	if iniCfg, err = ini.Load(cfgPathAbs); err != nil {
		log.Fatalln("Errors on parsing INI:", err)
		return nil, err
	}

	// mapping configuration into local struct
	if err = iniCfg.MapTo(cfg); err != nil {
		log.Fatalln("Errors on mapping INI:", err)
		return nil, err
	}

	return cfg, nil
}

/* EOF */
