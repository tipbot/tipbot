package main

import (
	"fmt"
	"log"
	//	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var app *App
var rootCmd *cobra.Command

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd.Execute()
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	rootCmd = &cobra.Command{
		Use:   "tipbot-back",
		Short: "tipbot backend server",
		Long: `tipbot backend server
=========================

Make sure config.toml file is in the working folder.
Required config values:
 `,
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Print("Reading config.toml file")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	config := Config{
		GITHUB_ACCESS_TOKEN:  viper.GetString("GITHUB_ACCESS_TOKEN"),
		BOT_GITHUB_NAME:      "@" + viper.GetString("BOT_GITHUB_NAME"),
		DB_CONNECTION_STRING: viper.GetString("DB_CONNECTION_STRING"),
		MIN_XML_BALANCE:      viper.GetFloat64("MIN_XLM_BALANCE"),
		HORIZON_URL:          viper.GetString("HORIZON_URL"),
		WEBSITE_URL:          viper.GetString("WEBSITE_URL"),
		FEDERATION_DOMAIN:    viper.GetString("FEDERATION_DOMAIN"),
		DEFAULT_TIP_AMOUNT:   viper.GetFloat64("DEFAULT_TIP_AMOUNT"),
		BRIDGE_URL:    	      viper.GetString("BRIDGE_URL"),
	}


	app, err = NewApp(config)

	if err != nil {
		log.Fatal(err.Error())
	}

	app.Run()
}
