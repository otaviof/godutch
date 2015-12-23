package godutch

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type BgCmd struct {
	Name       string
	SocketPath string
	Cmd        *exec.Cmd
	Env        []string
	stdout     io.ReadCloser
	stderr     io.ReadCloser
}

func NewBgCmd(name string, cmd *exec.Cmd) *BgCmd {
	var err error

	bg := &BgCmd{
		Name:       name,
		Cmd:        cmd,
		SocketPath: fmt.Sprintf("/tmp/godutch-%s.sock", name),
	}

	os.Setenv("GODUTCH_SOCKET_PATH", "")
	bg.Env = bg.setenv("GODUTCH_SOCKET_PATH", bg.SocketPath)
	bg.Cmd.Env = bg.Env

	if bg.stdout, err = bg.Cmd.StdoutPipe(); err != nil {
		log.Fatalln(err)
	}
	if bg.stderr, err = bg.Cmd.StderrPipe(); err != nil {
		log.Fatalln(err)
	}

	return bg
}

func (bg *BgCmd) setenv(key string, value string) []string {
	var env []string = os.Environ()
	var newEnv []string
	var newEntry string

	for _, entry := range env {
		name := strings.Split(entry, "=")
		if name[0] == key && name[1] != value {
			newEntry = fmt.Sprintf("%s=%s", key, value)
			newEnv = append(newEnv, newEntry)
			log.Println("ENV:", newEntry)
		}
	}

	return newEnv
}

func (bg *BgCmd) captureOutput() {
	var err error
	var multi io.Reader

	multi = io.MultiReader(bg.stdout, bg.stderr)
	output := bufio.NewScanner(multi)

	for output.Scan() {
		log.Println("[", bg.Name, "STDOUT ]:", output.Text())
	}

	if err = output.Err(); err != nil {
		log.Println("[", bg.Name, "STDERR ]:", err)
	}
}

func (bg *BgCmd) Serve() {
	var err error

	log.Println("Starting to 'serve':", bg.Name)

	if err = bg.Cmd.Start(); err != nil {
		log.Fatalln("Start error:", err)
	}

	bg.captureOutput()

	if err = bg.Cmd.Wait(); err != nil {
		log.Println("Wait error:", err)
	}
}

func (bg *BgCmd) Stop() {
	if err := bg.Cmd.Process.Kill(); err != nil {
		log.Fatalln("Error on kill: ", err)
	}
}

/* EOF */
