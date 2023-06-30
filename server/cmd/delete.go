package cmd

import "github.com/spf13/cobra"

func deleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete REPOSITORY",
		Aliases:           []string{"del", "remove", "rm"},
		Short:             "Delete a repository",
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: checkIfCollab,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			be, _ := fromContext(cmd)
			name := args[0]
			if err := be.DeleteRepository(ctx, name); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
