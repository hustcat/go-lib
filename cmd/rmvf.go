package main

import (
	"fmt"
	"github.com/hustcat/go-lib/sriov"
	"os"
)

var netNsPath string

func init() {
	flag.StringVar(&netNsPath, "ns", "", "container netns path")
}

func main() {

	conf := &sriov.NetConf{
		Master: "eth1",
		VF:     1,
		IfName: "eth1",
		NetNs:  netNsPath,
	}

	err := sriov.ReleaseVF(conf)
	if err != nil {
		fmt.Printf("release VF failed: %v", err)
		os.Exit(1)
	} else {
		fmt.Printf("release VF success")
	}
}
