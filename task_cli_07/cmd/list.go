package cmd

import (
	"fmt"
	"os"

	"github.com/hippo-an/tiny-go-challenges/task_cli_07/db"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of your tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("Add your first list 🤩")
			return
		}

		fmt.Println("You have the following tasks: ")
		for i, v := range tasks {
			fmt.Printf("%d. %s(%d)\n", i+1, v.Task, v.Id)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
