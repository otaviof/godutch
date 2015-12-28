package godutch_test

import (
	"fmt"
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
)

func TestNewNRPESrvc(t *testing.T) {
	var err error
	var cfg *Config
	var g *GoDutch
	var ns *NRPESrvc
	var conn net.Conn
	var listenOn string
	var n int

	g = NewGoDutch()
	cfg, _ = NewConfig("__etc/godutch/godutch.ini")
	// same string used on service object
	listenOn = fmt.Sprintf("%s:%d", cfg.NRPE.Interface, cfg.NRPE.Port)
	ns, err = NewNRPESrvc(&cfg.NRPE, g)

	g.Register(ns)
	go g.ServeBackground()

	Convey("Should be able to Onboard a Container", t, func() {
		err = g.Onboard(ns)
		So(err, ShouldEqual, nil)
	})

	Convey("Should be have the port open, and accepting payload", t, func() {
		conn, err = net.Dial("tcp", listenOn)
		So(err, ShouldEqual, nil)
		n, err = conn.Write(CHECK_NRPE_PAYLOAD)
		So(n, ShouldEqual, NRPE_PACKET_SIZE)
		So(err, ShouldEqual, nil)
		err = conn.Close()
		So(err, ShouldEqual, nil)
	})

	g.Stop()
}

/* EOF */
