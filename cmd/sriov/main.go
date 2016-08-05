package main

import (
	"flag"
	"fmt"
	"github.com/hustcat/go-lib/sriov"
	"log"
	"os"
)

var (
	op        = ""
	netNsPath = ""
	vfIndex   = -1
	master    = "eth1"
	ifName    = "eth1" //VF device in container net ns
	//for setup
	mac  = ""
	vlan = 0
	ip   = ""
	gw   = ""
)

func init() {
	flag.StringVar(&op, "op", "", "operation: add or del")
	flag.StringVar(&netNsPath, "ns", "", "container netns path")
	flag.IntVar(&vfIndex, "vf", -1, "vf index")
	flag.StringVar(&mac, "mac", "", "macaddress for vf device")
	flag.IntVar(&vlan, "vlan", 0, "vlan id")
	flag.StringVar(&ip, "ip", "", "ip address for vf device")
	flag.StringVar(&gw, "gw", "", "gateway address")

}

func add() {
	if netNsPath == "" {
		log.Fatal("ns can't be null")
		return
	}

	if vfIndex == -1 {
		log.Fatal("vf index can't be -1")
		return
	}

	if mac == "" {
		log.Fatal("mac address can't be null")
		return
	}

	if ip == "" {
		log.Fatal("ip address can't be null")
		return
	}

	if vlan == 0 {
		log.Fatal("vlan can't be 0")
		return
	}

	if gw == "" {
		log.Fatal("gateway address can't be null")
		return

	}

	conf := &sriov.NetConf{
		Master:  master,
		MAC:     mac,
		VF:      vfIndex,
		Vlan:    vlan,
		IPAddr:  ip,
		Gateway: gw,
		IfName:  ifName,
		NetNs:   netNsPath,
	}

	err := sriov.SetupVF(conf)
	if err != nil {
		fmt.Printf("setup VF failed: %v", err)
		os.Exit(1)
	} else {
		fmt.Printf("setup VF success")
	}

}

func del() {
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

func main() {
	flag.Parse()
	if op == "add" {
		add()
	} else if op == "del" {
		del()
	} else {
		log.Fatal("invalid operation")
	}
}
