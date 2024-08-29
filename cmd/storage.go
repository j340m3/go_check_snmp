package cmd

import (
	"fmt"
	"github.com/atc0005/go-nagios"
	humanize "github.com/dustin/go-humanize"
	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var storageCmd = &cobra.Command{
	Use: "storage",
	//Aliases: []string{"notification"},
	Short:        "display the snmp storage in seconds",
	Long:         `Show the SNMP storage in seconds`,
	SilenceUsage: true,
	RunE:         getStorage,
}

// SNMP Datas

// const storage_table = "1.3.6.1.2.1.25.2.3.1"
// const storagetype_table = "1.3.6.1.2.1.25.2.3.1.2"
// const index_table = "1.3.6.1.2.1.25.2.3.1.1"
const descr_table = "1.3.6.1.2.1.25.2.3.1.3"

//const size_table = "1.3.6.1.2.1.25.2.3.1.5."

//const used_table = "1.3.6.1.2.1.25.2.3.1.6."
//const alloc_units = "1.3.6.1.2.1.25.2.3.1.4."

// Storage types definition  - from /usr/share/snmp/mibs/HOST-RESOURCES-TYPES.txt
var hrStorage = map[string]string{
	"Other":                 "1.3.6.1.2.1.25.2.1.1",
	"1.3.6.1.2.1.25.2.1.1":  "Other",
	"Ram":                   "1.3.6.1.2.1.25.2.1.2",
	"1.3.6.1.2.1.25.2.1.2":  "Ram",
	"VirtualMemory":         "1.3.6.1.2.1.25.2.1.3",
	"1.3.6.1.2.1.25.2.1.3":  "VirtualMemory",
	"FixedDisk":             "1.3.6.1.2.1.25.2.1.4",
	"1.3.6.1.2.1.25.2.1.4":  "FixedDisk",
	"RemovableDisk":         "1.3.6.1.2.1.25.2.1.5",
	"1.3.6.1.2.1.25.2.1.5":  "RemovableDisk",
	"FloppyDisk":            "1.3.6.1.2.1.25.2.1.6",
	"1.3.6.1.2.1.25.2.1.6":  "FloppyDisk",
	"CompactDisk":           "1.3.6.1.2.1.25.2.1.7",
	"1.3.6.1.2.1.25.2.1.7":  "CompactDisk",
	"RamDisk":               "1.3.6.1.2.1.25.2.1.8",
	"1.3.6.1.2.1.25.2.1.8":  "RamDisk",
	"FlashMemory":           "1.3.6.1.2.1.25.2.1.9",
	"1.3.6.1.2.1.25.2.1.9":  "FlashMemory",
	"NetworkDisk":           "1.3.6.1.2.1.25.2.1.10",
	"1.3.6.1.2.1.25.2.1.10": "NetworkDisk",
}

func init() {
	//descriptionCmd.Flags().StringP("validate", "R", "", "regexp to validate the snmp description")
	//descriptionCmd.Flags().StringVarP(&warn, "warn", "w", "", "warning threshold")
	//descriptionCmd.Flags().StringVarP(&crit, "crit", "c", "", "critical threshold")
	RootCmd.AddCommand(storageCmd)
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
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("Close() err: %v", err)
		}
	}(g.Default.Conn)

	oids := []string{descr_table + ".1"}
	result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		NagiosResult("UNKNOWN", err2.Error(), "", nil)
		return err2
	}
	fmt.Println("variables", result.Variables)
	variable := result.Variables[0]
	switch variable.Type {
	case g.OctetString:
		value := variable.Value.([]byte)
		if strings.Contains(strconv.Quote(string(value)), "\\x") {
			tmp := ""
			for i := 0; i < len(value); i++ {
				tmp += fmt.Sprintf("%v", value[i])
				if i != (len(value) - 1) {
					tmp += " "
				}
			}
			fmt.Printf("Hex-String: %s\n", tmp)
		} else {
			fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
		}
	default:
		// ... or often you're just interested in numeric values.
		// ToBigInt() will return the Value as a BigInt, for plugging
		// into your calculations.
		fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
	}

	if result.Variables[0].Type == g.OctetString {
		fmt.Println("its a string")
	}
	if str, ok := result.Variables[0].Value.(string); ok {
		fmt.Println(string(str))
	} else {
		fmt.Println("Invalid variable")
	}
	//uptimeValue := result.Variables[0].Value
	fmt.Println("reflect", reflect.TypeOf(result.Variables[0].Value).String())
	//uptimeAsString := strconv.Itoa(int(uptimeValue))
	uptimeAsString := "0"
	//uptimeInSeconds := uptimeValue / 1000
	uptimeAsTime := time.Now().Add(-10 * time.Duration(28537) * time.Millisecond)
	p := GetPlugin()
	perfD := nagios.PerformanceData{
		Label:             "uptime",
		Value:             "uptimeAsString",
		UnitOfMeasurement: "ms",
		Warn:              warn,
		Crit:              crit,
		Min:               "0",
		Max:               "",
	}
	_ = p.AddPerfData(true, perfD)
	text := fmt.Sprintf("Stored %s", humanize.Time(uptimeAsTime))
	if critRange.CheckRange(uptimeAsString) {
		NagiosResult("CRITICAL", "CRITICAL: "+text, "", nil)
	} else if warnRange.CheckRange(uptimeAsString) {
		NagiosResult("WARNING", text, "", nil)
	} else {
		NagiosResult("OK", text, "", nil)
	}
	return nil
}
