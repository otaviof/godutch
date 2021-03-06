package godutch_test

import (
	"fmt"
	. "github.com/otaviof/godutch"
	"github.com/otaviof/gonrpe"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
	"time"
)

func TestNewNrpeService(t *testing.T) {
	var cfg *Config = mockNewConfig(t)
	var p *Panamax = mockPanamax(t)
	var listenOn string = fmt.Sprintf(
		"%s:%d",
		cfg.Service["nrpeservice"].Interface,
		cfg.Service["nrpeservice"].Port)
	var ns *NrpeService
	var conn net.Conn
	var wroteLen int
	var err error

	ns = NewNrpeService(cfg.Service["nrpeservice"], p)

	go ns.Serve()
	defer ns.Stop()
	time.Sleep(1e9)

	Convey("Should be able to Onboard a Service", t, func() {
		conn, err = net.Dial("tcp", listenOn)
		So(err, ShouldEqual, nil)

		wroteLen, err = conn.Write(gonrpe.SAMPLE_PACKET_NRPE_PAYLOAD)
		So(wroteLen, ShouldEqual, gonrpe.NRPE_PACKET_SIZE)
		So(err, ShouldEqual, nil)

		err = conn.Close()
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
