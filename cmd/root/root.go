package root

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	restoreDb "task-runner-cobra/cmd/db/restore"
	buildFrontend "task-runner-cobra/cmd/frontend/build"
	"task-runner-cobra/cmd/grpc/generate"
	"task-runner-cobra/cmd/grpc/info"
	envFinder "task-runner-cobra/utils/env/finder"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/config.yaml)")
	rootCmd.AddCommand(restoreDb.Cmd)
	rootCmd.AddCommand(buildFrontend.Cmd)
	rootCmd.AddCommand(generate.Cmd)
	rootCmd.AddCommand(info.Cmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	for _, val := range viper.AllKeys() {
		data := viper.GetString(val)
		if data != "" {
			envMap := envFinder.FindEnvBetweenQuotes(data)
			for withQuotes, value := range envMap {
				data = strings.ReplaceAll(data, withQuotes, value)
			}
			viper.Set(val, data)
		}
	}
}
