package godutch

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	command    []string
	Env        []string
	stdout     io.ReadCloser
	stderr     io.ReadCloser
}

// Creates a new BgCmd object, which will prepare socket and os/exec command to
// run in background, after "Bootstrap".
func NewBgCmd(containerCfg *ContainerConfig) *BgCmd {
	var okay bool = false
	var bg *BgCmd
	var socketName string = fmt.Sprintf("godutch-%s.sock", containerCfg.Name)
	var socketPath string = filepath.Join(containerCfg.SocketDir, socketName)

	if okay, _ = exists(socketPath); okay {
		log.Printf("[BgCmd] [WARN] Socket already found at: '%s'", socketPath)
	}

	bg = &BgCmd{
		Name:       containerCfg.Name,
		command:    containerCfg.Command,
		SocketPath: socketPath,
	}

	// socket information, basic commnicaton method with background process
	os.Setenv("GODUTCH_SOCKET_PATH", "")
	bg.Env = bg.setenv("GODUTCH_SOCKET_PATH", bg.SocketPath)

	return bg
}

// Start serving a background command. Spawn the command and handles the "wait"
// call, trowing stdout/stderr entries on log interface.
func (bg *BgCmd) Serve() {
	var err error
	log.Println("[BgCmd] Starting to 'serve':", bg.Name)

	if err = bg.spawnCmd(); err != nil {
		log.Fatalln("[BgCmd] Spawn error:", err)
	}

	if err = bg.Cmd.Start(); err != nil {
		log.Fatalln("[BgCmd] Start error:", err)
	}

	bg.captureOutput()

	if err = bg.Cmd.Wait(); err != nil {
		log.Println("[BgCmd] Wait error:", err)
	}
}

// Handles the creation of a new exec.Command instance with informed parameters
// and custom environment.
func (bg *BgCmd) spawnCmd() error {
	var err error

	bg.Cmd = exec.Command(bg.command[0], bg.command[1:]...)
	bg.Cmd.Env = bg.Env

	// leaving stdout and stderr pipes for capturing outputs
	if bg.stdout, err = bg.Cmd.StdoutPipe(); err != nil {
		return err
	}

	if bg.stderr, err = bg.Cmd.StderrPipe(); err != nil {
		return err
	}

	return nil
}

// Stop a background command.
func (bg *BgCmd) Stop() {
	var err error
	if err = bg.Cmd.Process.Kill(); err != nil {
		log.Println("[BgCmd] Error on kill: ", err)
	}
}

// Helper method to set-up a valid environment slice, adding the informed
// arguements in a key-value fashion.
func (bg *BgCmd) setenv(key string, value string) []string {
	var env []string = os.Environ()
	var newEnv []string
	var newEntry string
	var entry string
	var keyValue []string

	for _, entry = range env {
		keyValue = strings.Split(entry, "=")
		if keyValue[0] == key && keyValue[1] != value {
			newEntry = fmt.Sprintf("%s=%s", key, value)
			newEnv = append(newEnv, newEntry)
			log.Println("[BgCmd] ENV:", newEntry)
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
		log.Println("[BgCmd] [", bg.Name, "STDOUT ]:", output.Text())
	}

	if err = output.Err(); err != nil {
		log.Println("[BgCmd] [", bg.Name, "STDERR ]:", err)
	}
}

/* EOF */
