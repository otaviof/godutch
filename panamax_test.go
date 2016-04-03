package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func mockPanamax(t *testing.T) *Panamax {
	var p *Panamax
	var err error

	Convey("Should be able to instantiate Panamax", t, func() {
		p, err = NewPanamax()
		So(err, ShouldEqual, nil)
	})

	return p
}

func TestLoadAndExecute(t *testing.T) {
	var p *Panamax = mockPanamax(t)
	var cfg *Config = mockNewConfig(t)
	var req *Request
	var resp *Response
	var name string
	var err error

	Convey("Should load a container", t, func() {
		err = p.Load(cfg.Containers["rubycontainer"])
		So(err, ShouldEqual, nil)
	})

	Convey("Should Execute a Checks using the Panamax's routing", t, func() {
		for _, name = range []string{"check_test", "check_second_test"} {
			req, err = NewRequest(name, []string{})
			So(err, ShouldEqual, nil)

			resp, err = p.Execute(req)
			So(err, ShouldEqual, nil)
			So(resp.Name, ShouldEqual, req.Fields.Command)
		}
	})
}

/* EOF */
