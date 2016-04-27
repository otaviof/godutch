package godutch_test

import (
	. "github.com/otaviof/godutch"
	gocache "github.com/patrickmn/go-cache"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func mockPanamax(t *testing.T) *Panamax {
	var p *Panamax
	var err error
	var cache *gocache.Cache

	cache = gocache.New(time.Minute, 20*time.Second)

	Convey("Should be able to instantiate Panamax", t, func() {
		p, err = NewPanamax(cache)
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

	// Just loading a container in Panamax, which has the expected check names
	// and returned values.
	Convey("Should load a container on Panamax", t, func() {
		err = p.Load(cfg.Container["rubycontainer"])
		So(err, ShouldEqual, nil)
	})

	// We assume here the ruby container will have "check_test" and
	// "check_second_test" methods, it's hardcoded here and in other tests
	Convey("Should Execute a Checks using the Panamax's routing", t, func() {
		for _, name = range []string{"check_test", "check_second_test"} {
			req, err = NewRequest(name, []string{})
			So(err, ShouldEqual, nil)

			resp, err = p.Execute(req)
			So(err, ShouldEqual, nil)
			So(resp.Name, ShouldEqual, req.Fields.Command)
		}
	})

	/// After running the checks, we can measure how old is the last run in
	/// seconds
	Convey("Should calculate when the check has last ran", t, func() {
		for _, name = range []string{"check_test", "check_second_test"} {
			So(p.CheckLastRun(name), ShouldBeGreaterThanOrEqualTo, 0)
		}
	})
}

/* EOF */
