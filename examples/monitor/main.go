package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/r3labs/sse"
	"github.com/rmrobinson/nanoleaf-go"
)

func main() {
	var (
		host   = flag.String("host", "", "The IP or hostname of the panel")
		port   = flag.Int("port", 16021, "The port of the panel")
		apiKey = flag.String("apiKey", "", "The API key of the panel")
	)
	flag.Parse()

	client := sse.NewClient("http://" + *host + ":" + strconv.Itoa(*port) + "/api/v1/" + *apiKey + "/events?id=1,3")

	client.SubscribeRaw(func(msg *sse.Event) {
		id, err := strconv.Atoi(string(msg.ID))
		if err != nil {
			fmt.Printf("error converting ID to string: %s\n", err.Error())
			return
		}

		update := &nanoleaf.PanelUpdate{
			TypeID: id,
		}

		err = json.Unmarshal(msg.Data, update)
		if err != nil {
			fmt.Printf("error unmarshaling: %s\n", err.Error())
			return
		}

		spew.Dump(update)
	})

}
