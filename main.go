package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "macchange"
	app.Usage = "Change the mac address for a network interface"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "iface, i",
			Value: "",
			Usage: "Interface to change",
		},
		cli.StringFlag{
			Name:  "mac, m",
			Value: "",
			Usage: "MAC address to use",
		},
	}
	app.Action = func(c *cli.Context) {

		rand.Seed(time.Now().Unix())

		//Root is needed to change MAC address on linux
		// if os.Geteuid() != 0 {
		// 	fmt.Println("You must run this as root!")
		// 	return
		// }

		iface := c.String("iface")
		newmac := c.String("mac")

		j := -1
		ifaces, err := net.Interfaces()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		//Look for the user given interface in the list of actual interfaces
		for i := 0; iface != "" && i < len(ifaces); i++ {
			if iface == ifaces[i].Name {
				j = i
				break
			}
		}

		//The user given interfaces was not found, so prompt the user to pick a valid one
		if j < 0 {
			fmt.Println("Select a network interface:")

			for i := range ifaces {
				fmt.Printf("%d: %s\n", i, ifaces[i].Name)
			}

			for j < 0 || j >= len(ifaces) {
				fmt.Scanf("%d", &j)
			}

			iface = ifaces[j].Name
		}

		//If the user doesn't provide a mac address, we'll make one for them
		if newmac == "" {
			//0x00, 0x05, 0x69 is the start to a vmware MAC address
			digits := []interface{}{0x00, 0x05, 0x69, rand.Intn(0x7f), rand.Intn(0xff), rand.Intn(0xff)}

			newmac = fmt.Sprintf("%.2x:%.2x:%.2x:%.2x:%.2x:%.2x", digits...)
		}

		//Make sure the new MAC address is a valid one
		if _, err := net.ParseMAC(newmac); err != nil {
			fmt.Printf("Invalid MAC address: %s", newmac)
			return
		}

		fmt.Printf("Setting %s MAC address to %s\n", iface, newmac)

		// ip("link", "set", "down", "dev", iface)
		// ip("link", "set", "dev", iface, "address", newmac)
		// ip("link", "set", "up", "dev", iface)
	}

	app.Run(os.Args)
}

//ip wraps an iproute2 shell command and exits on any errors
func ip(params ...string) string {
	output, err := exec.Command("ip", params...).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return string(output)
}
