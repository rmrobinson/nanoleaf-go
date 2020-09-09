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

		isOn  = flag.Bool("isOn", false, "Whether to turn the light on or off")
		scene = flag.String("scene", "", "The name of the scene to apply")
		/*
			hue  = flag.Int("hue", 0, "The hue to set")
			sat  = flag.Int("sat", 0, "The saturation level to set")
			name = flag.String("name", "", "The name of the device to set")
		*/
		orientation = flag.Int("orientation", -1, "The orientation to set")

		setState  = flag.Bool("setState", false, "Whether to set the specified state fields")
		setConfig = flag.Bool("setConfig", false, "Whether to set the specified config fields")
	)
	flag.Parse()

	c := nanoleaf.NewClient(&http.Client{}, *host, *port, *apiKey)

	if *setState {
		err := c.SetOn(context.Background(), *isOn)
		if err != nil {
			fmt.Printf("error setting light on: %s\n", err.Error())
			return
		}
		fmt.Printf("set complete\n")

		if len(*scene) > 0 {
			err := c.SetScene(context.Background(), *scene)
			if err != nil {
				fmt.Printf("error setting light scene: %s\n", err.Error())
				return
			}
			fmt.Printf("scene selected\n")
		}
	}
	if *setConfig {
		if *orientation >= 0 {
			err := c.SetOrientation(context.Background(), *orientation)
			if err != nil {
				fmt.Printf("error setting orientation: %s\n", err.Error())
				return
			}
			fmt.Printf("orientation set\n")
		}
		fmt.Printf("set complete\n")
	}

	resp, err := c.GetPanel(context.Background())
	if err != nil {
		fmt.Printf("error getting light: %s\n", err.Error())
		return
	}
	spew.Dump(resp)
}
