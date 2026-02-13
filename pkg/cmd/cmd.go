package cmd

import (
	"fmt"
	"os"

	"github.com/feloy/tesh/pkg/run"
	"github.com/spf13/cobra"
)

type TeshCmdOptions struct {
	ScenariosFile    string
	SingleScenarioID string
	CoverageFile     string
}

func NewTesh() *cobra.Command {
	o := TeshCmdOptions{}
	cmd := &cobra.Command{
		Use:   "tesh <script file> [--scenarios <scenarios file> [--scenario <scenario id>] ]",
		Short: "Run and Test your shell scripts",
		Long:  "Tesh is a tool for running and testing shell scripts.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if o.SingleScenarioID != "" && o.ScenariosFile == "" {
				return fmt.Errorf("--scenario flag requires --scenarios flag")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			scriptFile, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("failed to open script file: %w", err)
			}
			defer scriptFile.Close()

			if o.ScenariosFile == "" {
				run.Script(scriptFile, run.ScriptOptions{WithCoverage: o.CoverageFile, FilePath: args[0]})
				return nil
			}

			scenariosFile, err := os.Open(o.ScenariosFile)
			if err != nil {
				return fmt.Errorf("failed to open scenarios file: %w", err)
			}
			defer scenariosFile.Close()

			var singleScenarioID *string
			if o.SingleScenarioID != "" {
				singleScenarioID = &o.SingleScenarioID
			}

			results := run.Scenarios(scriptFile, scenariosFile, run.ScenariosOptions{
				SingleScenarioID: singleScenarioID,
				WithCoverage:     o.CoverageFile,
				FilePath:         args[0],
			})

			// TODO move to specific outputter
			var exitCode int = 0
			for _, result := range results {
				fmt.Printf("Scenario: %s\n", result.ScenarioID)
				if result.ExitCode != nil {
					fmt.Printf("Exit Code: expected %d, actual %d\n", result.ExitCode.Expected, result.ExitCode.Actual)
				}
				if result.Stdout != nil {
					fmt.Printf("Stdout: expected %q, actual %q\n", result.Stdout.Expected, result.Stdout.Actual)
				}
				if result.Stderr != nil {
					fmt.Printf("Stderr: expected %q, actual %q\n", result.Stderr.Expected, result.Stderr.Actual)
				}
				for _, call := range result.Calls {
					fmt.Printf("Call: %s %v, expected %d calls, actual %d calls\n", call.Command, call.Args, call.ExpectedCalled, call.ActualCalled)
				}
				if !result.IsSuccess() {
					exitCode = 1
				}
			}
			os.Exit(exitCode)
			return nil
		},
	}

	cmd.Flags().StringVar(&o.ScenariosFile, "scenarios", "", "Scenarios file to use for testing the script, if not given the script is executed but not tested")
	cmd.Flags().StringVar(&o.SingleScenarioID, "scenario", "", "Single scenario ID to run, if not given all scenarios are run")
	cmd.Flags().StringVar(&o.CoverageFile, "coverage", "", "Display covered lines during execution or scenarios. Suppress normal stdout and stderr if empty and display colored script, otherwise write to the given file in Go coverage.txt format.")
	cmd.Flag("coverage").NoOptDefVal = "-"

	return cmd
}
