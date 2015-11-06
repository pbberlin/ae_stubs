package main

import "github.com/pbberlin/tools/net/http/routes"

const AllowedAppID = "libertarian-islands"
const DevServerPort = "8085"
const DevServerAdmin = "8000"

func init() {
	routes.InitAppHost(AllowedAppID, DevServerPort, DevServerAdmin)
}
