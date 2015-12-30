package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestNewBgCmd(t *testing.T) {
	var bg *BgCmd
	var containerCfg *ContainerConfig = &ContainerConfig{
		Name:    "TestNewBgCmd",
		Command: []string{"sleep", "1"},
	}

	bg = NewBgCmd(containerCfg)

	Convey("Should be albe have a custom ENV.", t, func() {
		So(bg.Name, ShouldEqual, "TestNewBgCmd")
		socketStr := "GODUTCH_SOCKET_PATH=/tmp/godutch-TestNewBgCmd.sock"
		So(strings.Join(bg.Env, ";"), ShouldContainSubstring, socketStr)
	})
}

/* EOF */
