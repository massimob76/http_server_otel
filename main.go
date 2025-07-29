package main

import (
	"os"
	"strconv"
)

func main() {
	server()
}

func getIdentity() int {
	if identity, ok := os.LookupEnv("IDENTITY"); ok {
		id, err := strconv.Atoi(identity)
		if err != nil {
			log.Error("Cannot covert identity to integer", "identity", identity, "err", err.Error())
			os.Exit(1)
		}
		return id
	} else {
		return 0
	}
}
