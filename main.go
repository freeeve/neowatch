package main

import (
	. "./neowatch"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var neopath string
	flag.StringVar(&neopath, "path", "", "path to the neo4j store")
	flag.Parse()
	if neopath == "" {
		flag.Usage()
		os.Exit(1)
	}
	fmt.Println("neowatch 0.1.0 started")
	fmt.Println("path:", neopath)

	ch := make(chan string, 100)
	go NodeStoreWatch(neopath, ch)

	time.Sleep(time.Second * 100)
}
