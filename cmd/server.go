/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/steviesama/nx/service/netfile"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "netfile server listens for client requests.",
	Long:  `netfile server launches a network server that listens for incoming netfile client store/fetch requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, portErr := cmd.Flags().GetInt("port")
		if portErr != nil {
			return errors.New(fmt.Sprintf("netfile server --port error: %s\n", portErr.Error()))
		}

		serverErr := netfile.ServerListen("localhost", port)

		return serverErr
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	// Setup flags
	serverCmd.Flags().IntP("port", "p", 8010, `Sets the port number for to listen on.`)
}
