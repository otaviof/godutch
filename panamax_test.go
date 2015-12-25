package godutch_test

import (
	"fmt"
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestOnboard(t *testing.T) {
	p := NewPanamax()
	c := mockBootstrappedContainer(t, "TestOnboard")

	Convey("Should be able to Onboard a Container", t, func() {
		err := p.Onboard(c)
		So(err, ShouldEqual, nil)
		c.Shutdown()
	})
}

func TestExecuteCheck(t *testing.T) {
	var err error
	var resp *Response
	var p *Panamax
	var c *Container

	p = NewPanamax()
	c = mockContainer(t, "TestExecuteCheck")

	p.RegisterService(c)
	go p.ServeBackground()

	Convey("Should be able to Onboard a Container.", t, func() {
		err = p.Onboard(c)
		So(err, ShouldEqual, nil)
	})

	for checkName, _ := range p.Checks {
		log.Println("TEST checkName:", checkName)
		Convey(fmt.Sprintf("Should be able to Execute '%s'", checkName), t, func() {
			resp, err = p.Execute(checkName, []string{})
			log.Printf("TEST Response: %#v", resp)
			So(err, ShouldEqual, nil)
			So(resp.Name, ShouldEqual, checkName)
		})
	}

	Convey("Should be able to offboard a container", t, func() {
		err = p.Offboard("TestExecuteCheck")
		So(err, ShouldEqual, nil)
	})

	p.Stop()

}

/* EOF */
