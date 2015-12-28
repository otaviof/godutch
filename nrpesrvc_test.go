package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewNRPESrvc(t *testing.T) {
	var err error
	var cfg *Config
	var g *GoDutch
	var ns *NRPESrvc

	g = NewGoDutch()
	cfg, _ = NewConfig("__etc/godutch/godutch.ini")
	ns, err = NewNRPESrvc(&cfg.NRPE, g)

	g.Register(ns)
	go g.ServeBackground()

	Convey("Should be able to Onboard a Container", t, func() {
		err = g.Onboard(ns)
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
