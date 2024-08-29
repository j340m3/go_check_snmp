package cmd

import (
	"fmt"
	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"regexp"
)

var descriptionCmd = &cobra.Command{
	Use: "description",
	//Aliases: []string{"notification"},
	Short:        "display the snmp description",
	Long:         `Show the snmp description`,
	SilenceUsage: true,
	RunE:         getDescription,
}

func init() {
	descriptionCmd.Flags().StringP("validate", "R", "", "regexp to validate the snmp description")
	RootCmd.AddCommand(descriptionCmd)
}

func getDescription(cmd *cobra.Command, _ []string) error {
	log.Debug("description called")

	validate, _ := cmd.Flags().GetString("validate")
	if validate != "" {
		log.Debugf("Validating description: %s", validate)
		_, err := regexp.Compile(validate)
		if err != nil {
			NagiosResult("UNKNOWN", "Regex isn't valid", "", nil)
			return err
		}
	}
	SetupSNMP(cmd)
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.1.0"}
	result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		NagiosResult("UNKNOWN", err2.Error(), "", nil)
		return err2
	}
	NagiosResult("OK", fmt.Sprintf("Description is: %s", result.Variables[0].Value), "", nil)
	return nil
}
