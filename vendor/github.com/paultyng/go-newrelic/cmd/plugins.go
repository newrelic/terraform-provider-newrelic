package cmd

import (
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

func makePluginsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "plugins",
		Aliases: []string{"plugin", "plug"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getPluginsCmd = makePluginsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		resources, err := client.ListPlugins()
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getPluginsCmd)
}
