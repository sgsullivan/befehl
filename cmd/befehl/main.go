package main

import (
	"flag"
	"fmt"
	"github.com/sgsullivan/befehl"
	"os"
)

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	targets := flag.String("targets", "REQUIRED", "The location to the txt file with ips of target machines run the payload on")
	payload := flag.String("payload", "REQUIRED", "The location to the payload shell script")
	routines := flag.Int("routines", 30, "Number of allowed concurrent Go routines")

	flag.Parse()
	if *targets == "REQUIRED" {
		showUsage()
	}
	if *payload == "REQUIRED" {
		showUsage()
	}

	fmt.Printf("\nusing targets: [%s] payload: [%s] routines: [%d]\n\n", *targets, *payload, *routines)
	befehl.Fire(targets, payload, routines)
}
