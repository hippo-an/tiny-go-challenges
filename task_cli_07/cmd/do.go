package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hippo-an/tiny-go-challenges/task_cli_07/db"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete.",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the argument:", arg)
				continue
			}

			ids = append(ids, id)
		}

		if len(ids) == 0 {
			fmt.Println("pass the id of the list")
			return
		}

		for _, id := range ids {
			err := db.DeleteTask(int64(id))
			if err != nil {
				fmt.Println("Failed to do:", id)
				continue
			}
		}

		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete! Take a vacation 🚀")
			return
		}

		fmt.Println("You have the following tasks: ")
		for i, v := range tasks {
			fmt.Printf("%d. %s(%d)\n", i+1, v.Task, v.Id)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
