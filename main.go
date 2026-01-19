package main

import (
	"log"
	"os"

	"github.com/feloy/tesh/pkg/run"
)

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
	defer scenariosFile.Close()

	var singleScenarioID *string
	if len(os.Args) > 3 {
		singleScenarioID = &os.Args[3]
	}

	run.Scenarios(file, scenariosFile, singleScenarioID)
}
