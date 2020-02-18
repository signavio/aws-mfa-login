package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/signavio/aws-mfa-login/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var cfgFile string
var Name string
var Destination string

var (
	conf *action.Clusters
)

var (
	VERSION = "0.1.1"
)

var rootCmd = &cobra.Command{
	Use:     "aws-mfa-login",
	Short:   "aws login with mfa",
	Long:    "CLI tool to update your temporary AWS credentials ",
	Version: VERSION,
	Run: func(cmd *cobra.Command, args []string) {
		action.PrintConfigWithoutClusterConfig()
		action.UpdateSessionCredentials()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aws-mfa.yaml)")
	rootCmd.PersistentFlags().StringVarP(&Name, "source", "s", "", "source profile where mfa is activated")
	rootCmd.PersistentFlags().StringVarP(&Destination, "destination", "d", "", "destination profile for temporary aws credentials")
	rootCmd.InitDefaultVersionFlag()
	err := viper.BindPFlag("source", rootCmd.PersistentFlags().Lookup("source"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("destination", rootCmd.PersistentFlags().Lookup("destination"))
	if err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".aws-mfa")
	}
	// read in environment variables
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	conf = &action.Clusters{}
	err := viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("unable to decode into config struct, %v", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
