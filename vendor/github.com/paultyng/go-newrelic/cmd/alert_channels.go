package cmd

import (
	"github.com/imdario/mergo"
	"github.com/paultyng/go-newrelic/api"
	"github.com/spf13/cobra"
)

func makeChannelsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "channels",
		Aliases: []string{"channel", "ch"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getAlertChannelsCmd = makeChannelsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}

		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}

		var resources []api.AlertChannel

		if id != 0 {
			resource, err := client.GetAlertChannel(id)
			if err != nil {
				return err
			}

			resources = []api.AlertChannel{*resource}
		} else {
			resources, err = client.ListAlertChannels()
			if err != nil {
				return err
			}
		}

		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getAlertChannelsCmd)
	getAlertChannelsCmd.Flags().Int("id", 0, "ID of the alert channel to get")
}
