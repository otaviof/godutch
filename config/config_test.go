package config_test

import (
	. "github.com/smartystreets/goconvey/convey"
    "github.com/otaviof/godutch/config"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg := config.LoadConfig("../__etc/godutch/godutch.ini")

	Convey("Should be able to read a String", t, func() {
		So(cfg.GoDutch.ListenAddress, ShouldEqual, "0.0.0.0")
	})

	Convey("Should be able to read a Integer", t, func() {
		So(cfg.GoDutch.ListenPort, ShouldEqual, 5666)
	})
}

/* EOF */
