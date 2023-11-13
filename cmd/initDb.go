/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"log"
	"os"

	"github.com/mikuta0407/library-manager/internal/database"
	"github.com/spf13/cobra"
)

// initDbCmd represents the initDb command
var initDbCmd = &cobra.Command{
	Use:   "initdb",
	Short: "Generate empty sqlite3 db file",
	Run: func(cmd *cobra.Command, args []string) {
		generateDbFile(filepath)
	},
}

var filepath string

func init() {
	rootCmd.AddCommand(initDbCmd)

	initDbCmd.Flags().StringVarP(&filepath, "filepath", "f", "./library.db", "filename of empty sqlite3 db file ")
}

var emptyLibraryDBFileBytes []byte

func generateDbFile(filepath string) {

	_, err := os.Stat(filepath)
	if err == nil {
		log.Println(filepath, "is exists. No modified.")
		return
	}

	if err := database.ConnectDB(filepath); err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(filepath, "is generated")

	defer database.DisconnectDB()
}
