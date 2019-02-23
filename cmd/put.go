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
	"bytes"
	"context"
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
	putCmd.Flags().StringVarP(&filename, "file", "f", "", "file path")
	putCmd.MarkFlagRequired("file")
}

// putCmd represents the put command
var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Store a template defined by filename",
	Run: func(cmd *cobra.Command, args []string) {
		defer esClient.Stop()
		file := cmd.Flag("file").Value.String()
		if len(args) < 1 {
			fmt.Printf("you need to give a name for your template")
			os.Exit(1)
		}
		name := args[0]
		content, err := getFileContent(file)
		if err != nil {
			fmt.Printf("✘ Error trying to read file %s\n", file)
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		res, err := esClient.IndexPutTemplate(name).BodyString(string(content)).Create(true).Do(context.Background())
		if err != nil {
			fmt.Printf("✘ Error trying to store template %s\n", name)
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Template successfully added %t\n", res.Acknowledged)
		os.Exit(0)
	},
}

func getFileContent(filepath string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Failed to Read the File %v\n", filepath)
		return nil, err
	}
	return fileContent, nil
}

func uploadElasticSearchTemplate(templateName string, templateFile string) error {

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
	return nil
}
