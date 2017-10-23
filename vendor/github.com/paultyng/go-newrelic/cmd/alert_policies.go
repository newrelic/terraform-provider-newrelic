package cmd

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/paultyng/go-newrelic/api"
	"github.com/spf13/cobra"
)

func makePoliciesCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "policies",
		Aliases: []string{"policy", "pol"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getAlertPoliciesCmd = makePoliciesCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		resources, err := client.ListAlertPolicies()
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

var createAlertPoliciesCmd = makePoliciesCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		incidentPreference, err := cmd.Flags().GetString("incident-preference")
		if err != nil {
			return err
		}

		_, err = client.CreateAlertPolicy(api.AlertPolicy{
			Name:               name,
			IncidentPreference: incidentPreference,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Alert policy '%v' created.\n", name)

		return nil
	},
})

var deleteAlertPoliciesCmd = makePoliciesCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}

		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			return err
		}

		err = client.DeleteAlertPolicy(id)
		if err != nil {
			return err
		}

		fmt.Printf("Alert policy '%v' deleted.\n", id)

		return nil
	},
})

func init() {
	getCmd.AddCommand(getAlertPoliciesCmd)

	createCmd.AddCommand(createAlertPoliciesCmd)
	createAlertPoliciesCmd.Flags().String("name", "", "Name of the alert policy")
	createAlertPoliciesCmd.Flags().String("incident-preference", "PER_POLICY", "Incident preference of the policy")

	deleteCmd.AddCommand(deleteAlertPoliciesCmd)
	deleteAlertPoliciesCmd.Flags().Int("id", 0, "ID of the alert policy to delete")
}
