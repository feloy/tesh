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
				run.Script(scriptFile)
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

			run.Scenarios(scriptFile, scenariosFile, singleScenarioID)
			return nil
		},
	}

	cmd.Flags().StringVar(&o.ScenariosFile, "scenarios", "", "Scenarios file to use for testing the script, if not given the script is executed but not tested")
	cmd.Flags().StringVar(&o.SingleScenarioID, "scenario", "", "Single scenario ID to run, if not given all scenarios are run")

	return cmd
}
