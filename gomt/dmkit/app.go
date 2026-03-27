package main

import (
	"fmt"
	"gomt/dmkit/cmd"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AppShortName   = "dmkit"
	AppName        = "DAS Data Model Kit"
	AppVersion     = ""
	AppBuild       = ""
	AppDescription = "DAS Data Model Kit"
)

var (
	Verbose bool
)

var rootCommand = &cobra.Command{
	Use:     fmt.Sprintf("%s", AppShortName),
	Version: fmt.Sprintf("%s (Build: %s)", AppVersion, AppBuild),
	Short:   AppShortName,
	Long:    AppDescription,
	Run:     run,
}

func init() {
	rootCommand.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose")
	rootCommand.CompletionOptions.DisableDefaultCmd = true

	rootCommand.AddCommand(cmd.GenerateCommand)

	viper.SetEnvPrefix(AppShortName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("logger_format", "text")
	viper.SetDefault("logger_level", "info")
	viper.SetDefault("logger_path", "./log")

	setupLoggger()
	//cmd.Help()
}

func setupLoggger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	logFormat := viper.GetString("logger_format")
	switch logFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	logLevel := viper.GetString("logger_level")
	logrus.ParseLevel(logLevel)

	if Verbose {
		logrus.SetLevel(logrus.TraceLevel)
	}

}

func run(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
