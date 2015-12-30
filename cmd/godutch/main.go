package main

import (
	"flag"
	"github.com/otaviof/godutch"
	"log"
	"time"
)

type Self struct {
	cfgPath    string
	cfg        *godutch.Config
	g          *godutch.GoDutch
	containers map[string]*godutch.Container
	ns         *godutch.NRPESrvc
}

func (self *Self) loadConfig() {
	var err error
	if self.cfg, err = godutch.NewConfig(self.cfgPath); err != nil {
		panic(err)
	}
}

func (self *Self) loadGoDutch() {
	self.g = godutch.NewGoDutch()
}

func (self *Self) loadContainers() {
	var err error
	var name string
	var cCfg *godutch.ContainerConfig
	var c *godutch.Container

	self.containers = make(map[string]*godutch.Container)

	for name, cCfg = range self.cfg.Containers {
		log.Printf("Loading container: '%s'", name)
		if !cCfg.Enabled {
			log.Println("-- skipping disabled container --")
			continue
		}

		if c, err = godutch.NewContainer(name, cCfg.Command); err != nil {
			log.Fatalln("NewContainer error:", err)
		}

		// keeping the container pointer for the onboard process
		self.containers[name] = c
		self.g.Register(c)
	}
}

func (self *Self) onboardContainers() {
	var err error
	var name string
	var c *godutch.Container

	for name, c = range self.containers {
		log.Printf("Onboarding container: '%s'", name)
		if err = self.g.Onboard(c); err != nil {
			log.Println("Error on onboarding:", err)
		}
	}
}

func (self *Self) loadNRPEService() {
	var err error
	if self.ns, err = godutch.NewNRPESrvc(&self.cfg.NRPE, self.g); err != nil {
		log.Fatalln(err)
	}
	self.g.Register(self.ns)
}

func (self *Self) onboardNRPEService() {
	var err error
	if err = self.g.Onboard(self.ns); err != nil {
		panic(err)
	}
}

func main() {
	// var err error
	var configFilePath string
	var self *Self

	flag.StringVar(
		&configFilePath,
		"config-path",
		"/etc/godutch/godutch.ini",
		"Path to configuration file, `godutch.ini`")
	flag.Parse()

	self = &Self{cfgPath: configFilePath}

	self.loadConfig()
	self.loadGoDutch()
	self.loadNRPEService()
	self.loadContainers()
	go func() {
		time.Sleep(1e9)
		self.onboardContainers()
		self.onboardNRPEService()
	}()
	self.g.Serve()
}

/* EOF */
