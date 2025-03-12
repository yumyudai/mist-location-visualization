package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"mist-location-visualization/internal/mistpoller"
)

func main() {
	var err error
	var configFile string
	var config mistpoller.Config

	rootCmd := &cobra.Command {
		Use: "mistpoller",
		Short: "Poll Mist API and push changes to database",
		// Main Entry Point
		Run: func(c *cobra.Command, args []string) {
			// Init 
			rcvr, err := mistpoller.New(config)
			if err != nil {
				log.Fatalf("Failed on init: %v", err)
			}

			err = rcvr.Run()
			if err != nil {
				log.Fatalf("Failed on start: %v", err)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.json", "Path to configuration")

	// Default Values
	viper.SetDefault("mist.endpoint", "api.mist.com")

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
