package cmd

import (
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

func makeConditionsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "conditions",
		Aliases: []string{"condition", "cond"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getAlertConditionsCmd = makeConditionsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		policyID, err := cmd.Flags().GetInt("policy-id")
		if err != nil {
			return err
		}
		resources, err := client.ListAlertConditions(policyID)
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getAlertConditionsCmd)
	getAlertConditionsCmd.Flags().IntP("policy-id", "p", 0, "ID of policy for which to get conditions")
}
