package main

import (
	"flag"
	"fmt"
	"github.com/sgsullivan/befehl"
	"os"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	targets := flag.String("targets", pwd+"/targets.txt", "The location to the txt file with ips of target machines run the payload on")
	payload := flag.String("payload", pwd+"/payload.sh", "The location to the payload shell script")
	routines := flag.Int("routines", 30, "Number of allowed concurrent Go routines")

	flag.Parse()

	fmt.Printf("\nusing targets: [%s] payload: [%s] routines: [%d]\n\n", *targets, *payload, *routines)
	befehl.Fire(targets, payload, routines)
}
