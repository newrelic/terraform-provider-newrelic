package cmd

import (
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

func makeComponentMetricDataCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "component-metric-data",
		Aliases: []string{"cmd"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getComponentMetricDataCmd = makeComponentMetricDataCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		componentID, err := cmd.Flags().GetInt("component-id")
		if err != nil {
			return err
		}
		names, err := cmd.Flags().GetStringSlice("name")
		if err != nil {
			return err
		}
		resources, err := client.ListComponentMetricData(componentID, names)
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getComponentMetricDataCmd)
	getComponentMetricDataCmd.Flags().IntP("component-id", "c", 0, "ID of component")
	getComponentMetricDataCmd.Flags().StringSliceP("name", "n", []string{}, "List of names")
}
