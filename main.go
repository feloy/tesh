package main

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/feloy/tesh/pkg/run"
)

type nopWriterCloser struct {
	*strings.Reader
}

func (nopWriterCloser) Write([]byte) (int, error) { return 0, io.EOF }
func (nopWriterCloser) Close() error              { return nil }

func main() {
	if len(os.Args) < 2 {
		log.Fatal("script file is required")
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if len(os.Args) == 2 {
		run.Script(file)
		return
	}

	scenariosFile, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	var singleScenarioID *string
	if len(os.Args) > 3 {
		singleScenarioID = &os.Args[3]
	}

	run.Scenarios(file, scenariosFile, singleScenarioID)
}
