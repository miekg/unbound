package main

// https://www.unbound.net/documentation/libunbound-tutorial-2.html

import (
	"dns"
	"fmt"
	"os"
	"unbound"
)

func main() {
	u := unbound.New()

	err := u.ResolvConf("/etc/resolv.conf")
	if err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	err = u.Hosts("/etc/hosts")
	if err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	r, err := u.Resolve("nlnetlabs.nl.", dns.TypeA, dns.ClassINET)
	if err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%+v\n", r)
	u.Delete()
}
