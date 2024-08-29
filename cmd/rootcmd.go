// Package cmd commands
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"strings"
	"time"
)

var (
	debugFlag       = false
	unitTestFlag    = false
	hostname        string
	snmpCommunity   string
	snmpUsername    string
	snmpAuthAlg     string
	snmpAuthPasswd  string
	snmpPrivAlg     string
	snmpPrivPasswd  string
	hmToken         string
	hmWarnThreshold string
	hmCritThreshold string

	// RootCmd entry point to start
	RootCmd = &cobra.Command{
		Use:   "snmpcli",
		Short: "snmpcli – SNMP Command Line and Icinga compatible Monitoring Tool",
		Long:  `Query Tool and Nagios/Icinga check plugin for SNMP-Hosts`,
		//SilenceErrors: true,
	}
)

const (
	// allows you to override any config values using
	// env APP_MY_VAR = "MY_VALUE"
	// e.g. export APP_LDAP_USERNAME test
	// maps to ldap.username
	configEnvPrefix = "SNMPCLI"
	configName      = "snmpcli"
	configType      = "yaml"
)

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "verbose debug output")
	RootCmd.PersistentFlags().BoolVarP(&unitTestFlag, "unit-test", "", false, "redirect output for unit tests")
	RootCmd.PersistentFlags().StringVarP(&hmToken, "token", "t", "", "Homematic XMLAPI Token")
	RootCmd.PersistentFlags().StringVarP(&hmWarnThreshold, "warn", "w", "", "warning level")
	RootCmd.PersistentFlags().StringVarP(&hmCritThreshold, "crit", "c", "", "critical level")
	RootCmd.PersistentFlags().StringVarP(&hostname, "hostname", "H", "", "hostname")
	RootCmd.PersistentFlags().StringVarP(&snmpCommunity, "community", "C", "", "SNMP Community")
	RootCmd.PersistentFlags().StringVarP(&snmpUsername, "username", "u", "", "SNMP Username")
	RootCmd.PersistentFlags().StringVarP(&snmpAuthAlg, "authprot", "a", "", "Authentication algorithm")
	RootCmd.PersistentFlags().StringVarP(&snmpAuthPasswd, "authpass", "A", "", "Authentication password")
	RootCmd.PersistentFlags().StringVarP(&snmpPrivAlg, "privprot", "x", "", "Privacy algorithm")
	RootCmd.PersistentFlags().StringVarP(&snmpPrivPasswd, "privpass", "X", "", "Privacy password")
	// don't have variables populated here
	if err := viper.BindPFlags(RootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}
}

// Execute run application
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		p := GetPlugin()
		p.Errors = append(p.Errors, err)
		log.Debugf("return UNKNOWN, errors: %v", p.Errors)
		NagiosResult("UNKNOWN", "", "", nil)
	}
}

func initConfig() {
	viper.SetConfigType(configType)
	viper.SetConfigName(configName)
	//home := homedir.Get(

	// env var overrides
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix(configEnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// check flags
	processFlags()

	// logger settings
	log.SetLevel(log.ErrorLevel)
	if debugFlag {
		// report function name
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
		//hmlib.SetDebug(true)
	}

	logFormatter := &prefixed.TextFormatter{
		DisableColors:   unitTestFlag,
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
	}
	log.SetFormatter(logFormatter)

	if unitTestFlag {
		log.SetOutput(RootCmd.OutOrStdout())
	}
	// debug config file

	// validate method
}

// processConfig reads in config file and ENV variables if set.
func processConfig() (bool, error) {
	err := viper.ReadInConfig()
	haveConfig := false
	if err == nil {
		//cfgFile = viper.ConfigFileUsed()
		haveConfig = true
	}
	return haveConfig, err
}

func processFlags() {
	//if common.CmdFlagChanged(RootCmd, "debug") {
	//	viper.Set("debug", debugFlag)
	//}
	//if common.CmdFlagChanged(RootCmd, "token") {
	//	viper.Set("token", hmToken)
	//}
	//if common.CmdFlagChanged(RootCmd, "url") {
	//	viper.Set("url", hmURL)
	//}
	debugFlag = viper.GetBool("debug")
	hmToken = viper.GetString("token")
	//hmlib.SetHmToken(hmToken)
	//hmlib.SetHmURL(hmURL)
}
