package go_dutch

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config_path := "./etc/go-dutch/go-dutch.ini"

	Convey("Should be able to load a INI configuration file", t, func() {
		config := LoadConfig(config_path)
		So(config.GoDutch.ListenAddress, ShouldEqual, "127.0.0.1")
	})
}

/* EOF */
