package godutch_test

import (
	"fmt"
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
	var sc *ServiceConfig
	var dialOn []string
	var dialStr string
	var host string
	var port int
	var i int

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
			"bin")
	})

	Convey("Should be able to detect NSCA configuration", t, func() {
		So(cfg.Service["nscaservice"].Type, ShouldEqual, "nsca")
		So(cfg.Service["nscaservice"].Port, ShouldEqual, 0)
		So(cfg.Service["nscaservice"].LastRunThreshold, ShouldBeGreaterThan, 0)
	})

	Convey("Should be able to parse Carbon's 'dial_on' string", t, func() {
		sc = cfg.Service["carbonrelay"]
		So(sc.Type, ShouldEqual, "carbon")
		dialOn = sc.ParseDialOn()
		So(len(dialOn), ShouldEqual, 2)

		for i, dialStr = range dialOn {
			// adapting position counter to match mock configuration
			i += 1
			// parsing into host (string) and port (int)
			host, port = sc.ParseDialString(dialStr)

			So(host, ShouldEqual, fmt.Sprintf("null%d.local", i))
			So(port, ShouldEqual, 2003)
		}
	})
}

/* EOF */
