package godutch_test

import (
	"fmt"
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
	"time"
)

func TestNewNrpeService(t *testing.T) {
	var err error
	var cfg *Config
	var g *GoDutch
	var s *Service
	var conn net.Conn
	var listenOn string
	var n int

	cfg, _ = NewConfig("__etc/godutch/godutch.ini")

	listenOn = fmt.Sprintf(
		"%s:%d",
		cfg.Services["nrpeservice"].Interface,
		cfg.Services["nrpeservice"].Port)

	g = NewGoDutch()
	s = NewService(cfg.Services["nrpeservice"], g)
	g.Register(s)

	go g.ServeBackground()

	Convey("Should be able to Onboard a Service", t, func() {
		err = g.Onboard(s)
		So(err, ShouldEqual, nil)

		time.Sleep(1e9)

		conn, err = net.Dial("tcp", listenOn)
		So(err, ShouldEqual, nil)

		n, err = conn.Write(CHECK_NRPE_PAYLOAD)
		So(n, ShouldEqual, NRPE_PACKET_SIZE)
		So(err, ShouldEqual, nil)

		err = conn.Close()
		So(err, ShouldEqual, nil)
	})

	Convey("Should be able to shutdown a service", t, func() {
		err = s.Shutdown()
		So(err, ShouldEqual, nil)
	})

	defer g.Stop()
}

/* EOF */
