package autofresh

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	rw "github.com/TerrenceHo/monorepo/autofresh/watcher"
	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
)

const logo = `
    ___         __        ______               __  
   /   | __  __/ /_____  / ____/_______  _____/ /_ 
  / /| |/ / / / __/ __ \/ /_  / ___/ _ \/ ___/ __ \
 / ___ / /_/ / /_/ /_/ / __/ / /  /  __(__  ) / / /
/_/  |_\__,_/\__/\____/_/   /_/   \___/____/_/ /_/ 
`

type Config struct {
	Cmd        string
	Args       []string
	Extensions []string
	HideBanner bool
	Watch      []string
}

func Start(c Config) {
	// if !c.HideBanner {
	// 	fmt.Println(logo)
	// }
	fmt.Println(stringifyCommand(c.Cmd, c.Args))
	fmt.Println(c.Watch)
	fmt.Println(c.Extensions)

	stopChan := make(chan bool)
	Run(stopChan, c.Cmd, c.Args)

	stopChan <- true

	time.Sleep(1 * time.Second)
	watcher, err := rw.New()
	if err != nil {
		log.Fatalf("Failed to start recursive watcher: %+v\n", err)
	}
	for _, path := range c.Watch {
		err := watcher.AddRecursive(path)
		if err != nil {
			log.Fatalf("Failed to start watching %s: %v", path, err)
		}
	}
}

// Run takes in a boolean channel, a command string, and its associated
// arguments, and starts a command asynchronously. It will redirect
// stdout/stderr output from the command back to the OS stdout/stderr. The
// stopChannel is used to kill the process when a true value is sent to the
// channel.  If starting the command fails, Run returns a wrapped error,
// otherwise Run returns nil.
func Run(stopChannel chan bool, cmd string, args []string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	// creates a closure that reads from stopChannel. If stopChannel receives a
	// true and the process has not already exited, then it will attempt to kill
	// the process. If killing the process fails, then it will log this event
	// fatally.
	go func() {
		for {
			stop := <-stopChannel
			fmt.Printf("STOP is %t\n", stop)
			if stop {
				pid := c.Process.Pid
				if err := c.Process.Kill(); err != nil {
					log.Fatalf("Failed to kill process %d with error: %v\n", pid, err)
				}
				break
			}
		}
	}()

	if err := c.Start(); err != nil {
		return stackerrors.Wrapf(err, "failed to start command (%s)", stringifyCommand(cmd, args))
	}
	return nil
}

func stringifyCommand(cmd string, args []string) string {
	return fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
}
