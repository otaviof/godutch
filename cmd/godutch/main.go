/*

Creates a GoDutch daemon based on command-line informed configuration. Good
proof of concept for the application server elements, like auto-healing when a
Container dies, and a very interesting performance, intiially.

Usage:

 $ go build -a -o bin/godutch cmd/godutch
 $ ./bin/godutch -config-path __etc/godutch/godutch.ini

*/
package main

import (
	"flag"
	"github.com/otaviof/godutch"
	"log"
	"time"
)

//
// Represents this command-line utility, with the necessary objects and data to
// boostrap a GoDutch Server.
//
type Self struct {
	cfgPath    string
	cfg        *godutch.Config
	g          *godutch.GoDutch
	services   map[string]*godutch.Service
	containers map[string]*godutch.Container
}

// Reads the configuration file and loads into itself, base step.
func (self *Self) loadConfig() {
	var err error
	if self.cfg, err = godutch.NewConfig(self.cfgPath); err != nil {
		panic(err)
	}
}

// Loads GoDutch.
func (self *Self) loadGoDutch() {
	self.g = godutch.NewGoDutch()
}

// Using configuration will load the containers and register them within the
// Supervisor. Containers pointers are kept for Onboarding step.
func (self *Self) loadContainers() {
	var err error
	var name string
	var containerCfg *godutch.ContainerConfig
	var c *godutch.Container

	self.containers = make(map[string]*godutch.Container)

	for name, containerCfg = range self.cfg.Containers {
		log.Printf("-- Loading container: '%s'", name)
		log.Println("Container Cfg:", containerCfg)

		if !containerCfg.Enabled {
			log.Println("--!!-- SKIPPING DISABLED CONTAINER --!!--")
			continue
		}

		// spawn a new container
		if c, err = godutch.NewContainer(containerCfg); err != nil {
			log.Fatalln("NewContainer error:", err)
		}

		// keeping the container pointer for the onboard step
		self.containers[name] = c

		if err = self.g.Register(c); err != nil {
			panic(err)
		}
	}
}

// Loads the containers into GoDutch, by setting up socket communication and
// taking inventory of what are the available checks per container.
func (self *Self) onboardContainers() {
	var err error
	var name string
	var c *godutch.Container

	for name, c = range self.containers {
		log.Printf("-- Onboarding container: '%s'", name)
		if err = c.Bootstrap(); err != nil {
			log.Fatalln("On bootstrap, error:", err)
			panic(err)
		}

		if err = self.g.Onboard(c); err != nil {
			log.Fatalln("Error on onboarding:", err)
			panic(err)
		}
	}
}

// Load the NRPE service interface.
func (self *Self) loadServices() {
	var err error
	var name string
	var serviceCfg *godutch.ServiceConfig
	var s *godutch.Service

	self.services = make(map[string]*godutch.Service)

	for name, serviceCfg = range self.cfg.Services {
		log.Println("Service: ", name)

		s = godutch.NewService(serviceCfg, self.g)
		self.services[name] = s

		if err = self.g.Register(s); err != nil {
			panic(err)
		}
	}
}

// Add the NRPE service into the Supervisor, to start listening and executing
// checks, linked by the informed GoDutch pointer.
func (self *Self) onboardServices() {
	var err error
	var s *godutch.Service

	for _, s = range self.services {
		if err = self.g.Onboard(s); err != nil {
			log.Fatalln("Error on onboarding:", err)
		}
	}
}

//
// Main
//
func main() {
	var configFilePath string
	var self Self

	flag.StringVar(
		&configFilePath,
		"config-path",
		"/etc/godutch/godutch.ini",
		"Path to configuration file, `godutch.ini`")
	flag.Parse()

	self = Self{cfgPath: configFilePath}

	self.loadConfig()
	self.loadGoDutch()
	self.loadServices()
	self.loadContainers()

	go func(self *Self) {
		// actions bellow suppose to happen after the Supervisor is loaded,
		// although, there's no confirmation of that state in place, yet.
		time.Sleep(1e9)
		self.onboardContainers()
		self.onboardServices()
	}(&self)

	self.g.Serve()
}

/* EOF */
