package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func mockContainer(t *testing.T, name string) *Container {
	var err error
	var c *Container
	var cwd string
	var command []string
	var testScriptPath string

	// determining the script location, based on current
	cwd, _ = os.Getwd()
	testScriptPath = path.Join(cwd, "godutch-checks.rb")
	command = []string{"/usr/bin/ruby", testScriptPath}

	c, err = NewContainer(name, command)
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
		time.Sleep(1e9)
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
		So(
			strings.Join(c.Checks, "::"),
			ShouldEqual,
			strings.Join([]string{"check_test", "check_second_test"}, "::"),
		)
	})

	Convey("Should be able to shutdown a running Container", t, func() {
		var err error
		err = c.Shutdown()
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
