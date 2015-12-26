package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewConfig(t *testing.T) {
	var err error
	var cfg *Config

	cfg, err = NewConfig("__etc/godutch/godutch.ini")

	Convey("Should be able to load config without errors", t, func() {
		So(err, ShouldEqual, nil)
	})

	Convey("Should be able to read a String", t, func() {
		So(cfg.GoDutch.UseUnixSockets, ShouldEqual, true)
	})

	Convey("Should be able to read a Integer", t, func() {
		So(cfg.NRPE.Port, ShouldEqual, 5666)
	})

	Convey("Should be able to load example containers", t, func() {
		So(len(cfg.Containers), ShouldBeGreaterThan, 0)
		So(
			cfg.Containers["rubycontainer"].Command[0],
			ShouldEqual,
			"/usr/bin/ruby",
		)
	})
}

/* EOF */
