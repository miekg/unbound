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
	defer u.Destroy()

	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	if err = u.Hosts("/etc/hosts"); err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}

	r, err := u.Resolve("www.nlnetlabs.nl.", dns.TypeA, dns.ClassINET)
	if err != nil {
		fmt.Printf("error %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%+v\n", r)
}
