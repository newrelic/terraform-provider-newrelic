package cmd

import (
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

func makeComponentMetricsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "component-metrics",
		Aliases: []string{"component-metric", "cm"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getComponentMetricsCmd = makeComponentMetricsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		componentID, err := cmd.Flags().GetInt("component-id")
		if err != nil {
			return err
		}
		resources, err := client.ListComponentMetrics(componentID)
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getComponentMetricsCmd)
	getComponentMetricsCmd.Flags().IntP("component-id", "c", 0, "ID of component")
}
