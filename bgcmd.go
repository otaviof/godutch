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

//
// BgCmd type is the representation of any command to run in background via
// GoDutch, with socket and other attributes to handle output and IPC.
//
type BgCmd struct {
	Name       string
	SocketPath string
	Cmd        *exec.Cmd
	Env        []string
	stdout     io.ReadCloser
	stderr     io.ReadCloser
}

// Creates a new BgCmd object, which will prepare socket and os/exec command to
// run in background, after "Bootstrap".
func NewBgCmd(name string, cmd *exec.Cmd) *BgCmd {
	var err error

	bg := &BgCmd{
		Name:       name,
		Cmd:        cmd,
		SocketPath: fmt.Sprintf("/tmp/godutch-%s.sock", name),
	}

	// socket information, basic commnicaton method with background process
	os.Setenv("GODUTCH_SOCKET_PATH", "")
	bg.Env = bg.setenv("GODUTCH_SOCKET_PATH", bg.SocketPath)
	bg.Cmd.Env = bg.Env

	// leaving stdout and stderr pipes for capturing outputs
	if bg.stdout, err = bg.Cmd.StdoutPipe(); err != nil {
		log.Fatalln(err)
	}
	if bg.stderr, err = bg.Cmd.StderrPipe(); err != nil {
		log.Fatalln(err)
	}

	return bg
}

// Helper method to set-up a valid environment slice, adding the informed
// arguements in a key-value fashion.
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

// Reads stdout and stderr io.Readers in a single Scanner, feeding the log
// interface with what's found. Bufferized IO is used here, avoid blocking.
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

// Start serving a background command. Spawn the command and handles the "wait"
// call, trowing stdout/stderr entries on log interface.
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

// Stop a background command.
func (bg *BgCmd) Stop() {
	if err := bg.Cmd.Process.Kill(); err != nil {
		log.Fatalln("Error on kill: ", err)
	}
}

/* EOF */
