package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	// "os"
	// "path"
	"strings"
	"testing"
)

func mockContainer(t *testing.T, name string) *Container {
	var err error
	var c *Container
	var cfg *Config

	if cfg, err = NewConfig("__etc/godutch/godutch.ini"); err != nil {
		panic(err)
	}

	c, err = NewContainer(name, cfg.Containers["rubycontainer"].Command)

	Convey("Should not return errors on NewContainer", t, func() {
		So(err, ShouldEqual, nil)
	})

	return c
}

func mockBootstrappedContainer(t *testing.T, name string) *Container {
	var err error
	var c *Container = mockContainer(t, name)

	go c.Bg.Serve()

	Convey("Should be able to bootstrap a container", t, func() {
		err = c.Bootstrap()
		So(err, ShouldEqual, nil)
	})

	return c
}

func TestNewContainer(t *testing.T) {
	Convey("Should not return errors on NewContainer", t, func() {
		_, err := NewContainer("TestNewContainer", []string{"sleep", "1"})
		So(err, ShouldEqual, nil)
	})
}

func TestBootstrapAndShutdown(t *testing.T) {
	c := mockBootstrappedContainer(t, "TestBootstrapAndShutdown")

	Convey("Should be able to bootstrap a container", t, func() {
		So(strings.Join(c.Checks, "::"), ShouldContainSubstring, "check")
	})

	Convey("Should be able to shutdown a running Container", t, func() {
		var err error
		err = c.Shutdown()
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
