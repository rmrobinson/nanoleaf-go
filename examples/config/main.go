package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/koron/go-ssdp"
	"github.com/rmrobinson/nanoleaf-go"
)

func main() {
	var (
		host   = flag.String("host", "", "The IP or hostname of the panels")
		port   = flag.Int("port", 16021, "The port to connect to")
		apiKey = flag.String("apiKey", "", "The API key of the panels")

		discover = flag.Bool("discover", false, "Whether to run the discovery tool")
		create   = flag.Bool("create", false, "Whether to create a new API key or not")
		delete   = flag.Bool("delete", false, "Whether to delete the other API key")
		delKey   = flag.String("delKey", "", "The API key to delete")
	)
	flag.Parse()

	if *discover {
		m := &ssdp.Monitor{
			Alive: func(msg *ssdp.AliveMessage) {
				fmt.Printf("alive; from=%s, type=%s, server=%s, location=%s\n", msg.From.String(), msg.Type, msg.Server, msg.Location)
			},
			Bye: func(msg *ssdp.ByeMessage) {
				fmt.Printf("bye; from=%s, type=%s\n", msg.From.String(), msg.Type)
			},
			Search: func(msg *ssdp.SearchMessage) {
				fmt.Printf("search; from=%s type=%s\n", msg.From.String(), msg.Type)
			},
		}

		if err := m.Start(); err != nil {
			fmt.Printf("error running ssdp: %s\n", err.Error())
			return
		}

		time.Sleep(time.Second * 90)
		return
	}

	c := nanoleaf.NewClient(&http.Client{}, *host, *port, *apiKey)

	if *create {
		key, err := c.CreateAPIKey(context.Background())
		if err != nil {
			fmt.Printf("error creating API key: %s\n", err.Error())
			return
		}

		fmt.Printf("created new API key %s\n", key)
		*apiKey = key
	}
	if *delete {
		err := c.DeleteAPIKey(context.Background(), *delKey)
		if err != nil {
			fmt.Printf("error deleting API key: %s\n", err.Error())
			return
		}

		fmt.Printf("deleted api key\n")
	}

	gw, err := c.GetPanel(context.Background())
	if err != nil {
		fmt.Printf("err getting panel: %s\n", err.Error())
		return
	}

	spew.Dump(gw)
}
