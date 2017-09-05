package common

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-imap/client"
	"log"
)

type ConnectInfo struct {
	Tls, Ssl         bool
	Host, User, Pass string
	Port             int
}

func Connection(parms *ConnectInfo) (c *client.Client, err error) {

	hostPort := fmt.Sprintf("%s:%v", parms.Host, parms.Port)

	if parms.Ssl {
		fmt.Println("ciao")
		c, err = client.DialTLS(hostPort, nil)

		if err != nil {
			return
		}
		log.Println("Connected ssl")

	} else {
		c, err = client.Dial(hostPort)

		if err != nil {
			return nil, err
		}
		log.Println("Connected")

		have_tls, _ := c.SupportStartTLS()

		switch {
		case have_tls != true && parms.Tls:
			log.Fatal("Source host does not support enforced TLS")
		case have_tls:
			tlsConfig := &tls.Config{ServerName: parms.Host}
			if err := c.StartTLS(tlsConfig); err != nil {
				log.Fatal(err)
			}
			log.Println("TLS started")
		}
	}

	// Login
	if err := c.Login(parms.User, parms.Pass); err != nil {
		return nil, err
	}
	log.Println("Logged in")

	return

}
