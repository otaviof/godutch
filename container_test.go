package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestNewContainer(t *testing.T) {
	Convey("Should not return errors on NewContainer", t, func() {
		_, err := NewContainer("TestNewContainer", []string{"sleep", "1"})
		So(err, ShouldEqual, nil)
	})
}

func TestBootstrap(t *testing.T) {
	c, _ := NewContainer(
		"TestBootstrap",
		[]string{
			"/usr/bin/ruby",
			"/home/otaviof/src/go/tmp/starlite/godutch_test.rb"},
	)

	go c.Bg.Serve()
	time.Sleep(1e9)

	Convey("Should be able to bootstrap a container", t, func() {
		err := c.Bootstrap()
		So(err, ShouldEqual, nil)
		c.Bg.Stop()
	})
}

/* EOF */
