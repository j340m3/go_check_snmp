package cmd

import (
	g "github.com/gosnmp/gosnmp"
	"github.com/spf13/cobra"
)

var (
	warn string
	crit string
)

func SetupSNMP(cmd *cobra.Command) {
	hostname := cmd.Flags().Lookup("hostname").Value.String()
	community := cmd.Flags().Lookup("community").Value.String()
	if community != "" {
		g.Default.Community = community
		g.Default.Version = g.Version2c
	} else {
		authprot := cmd.Flags().Lookup("authprot").Value.String()
		authpass := cmd.Flags().Lookup("authpass").Value.String()
		username := cmd.Flags().Lookup("username").Value.String()
		privprot := cmd.Flags().Lookup("privprot").Value.String()
		privpass := cmd.Flags().Lookup("privpass").Value.String()

		g.Default.SecurityParameters = &g.UsmSecurityParameters{
			UserName:                 username,
			AuthenticationProtocol:   getAuthProto(authprot),
			PrivacyProtocol:          getPrivProto(privprot),
			AuthenticationPassphrase: authpass,
			PrivacyPassphrase:        privpass,
		}
		g.Default.Version = g.Version3
		g.Default.MsgFlags = g.AuthPriv
		g.Default.SecurityModel = g.UserSecurityModel
	}
	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = hostname
}

func getPrivProto(privprot string) g.SnmpV3PrivProtocol {
	switch privprot {
	case "AES":
		return g.AES

	case "DES":
		return g.DES
	}
	return g.NoPriv
}

func getAuthProto(authprot string) g.SnmpV3AuthProtocol {
	switch authprot {
	case "MD5":
		return g.MD5
	case "SHA":
		return g.SHA
	}
	return g.NoAuth
}
