package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"mist-location-visualization/internal/locapiserver"
)

func main() {
	var err error
	var configFile string
	var config locapiserver.Config

	rootCmd := &cobra.Command {
		Use: "locapid",
		Short: "API server for BLE tag location visualization using data from Juniper Mist",
		// Main Entry Point
		Run: func(c *cobra.Command, args []string) {
			// Init 
			e, err := locapiserver.New(config)
			if err != nil {
				log.Fatalf("Failed on init: %v", err)
			}

			err = e.Run()
			if err != nil {
				log.Fatalf("Failed on start: %v", err)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.json", "Path to configuration")

	// Defaults
	viper.SetDefault("mist.endpoint", "api.mist.com")
	viper.SetDefault("mist.location_timeout", 60)
	viper.SetDefault("mist.refresh_time", 1800)

	// Read Configuration File Before Start
	cobra.OnInitialize(func() {
		_, err := os.Stat(configFile)
		if os.IsNotExist(err) {
			envConfFile := os.Getenv("CONFIG_FILE")
			if envConfFile != "" {
				_, err := os.Stat(envConfFile)
				if os.IsNotExist(err) {
					log.Fatalf("Config file %s does not exist!", envConfFile)
				}

				configFile = envConfFile
			} else {
				log.Fatalf("Config file %s does not exist!", configFile)
			}
		}

		viper.SetConfigFile(configFile)
		viper.SetConfigType("json")
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("Failed to parse config: %v", err)
		}

		log.Printf("Loaded config file: %s", configFile)
	})

	// Launch (cobra.OnInitializa -> rootCmd.Run)
	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

}
