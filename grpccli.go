/*
gRPC Client
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"io/ioutil"

	xr "github.com/nleiva/xrgrpc/client"
)

func main() {
	// Encoding option; defaults to JSON
	enc := flag.String("enc", "json", "Encoding: 'json' or 'text'")
	// CLI to issue; defaults to "show grpc status"
	cli := flag.String("cli", "show grpc status", "Command to execute")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "config.json", "Configuration file")
	// YANG path arguments; defaults to "yangpaths.json"
	ypath := flag.String("ypath", "yangpaths.json", "YANG path arguments")
	// RPC to call (as defined in proto file); get-config, etc.
	rpc := flag.String("rpc", "", "RPC to call")
	flag.Parse()

	id := rand.Int63n(1000)
	output := "Empty"

	// Define target parameters from the configuration file
	target := xr.NewCiscoGrpcClient()
	err := decodeJSONConfig(target, *cfg)
	if err != nil {
		log.Fatalf("Could not read the config: %v", err)
	}

	// Setup a connection to the target
	conn, err := xr.Connect(*target)
	if err != nil {
		log.Fatalf("Could not setup a client connection to the target: %v", err)
	}
	defer conn.Close()

	// Horrible hack to easily test GetConfig
	if *rpc == "get-config" {
		js, err := ioutil.ReadFile(*ypath)
		if err != nil {
			fmt.Printf("Couldn't read file: %v: %v\n", *ypath, err)
		}
		output, err = xr.GetConfig(conn, string(js), id)
		if err != nil {
			fmt.Printf("Couldn't get the config: %v\n", err)
		}
		fmt.Print(output)
		os.Exit(0)
	}

	// Return show command output based on encoding selected
	switch *enc {
	case "text":
		output, err = xr.ShowCmdTextOutput(conn, *cli, id)
	case "json":
		output, err = xr.ShowCmdJSONOutput(conn, *cli, id)
	default:
		fmt.Printf("Don't recognize encoding: %v\n", *enc)
	}
	if err != nil {
		fmt.Printf("Couldn't get the cli output: %v\n", err)
	}
	fmt.Print(output)

}