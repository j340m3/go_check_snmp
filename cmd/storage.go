package cmd

import (
	"fmt"
	"github.com/atc0005/go-nagios"
	"github.com/dustin/go-humanize"
	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"reflect"
	"time"
)

var storageCmd = &cobra.Command{
	Use: "uptime",
	//Aliases: []string{"notification"},
	Short:        "display the snmp uptime in seconds",
	Long:         `Show the SNMP uptime in seconds`,
	SilenceUsage: true,
	RunE:         getStorage,
}

// SNMP Datas

const storage_table= "1.3.6.1.2.1.25.2.3.1"
const storagetype_table = "1.3.6.1.2.1.25.2.3.1.2"
const index_table = "1.3.6.1.2.1.25.2.3.1.1"
const descr_table = "1.3.6.1.2.1.25.2.3.1.3"
const size_table = "1.3.6.1.2.1.25.2.3.1.5.";
const used_table = "1.3.6.1.2.1.25.2.3.1.6."
const alloc_units = "1.3.6.1.2.1.25.2.3.1.4."

//Storage types definition  - from /usr/share/snmp/mibs/HOST-RESOURCES-TYPES.txt
var hrStorage = map[string]string{
	"Other" : "1.3.6.1.2.1.25.2.1.1",
	"1.3.6.1.2.1.25.2.1.1" : "Other",
	"Ram" : "1.3.6.1.2.1.25.2.1.2",
	"1.3.6.1.2.1.25.2.1.2" : "Ram",
	"VirtualMemory" : "1.3.6.1.2.1.25.2.1.3",
	"1.3.6.1.2.1.25.2.1.3" : "VirtualMemory",
	"FixedDisk" : "1.3.6.1.2.1.25.2.1.4",
	"1.3.6.1.2.1.25.2.1.4" : "FixedDisk",
	"RemovableDisk" : "1.3.6.1.2.1.25.2.1.5",
}



$hrStorage{"1.3.6.1.2.1.25.2.1.5"} = 'RemovableDisk';
$hrStorage{"FloppyDisk"} = '1.3.6.1.2.1.25.2.1.6';
$hrStorage{"1.3.6.1.2.1.25.2.1.6"} = 'FloppyDisk';
$hrStorage{"CompactDisk"} = '1.3.6.1.2.1.25.2.1.7';
$hrStorage{"1.3.6.1.2.1.25.2.1.7"} = 'CompactDisk';
$hrStorage{"RamDisk"} = '1.3.6.1.2.1.25.2.1.8';
$hrStorage{"1.3.6.1.2.1.25.2.1.8"} = 'RamDisk';
$hrStorage{"FlashMemory"} = '1.3.6.1.2.1.25.2.1.9';
$hrStorage{"1.3.6.1.2.1.25.2.1.9"} = 'FlashMemory';
$hrStorage{"NetworkDisk"} = '1.3.6.1.2.1.25.2.1.10';
$hrStorage{"1.3.6.1.2.1.25.2.1.10"} = 'NetworkDisk';

func init() {
	//descriptionCmd.Flags().StringP("validate", "R", "", "regexp to validate the snmp description")
	descriptionCmd.Flags().StringVarP(&warn, "warn", "w", "", "warning threshold")
	descriptionCmd.Flags().StringVarP(&crit, "crit", "c", "", "critical threshold")
	RootCmd.AddCommand(uptimeCmd)
}

func getStorage(cmd *cobra.Command, _ []string) error {
	log.Debug("uptime called")

	warn, _ := cmd.Flags().GetString("warn")
	crit, _ := cmd.Flags().GetString("crit")
	warnRange := nagios.ParseRangeString(warn)
	critRange := nagios.ParseRangeString(crit)

	SetupSNMP(cmd)
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.3.0"}
	result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		NagiosResult("UNKNOWN", err2.Error(), "", nil)
		return err2
	}

	uptimeValue := g.ToBigInt(result.Variables[0].Value)
	fmt.Println(reflect.TypeOf(result.Variables[0].Value).String())
	//uptimeAsString := strconv.Itoa(int(uptimeValue))
	uptimeAsString := uptimeValue.String() + "0"
	//uptimeInSeconds := uptimeValue / 1000
	uptimeAsTime := time.Now().Add(-10 * time.Duration(uptimeValue.Int64()) * time.Millisecond)
	p := GetPlugin()
	perfD := nagios.PerformanceData{
		Label:             "uptime",
		Value:             uptimeAsString,
		UnitOfMeasurement: "ms",
		Warn:              warn,
		Crit:              crit,
		Min:               "0",
		Max:               "",
	}
	_ = p.AddPerfData(true, perfD)
	text := fmt.Sprintf("Booted %s", humanize.Time(uptimeAsTime))
	if critRange.CheckRange(uptimeAsString) {
		NagiosResult("CRITICAL", "CRITICAL: "+text, "", nil)
	} else if warnRange.CheckRange(uptimeAsString) {
		NagiosResult("WARNING", text, "", nil)
	} else {
		NagiosResult("OK", text, "", nil)
	}
	return nil
}
