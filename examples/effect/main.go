package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrobinson/nanoleaf-go"
)

func main() {
	var (
		host   = flag.String("host", "", "The IP or hostname of the panel")
		port   = flag.Int("port", 16021, "The port of the panel")
		apiKey = flag.String("apiKey", "", "The API key of the panel")

		name = flag.String("name", "", "The name of the effect to operate on")
	)
	flag.Parse()

	c := nanoleaf.NewClient(&http.Client{}, *host, *port, *apiKey)

	if len(*name) > 0 {
		effect, err := c.GetEffect(context.Background(), *name)
		if err != nil {
			fmt.Printf("err getting effect: %s\n", err.Error())
			return
		}

		spew.Dump(effect)
		return
	}

	effects, err := c.GetEffects(context.Background())
	if err != nil {
		fmt.Printf("error getting effects: %s\n", err.Error())
		return
	}
	spew.Dump(effects)
}
