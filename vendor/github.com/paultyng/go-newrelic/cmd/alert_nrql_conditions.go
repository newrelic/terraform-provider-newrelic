package cmd

import (
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

func makeAlertNrqlConditionsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "nrql-conditions",
		Aliases: []string{"nrql-condition", "nrql-cond"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getAlertNrqlConditionsCmd = makeAlertNrqlConditionsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		policyID, err := cmd.Flags().GetInt("policy-id")
		if err != nil {
			return err
		}
		resources, err := client.ListAlertNrqlConditions(policyID)
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

func init() {
	getCmd.AddCommand(getAlertNrqlConditionsCmd)
	getAlertNrqlConditionsCmd.Flags().IntP("policy-id", "p", 0, "ID of policy for which to get conditions")
}
