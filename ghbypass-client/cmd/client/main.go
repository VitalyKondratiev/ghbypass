package main

import (
	"flag"
	pwebsocket "ghbypass-client/pkg/websocket"
	"log"
)

var subdomain string
var proxy string
var expose string

func main() {
	flag.StringVar(&subdomain, "subdomain", "", "Subdomain on proxy host")
	flag.StringVar(&proxy, "proxy", "", "Proxy address")
	flag.StringVar(&expose, "expose", "", "URL on local machine to expose")
	flag.Parse()

	if subdomain == "" || proxy == "" || expose == "" {
		log.Fatal("--subdomain, --proxy and --expose options are required")
	}

	log.Println("Subdomain set to:", subdomain)
	log.Println("Local address for expose:", expose)

	conn := pwebsocket.GetWebsocketConnection(proxy, subdomain)
	defer conn.Close()

	pwebsocket.HandleWebsocketRequests(conn, expose)
}
