package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewService(t *testing.T) {
	var err error
	var cfg *Config
	var g *GoDutch = NewGoDutch()
	var s *Service

	cfg, _ = NewConfig("__etc/godutch/godutch.ini")
	s = NewService(cfg.Services["nrpeservice"], g)

	Convey("Should be able to register a Service instance", t, func() {
		err = g.Register(s)
		So(err, ShouldEqual, nil)
	})

	Convey("Should be able to find Execute as a dummy method", t, func() {
		_, err = s.Execute([]byte("whatever"))
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
