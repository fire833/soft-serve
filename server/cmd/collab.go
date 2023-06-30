package cmd

import (
	"github.com/spf13/cobra"
)

func collabCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collab",
		Aliases: []string{"collabs", "collaborator", "collaborators"},
		Short:   "Manage collaborators",
	}

	cmd.AddCommand(
	// collabAddCommand(),
	// collabRemoveCommand(),
	// collabListCommand(),
	)

	return cmd
}

//
// func collabAddCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:               "add REPOSITORY USERNAME",
// 		Short:             "Add a collaborator to a repo",
// 		Args:              cobra.ExactArgs(2),
// 		PersistentPreRunE: checkIfCollab,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			be, _ := fromContext(cmd)
// 			repo := args[0]
// 			username := args[1]
//
// 			return be.AddCollaborator(ctx, repo, username)
// 		},
// 	}
//
// 	return cmd
// }
//
// func collabRemoveCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:               "remove REPOSITORY USERNAME",
// 		Args:              cobra.ExactArgs(2),
// 		Short:             "Remove a collaborator from a repo",
// 		PersistentPreRunE: checkIfCollab,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			be, _ := fromContext(cmd)
// 			repo := args[0]
// 			username := args[1]
//
// 			return be.RemoveCollaborator(repo, username)
// 		},
// 	}
//
// 	return cmd
// }
//
// func collabListCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:               "list REPOSITORY",
// 		Short:             "List collaborators for a repo",
// 		Args:              cobra.ExactArgs(1),
// 		PersistentPreRunE: checkIfCollab,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			be, _ := fromContext(cmd)
// 			repo := args[0]
// 			collabs, err := be.Collaborators(repo)
// 			if err != nil {
// 				return err
// 			}
//
// 			for _, c := range collabs {
// 				cmd.Println(c)
// 			}
//
// 			return nil
// 		},
// 	}
//
// 	return cmd
// }
