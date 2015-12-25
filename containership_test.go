package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
	"time"
)

func TestOnboard(t *testing.T) {
	cs := NewContainerShip()
	c := mockBootstrappedContainer(t, "TestOnboard")

	Convey("Should be able to Onboard a Container", t, func() {
		err := cs.Onboard(c)
		So(err, ShouldEqual, nil)
		c.Shutdown()
	})
}

func TestExecuteCheck(t *testing.T) {
	cs := NewContainerShip()
	c := mockContainer(t, "TestExecuteCheck")

	cs.Add(c.Bg)
	go cs.ServeBackground()
	time.Sleep(1e9)

	Convey("Should be able to bootstrap a container", t, func() {
		err1 := c.Bootstrap()
		So(err1, ShouldEqual, nil)
		err2 := cs.Onboard(c)
		So(err2, ShouldEqual, nil)
	})

	Convey("Should be able to invoke a check", t, func() {
		time.Sleep(1e9)
		resp, err := cs.Execute("check_test", []string{})
		log.Println("TEST Response:", resp)
		So(err, ShouldEqual, nil)
		So(resp.Name, ShouldEqual, "check_test")
	})

	cs.Stop()

}

/* EOF */
