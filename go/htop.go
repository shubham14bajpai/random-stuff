// usage
// go run htop.go --name="kube-apiserver" | sh -

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

var (
	process  string
	duration int
)

func init() {
	flag.StringVar(&process, "name", "", "name of the process")
	flag.IntVar(&duration, "duration", 100, "duration to sleep and watch htop")
	flag.Parse()
}

func main() {

	fmt.Println(process)

	cmd := exec.Command("pgrep", "-f", process)

	stdout, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	pids := strings.ReplaceAll(strings.TrimSuffix(string(stdout), "\n"), "\n", ",")
	fmt.Println("htop -p", pids)

	time.Sleep(time.Duration(duration) * time.Second)
}
