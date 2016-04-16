package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func mockNewConfig(t *testing.T) *Config {
	var err error
	var cfg *Config

	cfg, err = NewConfig("test/etc/godutch.ini")

	Convey("Should be able to load config without errors", t, func() {
		So(err, ShouldEqual, nil)
	})

	return cfg
}

func TestNewConfig(t *testing.T) {
	var cfg *Config = mockNewConfig(t)

	Convey("Should be able to read a String", t, func() {
		So(cfg.GoDutch.UseUnixSockets, ShouldEqual, true)
	})

	Convey("Should be able to read a Integer", t, func() {
		So(cfg.Service["nrpeservice"].Port, ShouldEqual, 5666)
	})

	Convey("Should be able to load example containers", t, func() {
		So(len(cfg.Container), ShouldBeGreaterThan, 0)
		So(cfg.Container["rubycontainer"].Command[0],
			ShouldEqual,
			"/usr/bin/ruby")
		So(cfg.Container["perlcontainer"].Command[0],
			ShouldContainSubstring,
			"bin/godutch")
	})

	Convey("Should be able to detect NCSA configuration", t, func() {
		So(cfg.Service["ncsaservice"].Type, ShouldEqual, "NCSA")
	})
}

/* EOF */
