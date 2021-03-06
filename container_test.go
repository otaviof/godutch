package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func mockContainer(t *testing.T) *Container {
	var err error
	var c *Container
	var cfg *Config = mockNewConfig(t)

	c, err = NewContainer(cfg.Container["rubycontainer"])

	Convey("Should not return errors on NewContainer", t, func() {
		So(err, ShouldEqual, nil)
	})

	return c
}

// Returns a bootstrapped container with the option of using a defer stop, when
// informed by parameter.
func mockBootstrappedContainer(t *testing.T, deferStop bool) *Container {
	var err error
	var c *Container = mockContainer(t)

	c.Client()
	go c.Bg.Serve()

	if deferStop {
		defer c.Bg.Stop()
	}

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
		Name:      "TestNewBgCmd",
		SocketDir: "/tmp",
		Command:   []string{"sleep", "1"},
	}

	Convey("Should not return errors on NewContainer", t, func() {
		_, err = NewContainer(containerCfg)
		So(err, ShouldEqual, nil)
	})
}

func TestBootstrapAndComponentChecks(t *testing.T) {
	var err error
	var req *Request
	var resp *Response
	var c *Container = mockBootstrappedContainer(t, false)

	Convey("Should be able to bootstrap a container", t, func() {
		So(
			strings.Join(c.Inventory(), "::"),
			ShouldContainSubstring,
			"check_",
		)
	})

	// This test expect to find '{"okay": 1}' as Metrics returned by
	// "check_test" on the ruby container
	Convey("Should be able to execute a exiting check", t, func() {
		req, _ = NewRequest("check_test", []string{})
		resp, err = c.Execute(req)
		So(err, ShouldEqual, nil)
		So(resp.Metrics[0], ShouldContainKey, "okay")
		So(resp.Metrics[0]["okay"], ShouldEqual, 1)
	})

	Convey("Should be able to shutdown container.", t, func() {
		err = c.Shutdown()
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
