package godutch

/*
import (
	"fmt"
	"github.com/thejerf/suture"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

type Monitor struct {
	*suture.Supervisor
	name  string
	inbox chan []byte
	bg    *BgCommand
}

type BgCommand struct {
	cmd *exec.Cmd
}

func (bg *BgCommand) Serve() {
	env := os.Environ()
	bg.cmd.Env = append(
		env,
		fmt.Sprintf("GODUTCH_SOCKET_PATH=/tmp/%s.sock", "testMonitor"),
	)
	log.Println("Starting...")
	if err := bg.cmd.Start(); err != nil {
		log.Fatalln("Start: ", err)
	}
	if err := bg.cmd.Wait(); err != nil {
		log.Println("Wait: ", err)
	}
}

func (bg *BgCommand) Stop() {
	if err := bg.cmd.Process.Kill(); err != nil {
		log.Fatal("failed to kill: ", err)
	}
}

func NewMonitor(name string, cmd *exec.Cmd) (*Monitor, error) {
	m := &Monitor{
		Supervisor: suture.New(name, suture.Spec{
			Log: func(line string) {
				log.Println("Suture:", line)
			},
		}),
		bg:    &BgCommand{cmd: cmd},
		name:  name,
		inbox: make(chan []byte),
	}
	m.Add(m.bg)
	return m, nil
}

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("Client got:", string(buf[0:n]))
	}
}

func main() {
	env := os.Environ()
	env = append(env, fmt.Sprintf("GODUTCH_SOCKET_PATH=/tmp/%s.sock", "testMonitor"))
	m, err := NewMonitor(
		"testMonitor",
		exec.Command(
			"/opt/chefdk/embedded/bin/ruby",
			"/Users/ofernandes/src/go/tmp/starlite/godutch_test.rb",
		))
	log.Println("m:", m)
	log.Println("err:", err)

	go func() {
		time.Sleep(3 * 1e9)
		c, err := net.Dial("unix", "/tmp/godutch_test.sock")
		if err != nil {
			panic(err)
		}
		defer c.Close()

		go reader(c)
		for {
			log.Println("*")
			_, err := c.Write([]byte("{ \"command\": \"__list_check_methods\", \"arguments\": [] }"))
			if err != nil {
				log.Fatal("write error:", err)
				break
			}
			time.Sleep(1e9)
		}
	}()

	m.Serve()
}

*/
/* EOF */
