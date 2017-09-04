// Copyright © 2017 Alessandro Ogier <alessandro.ogier@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/emersion/go-imap/client"
	"github.com/spf13/cobra"
	"log"
)

type MailboxInfo struct {
	Folders []FolderInfo
}

type FolderInfo struct {
	Name        string
	Size, Count int
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		hostPort := fmt.Sprintf("%s:%v", host1, port1)

		log.Printf("getting info from %s@%v\n", user1, hostPort)

		switch {
		case ssl1:
			fmt.Println("ciao")
		default:
			// Connect to server
			c, err := client.Dial(hostPort)
			if err != nil {
				log.Print("error connecting")
				log.Fatal(err)
			}
			log.Println("Connected")

			// Don't forget to logout
			defer c.Logout()

			have_tls, err := c.SupportStartTLS()

			switch {
			case have_tls != true && tls1:
				log.Fatal("source host does not support enforced TLS")
			case have_tls:
				c, err = client.DialTLS(hostPort, nil)
				if err != nil {
					log.Print("error connecting")
					log.Fatal(err)
				}
			}

			capabilities, err := c.Capability()
			log.Println(capabilities)
		}

	},
}

func init() {
	RootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
