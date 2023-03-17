package main

// https://unbound.docs.nlnetlabs.nl/en/latest/developer/libunbound-tutorial/setup-context.html

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
	"github.com/miekg/unbound"
)

func main() {
	u := unbound.New()
	defer u.Destroy()

	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		log.Fatalf("error %s\n", err.Error())
	}

	if err := u.Hosts("/etc/hosts"); err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	r, err := u.Resolve("www.nlnetlabs.nl.", dns.TypeA, dns.ClassINET)
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}
	fmt.Printf("%+v\n", r)
}
