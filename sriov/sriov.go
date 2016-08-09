package sriov

import (
	"fmt"
	"github.com/containernetworking/cni/pkg/ip"
	"github.com/containernetworking/cni/pkg/ns"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	"net"
	"os"
	"runtime"
)

type NetConf struct {
	Master string `json:"master"`
	MAC    string `json:"mac"`
	VF     int    `json:"vf"`
	Vlan   int    `json:"vlan"`
	// ex: "192.168.1.2/24"
	IPAddr  string `json:"ip"`
	Gateway string `json:"gateway"`
	IfName  string `json:"ifname"`
	NetNs   string `json:"netns"`
}

func init() {
	runtime.LockOSThread()
}

func setupVF(conf *NetConf, ifName string, netns ns.NetNS) error {

	masterName := conf.Master
	vfIdx := conf.VF

	m, err := netlink.LinkByName(masterName)
	if err != nil {
		return fmt.Errorf("failed to lookup master %q: %v", conf.Master, err)
	}

	vfDir := fmt.Sprintf("/sys/class/net/%s/device/virtfn%d/net", masterName, vfIdx)
	if _, err := os.Lstat(vfDir); err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(vfDir)
	if err != nil {
		return err
	}

	if len(infos) != 1 {
		return fmt.Errorf("Mutiple network devices in directory %s", vfDir)
	}

	// VF NIC name
	vfDevName := infos[0].Name()
	vfDev, err := netlink.LinkByName(vfDevName)
	if err != nil {
		return fmt.Errorf("failed to lookup vf device %q: %v", vfDevName, err)
	}

	// set hardware address
	if conf.MAC != "" {
		macAddr, err := net.ParseMAC(conf.MAC)
		if err != nil {
			return err
		}
		if err = netlink.LinkSetVfHardwareAddr(m, conf.VF, macAddr); err != nil {
			return fmt.Errorf("failed to set vf %d macaddress: %v", conf.VF, err)
		}
	}

	if conf.Vlan != 0 {
		if err = netlink.LinkSetVfVlan(m, conf.VF, conf.Vlan); err != nil {
			return fmt.Errorf("failed to set vf %d vlan: %v", conf.VF, err)
		}
	}

	if err = netlink.LinkSetUp(vfDev); err != nil {
		return fmt.Errorf("failed to setup vf %d device: %v", conf.VF, err)
	}

	// move VF device to ns
	if err = netlink.LinkSetNsFd(vfDev, int(netns.Fd())); err != nil {
		return fmt.Errorf("failed to move vf %d to netns: %v", conf.VF, err)
	}

	return netns.Do(func(_ ns.NetNS) error {
		err := renameLink(vfDevName, ifName)
		if err != nil {
			return fmt.Errorf("failed to rename vf %d device %q to %q: %v", conf.VF, vfDevName, ifName, err)
		}
		return nil
	})
}

func releaseVF(conf *NetConf, ifName string, initns ns.NetNS) error {
	// get VF device
	vfDev, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("failed to lookup vf %d device %q: %v", conf.VF, ifName, err)
	}

	// device name in init netns
	index := vfDev.Attrs().Index
	devName := fmt.Sprintf("dev%d", index)

	// shutdown VF device
	if err = netlink.LinkSetDown(vfDev); err != nil {
		return fmt.Errorf("failed to down vf % device: %v", conf.VF, err)
	}

	// rename VF device
	err = renameLink(ifName, devName)
	if err != nil {
		return fmt.Errorf("failed to rename vf %d evice %q to %q: %v", conf.VF, ifName, devName, err)
	}

	// move VF device to init netns
	if err = netlink.LinkSetNsFd(vfDev, int(initns.Fd())); err != nil {
		return fmt.Errorf("failed to move vf %d to init netns: %v", conf.VF, err)
	}

	return nil
}

func SetupVF(conf *NetConf) error {
	netns, err := ns.GetNS(conf.NetNs)
	if err != nil {
		return fmt.Errorf("failed to open netns %q: %v", netns, err)
	}
	defer netns.Close()

	if err = setupVF(conf, conf.IfName, netns); err != nil {
		return err
	}

	err = netns.Do(func(_ ns.NetNS) error {
		return configureIface(conf)
	})
	if err != nil {
		return err
	}

	return nil
}

func ReleaseVF(conf *NetConf) error {
	netns, err := ns.GetNS(conf.NetNs)
	if err != nil {
		return fmt.Errorf("failed to open netns %q: %v", netns, err)
	}
	defer netns.Close()

	initns, err := ns.GetCurrentNS()
	if err != nil {
		return fmt.Errorf("failed to open init ns: %v", err)
	}
	defer initns.Close()

	err = netns.Do(func(_ ns.NetNS) error {
		return releaseVF(conf, conf.IfName, initns)
	})

	return nil
}

func renameLink(curName, newName string) error {
	link, err := netlink.LinkByName(curName)
	if err != nil {
		return err
	}

	return netlink.LinkSetName(link, newName)
}

// ConfigureIface takes the result of IPAM plugin and
// applies to the ifName interface
func configureIface(conf *NetConf) error {
	ifName := conf.IfName
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("failed to lookup %q: %v", ifName, err)
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed to set %q UP: %v", ifName, err)
	}

	i, n, err := net.ParseCIDR(conf.IPAddr)
	if err != nil {
		return fmt.Errorf("failed to parse ip address :%s", conf.IPAddr)
	}

	addr := &net.IPNet{IP: i, Mask: n.Mask}
	nlAddr := &netlink.Addr{IPNet: addr, Label: ""}
	if err = netlink.AddrAdd(link, nlAddr); err != nil {
		return fmt.Errorf("failed to add IP addr to %q: %v", ifName, err)
	}

	gw := net.ParseIP(conf.Gateway)
	if gw == nil {
		return fmt.Errorf("parse gateway: %s return nil", conf.Gateway)
	}

	if err = ip.AddDefaultRoute(gw, link); err != nil {
		// we skip over duplicate routes as we assume the first one wins
		if !os.IsExist(err) {
			return fmt.Errorf("failed to add default route via %v dev %v: %v", gw, ifName, err)
		}
	}
	return nil
}
