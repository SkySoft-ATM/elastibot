// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(delCmd)
	delCmd.AddCommand(delAllCmd)
}

// clearCmd represents the clear command
var delCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete index",
	Run: func(cmd *cobra.Command, args []string) {
		defer esClient.Stop()
		if len(args) < 1 {
			fmt.Printf("✘ you need to specify an index name\n")
			os.Exit(1)
		}
		index := args[0]
		err := deleteIndex(index)
		if err != nil {
			fmt.Printf("✘ Error trying to delete index %s\n", index)
			fmt.Printf("%s\n", err)
		}
	},
}

func deleteIndex(name string) error {
	fmt.Printf("Deleting index name %s...\n", name)
	_, err := esClient.DeleteIndex(name).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("✔ index [ %s ] deleted\n", name)
	return nil
}

var delAllCmd = &cobra.Command{
	Use:   "all",
	Short: `Delete all indexes`,
	Run: func(cmd *cobra.Command, args []string) {
		names, err := esClient.IndexNames()
		if err != nil {
			fmt.Printf("✘ Error trying to delete all indexes")
			os.Exit(1)
		}
		for _, name := range names {
			err := deleteIndex(name)
			if err != nil {
				fmt.Printf("✘ Error trying to delete index [ %s ]\n", name)
			}
		}
	},
}

//✘ ✔
