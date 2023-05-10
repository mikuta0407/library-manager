/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mikuta0407/library-manager/internal/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start api server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		server.HandleRequests(httpHost, httpPort, sqliteDBPath)
	},
}

var httpHost string
var httpPort string
var sqliteDBPath string

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&httpHost, "httphost", "l", "0.0.0.0", "Option: HTTP Host (defaul: 0.0.0.0)")
	serverCmd.Flags().StringVarP(&httpPort, "httpport", "p", "8080", "Option: HTTP Port (defaul: 8080)")
	serverCmd.Flags().StringVarP(&sqliteDBPath, "dbfile", "f", "./library.db", "Option: SQLite3 DB File (defaul: ./library.db)")
}
