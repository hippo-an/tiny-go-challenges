package cmd

import (
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/task_cli_07/db"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list.",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(task, "added to your list.")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
