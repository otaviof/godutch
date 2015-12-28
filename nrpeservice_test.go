package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestNewNRPEService(t *testing.T) {
	var err error
	var cfg *Config
	var ns *NRPEService
	var resp *Response

	p := NewPanamax()
	c := mockContainer(t, "TestNewNRPEService")

	cfg, _ = NewConfig("__etc/godutch/godutch.ini")
	ns, err = NewNRPEService(&cfg.NRPE, p)

	p.Add(ns)
	p.RegisterService(c)
	go p.ServeBackground()

	Convey("Should be able to Onboard a Container", t, func() {
		err = p.Onboard(c)
		So(err, ShouldEqual, nil)

		resp, err = p.Execute("check_test", []string{})
		So(err, ShouldEqual, nil)
		log.Printf("TEST Response: %#v", resp)
		log.Println("Checks: ", p.Checks)
	})

	p.Stop()
}

/* EOF */
