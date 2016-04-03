package main

import (
	"flag"
	"github.com/otaviof/godutch"
	"log"
)

func main() {
	var configFilePath string
	var cfg *godutch.Config
	var g *godutch.GoDutch
	var err error

	flag.StringVar(
		&configFilePath,
		"config-path",
		"/etc/godutch/godutch.ini",
		"Path to primary GoDutch configuration file.",
	)
	flag.Parse()

	if cfg, err = godutch.NewConfig(configFilePath); err != nil {
		log.Fatalln(err)
	}

	if g, err = godutch.NewGoDutch(cfg); err != nil {
		log.Fatalln(err)
	}

	if err = g.LoadContainers(); err != nil {
		log.Fatalln(err)
	}

	if err = g.LoadNrpeService(); err != nil {
		log.Fatalln(err)
	}

	g.Serve()
}

/* EOF */
