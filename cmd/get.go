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

const availableRes = `
	* index
	* template
`

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",

	Run: func(cmd *cobra.Command, args []string) {
		defer esClient.Stop()
		if len(args) < 1 {
			fmt.Printf("✘ You must specify the type of resource to get. Valid resource types include: %s\n", availableRes)

			os.Exit(1)
		}
		resource := args[0]
		if resource == "index" {
			names, err := esClient.IndexNames()
			if err != nil {
				// Handle error
				panic(err)
			}
			for _, name := range names {
				fmt.Printf("%s\n", name)
			}
			return
		}

		if resource == "template" {
			esClient.IndexGetTemplate("_all").Do(context.Background())
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
