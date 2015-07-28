package main

import "fmt"
import "github.com/otaviof/godutch/config"

func main() {
	config_path := "../etc/godutch/godutch.ini"
	cfg := config.LoadConfig(config_path)
	fmt.Println(cfg)
}

/* EOF */
