package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	CFGFile string

	next    func()
	rootCmd = &cobra.Command{
		Use:   "homework",
		Short: "homework service",
		Long:  `Storage for homeworks`,
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			if next != nil {
				next()
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&CFGFile, "config", "c", "", "Config file path")

	rootCmd.PersistentFlags().String("http.host", "", "API host")
	rootCmd.PersistentFlags().Int("http.port", 8585, "API port")

	rootCmd.PersistentFlags().String("db.host", "127.0.0.1", "Database host")
	rootCmd.PersistentFlags().Int("db.port", 27017, "Database port")
	rootCmd.PersistentFlags().String("db.database", "homework", "Database name")
	rootCmd.PersistentFlags().String("db.login", "", "Database login")
	rootCmd.PersistentFlags().String("db.password", "", "Database password")

	rootCmd.PersistentFlags().String("logging.output", "STDOUT", "Logging output")
	rootCmd.PersistentFlags().String("logging.level", "INFO", "Logging level")
	rootCmd.PersistentFlags().String("logging.format", "TEXT", "Logging format: TEXT or JSON")

	rootCmd.PersistentFlags().String("sentryDSN", "", "Sentry DSN")

	// Bind command flags to config variables
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	// Use config file from the flag if provided.
	if CFGFile != "" {
		viper.SetConfigFile(CFGFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Can't read config:", err)
		}
	}

	viper.SetEnvPrefix("HOMEWORK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func SetVersion(version string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("Version: %s\n", rootCmd.Version))
}

func Execute(nextFn func()) {
	next = nextFn

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).
			Fatal("Execute error")
	}

}
