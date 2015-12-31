package godutch_test

import (
	"fmt"
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

func TestNewBgCmd(t *testing.T) {
	var bg *BgCmd
	var bgCmdName string = "TestNewBgCmd"
	var socketEnvStr string = fmt.Sprintf("godutch-%s.sock", bgCmdName)
	var containerCfg *ContainerConfig = &ContainerConfig{
		Name:    bgCmdName,
		Command: []string{"sleep", "1"},
	}

	bg = NewBgCmd(containerCfg)

	Convey("Should be albe have a custom ENV.", t, func() {
		So(bg.Name, ShouldEqual, bgCmdName)
		So(strings.Join(bg.Env, ";"), ShouldContainSubstring, socketEnvStr)
	})
}

/* EOF */
