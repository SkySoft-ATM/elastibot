// Copyright © 2019 SkySoft-ATM <chambodn@skysoft-atm.com>
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
	* color (aka 'c')
	* index (aka 'i')
	* template (aka 'tpl')
	* version (aka 'v')
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

		switch args[0] {
		case "color", "c":
			res, err := esClient.ClusterHealth().Do(context.Background())
			if err != nil {
				fmt.Printf("✘ error trying to retrieve indexes on %s\n", url)
				fmt.Printf("✘ %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", res.Status)
			os.Exit(0)
		case "index", "i":
			names, err := esClient.IndexNames()
			if err != nil {
				fmt.Printf("✘ error trying to retrieve indexes on %s\n", url)
				fmt.Printf("✘ %s\n", err)
				os.Exit(1)
			}
			for _, name := range names {
				fmt.Printf("%s\n", name)
			}
			os.Exit(0)
		case "template", "tpl":
			templates, err := esClient.IndexGetTemplate("_all").Do(context.Background())
			if err != nil {
				fmt.Printf("✘ error trying to retrieve templates on %s\n", url)
				fmt.Printf("✘ %s\n", err)
				os.Exit(1)
			}
			for _, template := range templates {
				fmt.Printf("%v\n", template.IndexPatterns)
			}
			os.Exit(0)
		case "version", "v":
			esVersion, err := esClient.ElasticsearchVersion(url)
			if err != nil {
				fmt.Printf("✘ no elasticsearch found in [ %s ]\n", url)
				os.Exit(1)
			}
			fmt.Printf("%s\n", esVersion)
		}

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
