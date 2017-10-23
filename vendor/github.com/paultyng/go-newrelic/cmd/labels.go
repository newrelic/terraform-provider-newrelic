package cmd

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/paultyng/go-newrelic/api"
	"github.com/spf13/cobra"
)

func makeLabelsCmd(dst cobra.Command) *cobra.Command {
	src := cobra.Command{
		Use:     "labels",
		Aliases: []string{"label", "l"},
	}

	if err := mergo.Merge(&dst, src); err != nil {
		panic(err)
	}

	return &dst
}

var getLabelsCmd = makeLabelsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}
		resources, err := client.ListLabels()
		if err != nil {
			return err
		}

		return outputList(cmd, resources)
	},
})

var createLabelsCmd = makeLabelsCmd(cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newAPIClient(cmd)
		if err != nil {
			return err
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		category, err := cmd.Flags().GetString("category")
		if err != nil {
			return err
		}

		err = client.CreateLabel(api.Label{
			Name:     name,
			Category: category,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Label '%v' created successfully.\n", name)

		return nil
	},
})

func init() {
	getCmd.AddCommand(getLabelsCmd)
	createCmd.AddCommand(createLabelsCmd)

	createLabelsCmd.Flags().String("name", "", "Name of the label")
	createLabelsCmd.Flags().String("category", "", "Category of the label")
}
