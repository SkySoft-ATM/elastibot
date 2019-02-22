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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var filename string

func init() {
	rootCmd.AddCommand(putCmd)
	putCmd.Flags().StringVarP(&filename, "filename", "f", "", "file path")
	putCmd.MarkFlagRequired("filename")
}

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Store a configuration to a resource by filename (i.e. templates, mapping etc...)",
	Run: func(cmd *cobra.Command, args []string) {
		defer esClient.Stop()
		filename := cmd.Flag("filename").Value.String()
		if len(args) < 1 {
			fmt.Printf("you need to give a name for your template")
			os.Exit(1)
		}
		UploadElasticSearchTemplate()

	},
}

func UploadElasticSearchTemplate(templateName string, templateFile string) error {

	fileContent, err := ioutil.ReadFile(templateFile)
	if err != nil {
		fmt.Printf("Failed to Read the File %v\n", templateFile)
		return err
	}

	client := &http.Client{}
	client.Timeout = time.Second * 15

	uri := url + "/_template/" + templateName
	body := bytes.NewBuffer(fileContent)
	req, err := http.NewRequest(http.MethodPut, uri, body)
	if err != nil {
		fmt.Printf("http.NewRequest() failed with %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do() failed with %v\n", err)
		return err
	}

	defer resp.Body.Close()
	var response []byte
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() failed with %v\n", err)
		return err
	}

	fmt.Printf("Response status code: %v, text:%v\n", resp.StatusCode, string(response))
	if resp.StatusCode == 200 {
		fmt.Printf("Template has been uploaded to ES: %v\n", string(fileContent))
	} else {
		fmt.Printf("Template has NOT been uploaded to ES\n")
	}
}
