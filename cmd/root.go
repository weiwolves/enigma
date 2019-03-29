package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var dbType string
var dbHost string
var dbPort int32
var dbUser string
var dbPass string
var dbName string
var dbTable string
var protoFile string
var caseFlag string
var goPkg string

func init() {
	genCmd.Flags().StringVarP(&dbType, "type", "g", "mysql", "db type: mysql, postgres")
	genCmd.Flags().StringVarP(&dbHost, "host", "o", "localhost", "db host")
	genCmd.Flags().Int32VarP(&dbPort, "port", "p", 3306, "db port")
	genCmd.Flags().StringVarP(&dbUser, "user", "u", "root", "db user")
	genCmd.Flags().StringVarP(&dbPass, "pass", "w", "secret", "db pass")
	genCmd.Flags().StringVarP(&dbName, "db", "d", "db_name", "db name")
	genCmd.Flags().StringVarP(&dbTable, "table", "t", "table_name", "db name")
	genCmd.Flags().StringVarP(&protoFile, "file", "f", "proto", "proto file name")
	genCmd.Flags().StringVarP(&caseFlag, "case", "c", "0", "case:0, default;1 CamelCase")
	genCmd.Flags().StringVarP(&goPkg, "gopkg", "k", "", "golang package")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(genCmd)
}

var rootCmd = &cobra.Command{
	Use:   "enigma",
	Short: "Enigma is small tool to generate protobuf file from MySQL",
	Long: `Enigma is small tool to generate protobuf file from MySQL,
You don't need to create proto file manual`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.Join(args, " "))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
