package godutch_test

import (
	"fmt"
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestOnboard(t *testing.T) {
	var err error
	var g *GoDutch
	var c *Container

	g = NewGoDutch()
	c = mockBootstrappedContainer(t, "TestOnboard")

	Convey("Should be able to Onboard a Container", t, func() {
		err = g.Onboard(c)
		So(err, ShouldEqual, nil)
		c.Shutdown()
	})
}

func TestExecuteCheck(t *testing.T) {
	var err error
	var resp *Response
	var g *GoDutch
	var c *Container

	g = NewGoDutch()
	c = mockContainer(t, "TestExecuteCheck")

	g.Register(c)
	go g.ServeBackground()

	Convey("Should be able to Onboard a Container.", t, func() {
		err = g.Onboard(c)
		So(err, ShouldEqual, nil)
	})

	for _, checkName := range c.ComponentChecks() {
		log.Println("TEST checkName:", checkName)
		Convey(fmt.Sprintf("Should be able to Execute '%s'", checkName), t, func() {
			resp, err = g.Execute(checkName, []string{})
			log.Printf("TEST Response: %#v", resp)
			So(err, ShouldEqual, nil)
			So(resp.Name, ShouldEqual, checkName)
		})
	}

	Convey("Should be able to offboard a container", t, func() {
		err = g.Offboard("TestExecuteCheck")
		So(err, ShouldEqual, nil)
	})

	g.Stop()

}

/* EOF */
