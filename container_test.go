package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	// "os"
	// "path"
	"strings"
	"testing"
	"time"
)

func mockContainer(t *testing.T) *Container {
	var err error
	var c *Container
	var cfg *Config

	if cfg, err = NewConfig("__etc/godutch/godutch.ini"); err != nil {
		panic(err)
	}

	c, err = NewContainer(cfg.Containers["rubycontainer"])

	Convey("Should not return errors on NewContainer", t, func() {
		So(err, ShouldEqual, nil)
	})

	return c
}

func mockBootstrappedContainer(t *testing.T) *Container {
	var err error
	var c *Container = mockContainer(t)

	go c.Bg.Serve()
	defer c.Bg.Stop()

	Convey("Should be able to bootstrap a container", t, func() {
		time.Sleep(1e9)
		err = c.Bootstrap()
		So(err, ShouldEqual, nil)
	})

	return c
}

func TestNewContainer(t *testing.T) {
	var err error
	var containerCfg *ContainerConfig = &ContainerConfig{
		Name:    "TestNewBgCmd",
		Command: []string{"sleep", "1"},
	}

	Convey("Should not return errors on NewContainer", t, func() {
		_, err = NewContainer(containerCfg)
		So(err, ShouldEqual, nil)
	})
}

func TestBootstrapAndComponentChecks(t *testing.T) {
	var err error
	var c *Container = mockBootstrappedContainer(t)
	var component *Component = c.ComponentInfo()

	Convey("Should be able to bootstrap a container", t, func() {
		So(
			strings.Join(component.Checks, "::"),
			ShouldContainSubstring,
			"check",
		)
	})

	Convey("Should be able to shutdown container.", t, func() {
		err = c.Shutdown()
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
