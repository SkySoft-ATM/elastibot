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
	"fmt"
	"log"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/olivere/elastic"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var url string
var esClient *elastic.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elasticbot",
	Short: "Little CLI to facilitate interactions with Elasticsearch clusters",
	Long: `elasticbot allows you to quickly interact with your Elasticsearch clusters.
It is mainly used for administration purposes.`,
	PersistentPreRun:  elasticConnect,
	PersistentPostRun: elasticDisconnect,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.elasticbot.yaml)")
	rootCmd.PersistentFlags().StringVar(&url, "url", "http://localhost:9200", "Elasticsearch url")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".elasticbot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".elasticbot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func elasticDisconnect(cmd *cobra.Command, args []string) {
	esClient.Stop()
}

func elasticConnect(cmd *cobra.Command, args []string) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetRetrier(NewCustomRetrier()),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	if err != nil {
		fmt.Printf("✘ Failed to connect to Elasticsearch, configuration is not correct\nCheck url [ %s ]\n", url)
		os.Exit(1)
	}
	esClient = client
}
