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
	containers map[string]*godutch.Container
	ns         *godutch.NRPESrvc
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
	var cCfg *godutch.ContainerConfig
	var c *godutch.Container

	self.containers = make(map[string]*godutch.Container)

	for name, cCfg = range self.cfg.Containers {
		log.Printf("Loading container: '%s'", name)

		if !cCfg.Enabled {
			log.Println("-- skipping disabled container --")
			continue
		}

		// spawn a new container
		if c, err = godutch.NewContainer(name, cCfg.Command); err != nil {
			log.Fatalln("NewContainer error:", err)
		}

		// keeping the container pointer for the onboard step
		self.containers[name] = c
		self.g.Register(c)
	}
}

// Loads the containers into GoDutch, by setting up socket communication and
// taking inventory of what are the available checks per container.
func (self *Self) onboardContainers() {
	var err error
	var name string
	var c *godutch.Container

	for name, c = range self.containers {
		log.Printf("Onboarding container: '%s'", name)
		if err = self.g.Onboard(c); err != nil {
			log.Fatalln("Error on onboarding:", err)
		}
	}
}

// Load the NRPE service interface.
func (self *Self) loadNRPEService() {
	var err error
	if self.ns, err = godutch.NewNRPESrvc(&self.cfg.NRPE, self.g); err != nil {
		log.Fatalln(err)
	}
	self.g.Register(self.ns)
}

// Add the NRPE service into the Supervisor, to start listening and executing
// checks, linked by the informed GoDutch pointer.
func (self *Self) onboardNRPEService() {
	var err error
	if err = self.g.Onboard(self.ns); err != nil {
		panic(err)
	}
}

//
// Main
//
func main() {
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
		// actions bellow suppose to happen after the Supervisor is loaded,
		// although, there's no confirmation of that state in place, yet.
		time.Sleep(1e9)
		self.onboardContainers()
		self.onboardNRPEService()
	}()

	self.g.Serve()
}

/* EOF */
