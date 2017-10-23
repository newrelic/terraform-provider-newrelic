package cmd

import (
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

func makeComponentsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "components",
		Aliases: []string{"component", "comp"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getComponentsCmd = makeComponentsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		pluginID, err := cmd.Flags().GetInt("plugin-id")
		if err != nil {
			return err
		}
		resources, err := client.ListComponents(pluginID)
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getComponentsCmd)
	getComponentsCmd.Flags().IntP("plugin-id", "p", 0, "ID of policy for which to get conditions")
}
