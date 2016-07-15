package main

import "github.com/pbberlin/tools/net/http/routes"

const AllowedAppID = "dom-clean"
const DevServerPort = "8088"
const DevServerAdmin = "8008"

func init() {
	routes.InitAppHost(AllowedAppID, DevServerPort, DevServerAdmin)
}
