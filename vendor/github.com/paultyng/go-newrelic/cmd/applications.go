package cmd

import (
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

func makeApplicationsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "applications",
		Aliases: []string{"application", "apps", "app"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getApplicationsCmd = makeApplicationsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		resources, err := client.ListApplications()
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getApplicationsCmd)
}
