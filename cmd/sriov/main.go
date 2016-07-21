package main

import (
	"flag"
	"fmt"
	"github.com/hustcat/go-lib/sriov"
	"log"
	"os"
)

var (
	netNsPath = ""
	vfIndex   = -1
	master    = "eth1"
	ifName    = "eth1" //VF device in container net ns
)

func init() {
	flag.StringVar(&netNsPath, "ns", "", "container netns path")
	flag.IntVar(&vfIndex, "vf", -1, "vf index")
}

func main() {
	flag.Parse()

	if netNsPath == "" {
		log.Fatal("ns can't be null")
	}

	if vfIndex == -1 {
		log.Fatal("vf index can't be -1")
	}

	conf := &sriov.NetConf{
		Master: master,
		VF:     vfIndex,
		IfName: ifName,
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
