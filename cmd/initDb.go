/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"log"

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

//go:embed library.db
var emptyLibraryDBFileBytes []byte

func generateDbFile(filepath string) {
	if err := database.ConnectDB; err != nil {
		log.Fatalln(err)
		return
	}

	defer database.DisconnectDB()
}
