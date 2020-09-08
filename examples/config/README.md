# Config example

This example shows how to configure the Nanoleaf; creating or deleting API keys and retrieving state. An example way to run this command would be to execute:

```
$ go run main.go --host=<IP of your Nanoleaf> --apiKey=<API key of the Nanoleaf>
```

It is possible to use the --discovery=true argument to this tool in order to discover where the panel is located via UPnP. The type to look for is `nanoleaf_aurora`, and the `location` field will have both the IP address and port to supply for subsequent operations.

Creating an API key will be allowed by the panel by first pressing and holding the Power button for ~5 seconds - the 2 status LEDs on the controller will begin flashing in an alternating pattern, at this point the CreateAPIKey call will succeed.

If necessary, it is possible to get both the IP by following the documentation [here](https://forum.nanoleaf.me/docs/openapi).
