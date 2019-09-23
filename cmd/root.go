package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/signavio/aws-mfa-login/action"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var cfgFile string
var Name string

var rootCmd = &cobra.Command{
	Use:   "aws-mfa-login",
	Short: "aws login with mfa",
	Long:  "CLI tool to update your temporary AWS credentials ",
	Run: func(cmd *cobra.Command, args []string) {
		PrintConfig()
		action.UpdateSessionCredentials()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aws-mfa.yaml)")
	rootCmd.Flags().StringVarP(&Name, "source", "s", "", "source profile where mfa is activated")
	rootCmd.Flags().StringVarP(&Name, "destination", "d", "", "destination profile for temporary aws credentials")
	viper.BindPFlag("source", rootCmd.Flags().Lookup("source"))
	viper.BindPFlag("destination", rootCmd.Flags().Lookup("destination"))
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
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func PrintConfig() {
	fmt.Println("Current Config located in ~/.aws-mfa.yaml\n#####")
	for _, key := range viper.AllKeys() {
		fmt.Printf("%v: %v\n", key, viper.Get(key))
	}
	fmt.Print("\n")
}
