package main

import "github.com/pbberlin/tools/net/http/routes"

const AllowedAppID = "credit-expansion"
const DevServerPort = "8086"
const DevServerAdmin = "8001"

func init() {
	routes.InitAppHost(AllowedAppID, DevServerPort, DevServerAdmin)
}
