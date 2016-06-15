package main

// https://www.unbound.net/documentation/libunbound-tutorial-1.html

import (
	"fmt"
	"log"

	"github.com/miekg/unbound"
)

func main() {
	u := unbound.New()
	defer u.Destroy()

	u.ResolvConf("/etc/resolv.conf")
	a, err := u.LookupHost("www.nlnetlabs.nl")
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	for _, a1 := range a {
		fmt.Printf("The address is %s\n", a1)
	}
}
