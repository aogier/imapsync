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
	//	"crypto/tls"
	"fmt"
	"github.com/aogier/imapsync/common"
	"github.com/emersion/go-imap"
	//	"github.com/emersion/go-imap/client"
	"github.com/spf13/cobra"
	"log"
	"sync"
	//	"time"
)

type MailboxInfo struct {
	Folders []FolderInfo
}

type FolderInfo struct {
	Name        string
	Size, Count int
}

var wg sync.WaitGroup

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

		log.Printf("Getting info from %s@%v\n", user1, hostPort)

		connInfo := common.ConnectInfo{
			Tls: tls1, Ssl: ssl1,
			Host: host1, Port: port1,
			User: user1, Pass: pass1,
		}

		c, err := common.Connection(&connInfo)
		if err != nil {
			log.Print("Error connecting")
			log.Fatal(err)
		}
		defer c.Logout()

		// List mailboxes
		mailboxes := make(chan *imap.MailboxInfo, 10)
		done := make(chan error, 1)

		//XXX: does not work because (I can't understand) library design :P
		go func() {
			if folders != nil {
				for _, folder := range folders {
					done <- c.List("", folder, mailboxes)
				}
			}
			if foldersRec != nil {
				for _, folderRec := range foldersRec {
					done <- c.List("", fmt.Sprintf("%v*", folderRec), mailboxes)
				}
			}

			if folders == nil && foldersRec == nil {
				done <- c.List("", "*", mailboxes)
			}

		}()

		//		if err := <-done; err != nil {
		//			log.Fatal(err)
		//		}

		for w := 1; w <= pool1; w++ {
			wg.Add(1)
			go func(
				id int,
				connInfo *common.ConnectInfo,
				mailboxes <-chan *imap.MailboxInfo) {

				defer wg.Done()
				c, err := common.Connection(connInfo)
				if err != nil {
					log.Print("Error connecting")
					log.Fatal(err)
				}
				defer c.Logout()

				for m := range mailboxes {
					log.Printf("[worker %v]: %s", id, m)

					//					time.Sleep(1 * time.Second)

					mbox, err := c.Select(m.Name, false)
					if err != nil {
						log.Println("CANE")
						log.Fatal(err)
					}

					log.Printf("Flags for %v: %v\n", m.Name, mbox.Flags)

					//////////////////////

					// Get the last 4 messages
					//					from := uint32(1)
					//					to := mbox.Messages

					log.Printf("messages: %v\n", mbox.Messages)

					//					if mbox.Messages > 3 {
					//						// We're using unsigned integers here, only substract if the result is > 0
					//						from = mbox.Messages - 3
					//					}
					//					seqset := new(imap.SeqSet)
					//					seqset.AddRange(from, to)
					//
					//					messages := make(chan *imap.Message, 10)
					//					done = make(chan error, 1)
					//					go func() {
					//						done <- c.Fetch(seqset, []string{imap.EnvelopeMsgAttr}, messages)
					//					}()
					//
					//					log.Println("Last 4 messages:")
					//					for msg := range messages {
					//						log.Println("* " + msg.Envelope.Subject)
					//					}
					//
					//					if err := <-done; err != nil {
					//						log.Fatal(err)
					//					}

				}
			}(w, &connInfo, mailboxes)
		}

		wg.Wait()

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
