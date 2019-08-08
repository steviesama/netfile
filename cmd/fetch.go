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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/steviesama/nx/service/netfile"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "netfile fetch requests a file transfer.",
	Long:  `netfile fetch requests a file transfer from the provided netfile server connection.`,
	RunE: func(cmd *cobra.Command, args []string) (RunError error) {
		host, hostErr := cmd.Flags().GetString("host")
		if hostErr != nil {
			return errors.New(fmt.Sprintf("netfile fetch --host error: %s\n", hostErr.Error()))
		} else if host == "" {
			return errors.New("netfile fetch --host error: no host address provided.\n")
		}

		port, portErr := cmd.Flags().GetInt("port")
		if portErr != nil {
			return errors.New(fmt.Sprintf("netfile fetch --port error: %s\n", portErr.Error()))
		} else if port == -1 {
			return errors.New("netfile fetch --port error: no port value provided.\n")
		}

		addrString := fmt.Sprintf("%s:%d", host, port)

		server, dialErr := net.Dial("tcp", addrString)

		if dialErr != nil {
			return errors.New(fmt.Sprintf("netfile fetch server dial error: %s\n", dialErr.Error()))
		}
		defer server.Close()

		for {
			command, readErr := netfile.ReadMsg(server)

			switch {
			case readErr == io.EOF:
				fmt.Println("netfile server closed connection.")
				return nil
			case readErr != nil:
				return errors.New(fmt.Sprintf("netfile.ReadMsg() error: %s\n", readErr.Error()))
			}

			filepath := "Streets.of.Fire.1984.720p.BluRay.mkv"
			var filesize int64 = 0
			var bytesCopied int64 = 0

			switch command {
			case "server.ready":
				fmt.Println("netfile server is ready.")
				// Send command request...
				_ = netfile.SendMsg(server, "client.fetch")
				// Send expected command param
				_ = netfile.SendMsg(server, filepath)
			case "server.fetch.file":
				file, createErr := os.Create(fmt.Sprintf("%s.tmp", filepath))
				defer file.Close()
				if createErr != nil {
					_ = netfile.SendMsg(server, "client.quit")
					return errors.New(fmt.Sprintf("netfile fetch file create error: %s\n", createErr.Error()))
				}

				msg, _ := netfile.ReadMsg(server)
				filesize, _ = strconv.ParseInt(msg, 10, 64)
				rw := netfile.NewConnReadWriter(server)

				data := make([]byte, netfile.BufferSize)
				fileWriter := bufio.NewWriter(file)
				bar := pb.Full.Start64(filesize)
				barReader := bar.NewProxyReader(rw)
				for {
					n, readErr := barReader.Read(data)
					// Track bytes copied so far
					// TODO: May need to remove this...proxy reader seems to handle it.
					bytesCopied += int64(n)
					// Blank assignment to signal return values are being ignored.
					numBytes, writeErr := io.CopyN(fileWriter, bytes.NewReader(data), int64(n))

					if readErr != nil {
						fmt.Printf("numBytes: %d;\nreadErr: %+v\nwriteErr: %+v\n", numBytes, readErr, writeErr)
						if readErr != io.EOF {
							RunError = errors.New(fmt.Sprintf("fetchCmd: netfile fetch file read error: %s\n", readErr.Error()))
						}
						goto ExitFor
					}
				}
				// break inside the for seems to be breaking completely out of the switch instead...
			ExitFor:
				_ = os.Rename(fmt.Sprintf("%s.tmp", filepath), filepath)
				bar.Finish()
				_ = netfile.SendMsg(server, "client.quit")
			case "server.fetch.nofile":
				// Might not need this.
				_ = netfile.SendMsg(server, "client.quit")
				return errors.New(fmt.Sprintf("'%s' does not exist...exiting.", filepath))
			}

			if command == "server.quit" {
				fmt.Println("netfile fetch...server.quit received...")
				break
			}

		}

		return RunError
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	// Setup flags
	fetchCmd.Flags().StringP("host", "H", "", `Sets the host address to connect to.`)
	fetchCmd.Flags().IntP("port", "p", -1, `Sets the port number to connect to on the host.`)
}
