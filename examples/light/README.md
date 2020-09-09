# Light example

This example shows how to interact with a Nanoleaf light panel. To turn on the light, you would run:

```
$ go run main.go --host=<IP of your gateway> --apiKey=<API key of the gateway> --setState=true --isOn=true
```

You likely want to also select a scene; that can be done by specifying the name of the scene and the `--scene=<scene name>` argument.
