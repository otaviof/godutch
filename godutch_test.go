package godutch_test

import (
	"fmt"
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
	"time"
)

// Test the execution of every check known, from GoDutch Execute method, which
// calls other method down the stack.
func TestExecuteChecks(t *testing.T) {
	var err error
	var g *GoDutch
	var c *Container
	var component *Component
	var resp *Response

	g = NewGoDutch()
	c = mockContainer(t)
	component = c.ComponentInfo()

	g.Register(c)
	go g.ServeBackground()

	c.Bootstrap()

	Convey("Should be able to Onboard a Container.", t, func() {
		err = g.Onboard(c)
		So(err, ShouldEqual, nil)
	})

	// caling every known check, making sure there's response
	for _, checkName := range component.Checks {
		log.Println("TEST checkName:", checkName)
		Convey(fmt.Sprintf("Should be able to Execute '%s'", checkName), t, func() {
			resp, err = g.Execute(checkName, []string{})
			log.Printf("TEST Response: %#v", resp)
			So(err, ShouldEqual, nil)
			So(resp.Name, ShouldEqual, checkName)
		})
	}

	Convey("Should be able to offboard a container", t, func() {
		time.Sleep(1e9)
		err = g.Offboard(component.Name)
		So(err, ShouldEqual, nil)
	})

	defer g.Stop()
}

/* EOF */
