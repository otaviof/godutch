package godutch

//
// Config handles INI files spread over a few directories. The contents will be
// parsed into specific structs, regarding containers or services.
//

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//
// Holds INI configuration file contents mapped into Config data struture.
//
type Config struct {
	GoDutch   GoDutchConfig
	Service   map[string]*ServiceConfig
	Container map[string]*ContainerConfig
}

//
// Auxiliary types to compose configuration, and to be imported and used
// further on the services implementation.
//

type GoDutchConfig struct {
	UseUnixSockets bool   `ini:"use_unix_sockets"`
	ContainersDir  string `ini:"containers_dir"`
	ServicesDir    string `ini:"services_dir"`
	TCPPortsRange  string `ini:"tcp_ports_range"`
}

type ContainerConfig struct {
	Enabled   bool     `ini:"enabled"`
	Name      string   `ini:"name"`
	Command   []string `ini:"command"`
	SocketDir string   `ini:"socket_dir"`
}

type ServiceConfig struct {
	Enabled          bool   `ini:"enabled"`
	Type             string `ini:"type"`
	Name             string `ini:"name"`
	Interface        string `ini:"interface"`
	Port             int    `ini:"port"`
	DialOn           string `ini:"dial_on"`
	Ssl              bool   `ini:"ssl"`
	LastRunThreshold int64  `ini:"last_run_threshold"`
}

// Instantiate a new Config type, by loading informed configuration file and
// following Container's configuration directory, which might contain more
// files.
func NewConfig(configPath string) (*Config, error) {
	var err error
	var cfg *Config
	var cfgPathAbs string
	var dirPath string

	// definining absolute location for informed file
	cfgPathAbs, _ = filepath.Abs(configPath)

	// checking if informed config file indeed exists
	if _, err = exists(configPath); err != nil {
		log.Printf("[Config] Can't find config file at: '%s'", cfgPathAbs)
		return nil, err
	}

	// transforming INI file into local Config object
	if cfg, err = parseConfigINI(cfgPathAbs); err != nil {
		return nil, err
	}

	cfg.Service = make(map[string]*ServiceConfig)
	cfg.Container = make(map[string]*ContainerConfig)

	for _, dirPath = range []string{
		cfg.GoDutch.ServicesDir,
		cfg.GoDutch.ContainersDir,
	} {
		// loading containers configuration
		if err = cfg.globIniConfigFIles(
			filepath.Dir(cfgPathAbs),
			dirPath,
		); err != nil {
			log.Println("[Config] Error on loading config file: ", err)
			return nil, err
		}
	}

	return cfg, nil
}

// Identifies absolute path for containers' directory and glob for INI files in
// there, composing a list of INI files on that directory.
func (cfg *Config) globIniConfigFIles(baseDir string, cfgDir string) error {
	var err error
	var cfgDirAbs string
	var glob string
	var cfgPaths []string

	cfgDirAbs, _ = filepath.Abs(filepath.Join(baseDir, cfgDir))

	if _, err = exists(cfgDirAbs); err != nil {
		log.Printf("[Config] Containers dir not found at: '%s'", cfgDirAbs)
		return err
	}

	// listing files on containers' config directory
	glob = fmt.Sprintf("%s/*.ini", cfgDirAbs)
	if cfgPaths, err = filepath.Glob(glob); err != nil {
		log.Println("[Config] Errors on directory glob:", err)
		return err
	}

	// loading container configuration files
	if len(cfgPaths) > 0 {
		if err = cfg.loadIniConfigs(cfgPaths); err != nil {
			log.Println("[Config] During config-file load:", err)
			return err
		}
	}

	return nil
}

// Receives a list of Container INI configuration files and load them into
// "Container" section of primary Config instance.
func (cfg *Config) loadIniConfigs(cfgPaths []string) error {
	var err error
	var iniCfg *ini.File
	var section *ini.Section
	var sectionName string
	var cfgPath string
	var serviceCfg *ServiceConfig
	var containerCfg *ContainerConfig
	var name string
	var match bool

	for _, cfgPath = range cfgPaths {
		log.Printf("[Config] Loading: '%s'", cfgPath)

		// avoiding dummy files
		if match, _ = path.Match("\\.\\#*\\.ini", path.Base(cfgPath)); match {
			log.Printf("[Config] Ignoring config file: '%s'", cfgPath)
			continue
		}

		if iniCfg, err = ini.Load(cfgPath); err != nil {
			log.Println("[Config] Config load error:", err)
			return err
		}

		// the section names will determine witch kind of configuratio this is
		for _, sectionName = range iniCfg.SectionStrings() {
			section = iniCfg.Section(sectionName)

			switch sectionName {
			case "Service":
				serviceCfg = new(ServiceConfig)

				if err = section.MapTo(serviceCfg); err != nil {
					log.Println("[Config] Error on mapTo ServiceConfig:", err)
					return err
				}

				if name, err = sanitizeName(serviceCfg.Name); err != nil {
					log.Println("[Config] Error on sanitize name:", err)
					return err
				}

				log.Printf("[Config] Adding service: '%s'", name)
				serviceCfg.Name = name
				cfg.Service[name] = serviceCfg

				log.Printf("[Config] DEBUG serviceCfg: '%+v'", serviceCfg)
			case "Container":
				containerCfg = new(ContainerConfig)

				if err = section.MapTo(containerCfg); err != nil {
					log.Println("[Config] Error on mapTo Container:", err)
					return err
				}

				if name, err = sanitizeName(containerCfg.Name); err != nil {
					log.Println("[Config] Error on sanitize name:", err)
					return err
				}

				log.Printf("[Config] Adding container: '%s'", name)
				containerCfg.Name = name
				cfg.Container[name] = containerCfg

				log.Printf("[Config] DEBUG containerCfg: '%+v'", containerCfg)
			case "DEFAULT":
				continue
			default:
				log.Printf("[Config] Ignored section: '%s'", sectionName)
			}
		}
	}

	return nil
}

// Parses the "dial_on" string present on services that need to dial a external
// network communication.
func (sc *ServiceConfig) ParseDialOn() []string {
	return strings.Split(sc.DialOn, ", ")
}

func (sc *ServiceConfig) ParseDialString(dialStr string) (string, int) {
	var str []string = strings.Split(dialStr, ":")
	var host string = str[0]
	var port string = str[1]
	var portInt int
	portInt, _ = strconv.Atoi(port)
	return host, portInt
}

// Returns a sanitized name based on input raw input string. By a sanitized name
// it means only alpha-numeric cachacters, all lower.
func sanitizeName(rawName string) (string, error) {
	var err error
	var reg *regexp.Regexp
	var safe string

	if reg, err = regexp.Compile("[^A-Za-z0-9]+"); err != nil {
		return "", err
	}

	safe = reg.ReplaceAllString(rawName, "")
	safe = strings.ToLower(strings.Trim(safe, ""))

	if len(safe) <= 1 {
		err = errors.New("Result string is too short.")
		return "", err
	}

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
		log.Fatalln("[Config] Errors on parsing INI file:", err)
		return nil, err
	}

	// mapping configuration into local struct
	if err = iniCfg.MapTo(cfg); err != nil {
		log.Fatalln("[Config] Errors on mapping INI:", err)
		return nil, err
	}

	return cfg, nil
}

/* EOF */
