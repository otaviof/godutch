package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"os/exec"
	"strings"
	"testing"
)

func TestNewBgCmd(t *testing.T) {
	bg := NewBgCmd("TestNewBgCmd", exec.Command("sleep", "1"))

	Convey("Should be albe have a new obj. with custom Env.", t, func() {
		So(bg.Name, ShouldEqual, "TestNewBgCmd")
		socketStr := "GODUTCH_SOCKET_PATH=/tmp/godutch-TestNewBgCmd.sock"
		So(strings.Join(bg.Env, ";"), ShouldContainSubstring, socketStr)
	})
}

/* EOF */
