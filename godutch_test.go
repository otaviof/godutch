package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func mockGoDutch(t *testing.T) (*GoDutch) {
	var cfg *Config = mockNewConfig(t)
	var g *GoDutch
	var err error

	Convey("Should be able to instantiate GoDutch.", t, func () {
		g, err = NewGoDutch(cfg)
		So(err, ShouldEqual, nil)
	})

	return g
}

func TestLoadContainers(t *testing.T) {
	var g *GoDutch = mockGoDutch(t)
	var err error

	Convey("Should be able to load containers based on config.", t, func () {
		err = g.LoadContainers()
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
