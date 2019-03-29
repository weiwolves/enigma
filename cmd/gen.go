package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate proto file",
	Long:  `Generate proto file from MySQL`,
	Run: func(cmd *cobra.Command, args []string) {
		GenerateProto(args)
	},
}

func GenerateProto(args []string) {

	schemaTable := "information_schema"

	// example: root:pass@tcp(host:port)/database?param=value
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, schemaTable)
	mysqlDB, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Can not connect to mysql, detail: ", err)
	}

	defer mysqlDB.Close()

	stmt, err := mysqlDB.Prepare("SELECT COLUMN_NAME, DATA_TYPE FROM COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME =?")
	if err != nil {
		log.Println("Error when prepare query, detail: ", err)
		return
	}

	rows, err := stmt.Query(dbName, dbTable)
	if err != nil {
		log.Println("Error when exec query, detail: ", err)
		return
	}

	// open output file
	fo, err := os.Create(protoFile)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			log.Println("Error when exec query, detail: ", err)
			return
		}
	}()
	// make a write buffer
	w := bufio.NewWriter(fo)

	// write message name
	if _, err := w.Write([]byte("syntax = \"proto3\";\n\n")); err != nil {
		log.Println("Error when write data to file, detail: ", err)
		return
	}

	// write golang package name
	if len(goPkg) > 0 {
		if _, err := w.Write([]byte(fmt.Sprintf("option go_package = \"%s\";\n\n", goPkg))); err != nil {
			log.Println("Error when write data to file, detail: ", err)
			return
		}
	}

	var dt string
	if caseFlag == "0" {
		dt = dbTable
	} else {
		dt = GetCamelCase(dbTable)
	}
	if _, err := w.Write([]byte(fmt.Sprintf("message %s {\n", dt))); err != nil {
		log.Println("Error when write data to file, detail: ", err)
		return
	}

	counter := 0
	for rows.Next() {
		var columnName, dataType string
		err = rows.Scan(&columnName, &dataType)
		if err != nil {
			log.Println("Error when scan rows, detail: ", err)
			return
		}
		// log.Println("ColumnName: ", columnName, ", DataType: ", dataType)
		counter++
		buf := []byte(fmt.Sprintf("\t%s %s = %d;\n", GetProtoDataType(dataType), columnName, counter))

		// write a chunk
		if _, err := w.Write(buf); err != nil {
			log.Println("Error when write data to file, detail: ", err)
			return
		}
	}

	if _, err := w.Write([]byte("}\n")); err != nil {
		log.Println("Error when write data to file, detail: ", err)
		return
	}

	if err = w.Flush(); err != nil {
		log.Println("Error when flush data to file, detail: ", err)
		return
	}
}

func GetProtoDataType(mysqlType string) string {
	switch mysqlType {
	case "char", "varchar", "longtext", "text", "mediumtext":
		return "string"
	case "tinyblob", "blob", "mediumblob", "longblob":
		return "string"
	case "smallint":
		return "int32"
	case "int", "bigint", "date", "datetime", "timestamp":
		return "int64"
	case "tinyint":
		return "bool"
	case "decimal", "float", "double", "real":
		return "double"
	default:
		return ""
	}
}

func GetCamelCase(input string) string {
	output := ""
	capNext := true
	for _, v := range input {
		if v >= 'A' && v <= 'Z' {
			output += string(v)
		}
		if v >= '0' && v <= '9' {
			output += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				output += strings.ToUpper(string(v))
			} else {
				output += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}

	return output
}
