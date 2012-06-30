package main

// https://www.unbound.net/documentation/libunbound-tutorial-6.html

import (
	"dns"
	"fmt"
	"os"
	"unbound"
)

func main() {
	u := unbound.New()

	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	if err := u.Hosts("/etc/hosts"); err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	if err := u.AddTaFile("keys"); err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	r, err := u.Resolve("nlnetlabs.nl.", dns.TypeA, dns.ClassINET)
	if err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	// show first result
	if r.HaveData {
		fmt.Printf("The address is %v\n", r.Data[0])
		// show security status
		if r.Secure {
			fmt.Printf("Result is secure\n")
		} else if r.Bogus {
			fmt.Printf("Result is bogus: %s\n", r.WhyBogus)
			fmt.Printf("Result is insecure\n")
		}
	}

	u.Delete()
}
