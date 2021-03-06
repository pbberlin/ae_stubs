package main

import "github.com/pbberlin/tools/net/http/routes"

const AllowedAppID = "tec-news"
const DevServerPort = "8087"
const DevServerAdmin = "8002"

func init() {
	routes.InitAppHost(AllowedAppID, DevServerPort, DevServerAdmin)
}
