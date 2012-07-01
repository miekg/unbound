package main

// https://www.unbound.net/documentation/libunbound-tutorial-2.html

import (
	"dns"
	"flag"
	"fmt"
	"os"
	"unbound"
)

// Examine the result structure in detail
func examineResult(query string, r *unbound.Result) {
	fmt.Printf("The query is for: %s\n", query)
	fmt.Printf("The result has:\n")
	fmt.Printf("qname: %s\n", r.Qname)
	fmt.Printf("qtype: %d\n", r.Qtype)
	fmt.Printf("qclass: %d\n", r.Qclass)
	if r.CanonName != "" {
		fmt.Printf("canonical name: %s\n", r.CanonName)
	} else {
		fmt.Printf("canonical name: <none>\n")
	}

	if r.HaveData {
		fmt.Printf("has data\n")
	} else {
		fmt.Printf("has no data\n")
	}

	if r.NxDomain {
		fmt.Printf("nxdomain (name does not exist)\n")
	} else {
		fmt.Printf("not an nxdomain (name exists)\n")
	}

	if r.Secure {
		fmt.Printf("validated to be secure\n")
	} else {
		fmt.Printf("not validated as secure\n")
	}

	if r.Bogus {
		fmt.Printf("a security failure! (bogus)\n")
	} else {
		fmt.Printf("not a security failure (not bogus)\n")
	}

	fmt.Printf("DNS rcode: %d\n\n", r.Rcode)

	for i, d := range r.Data {
		fmt.Printf("result data element %d has length %d\n", i, len(d))
		fmt.Printf("result data element %d is: %v\n", i, d)
		fmt.Printf("result data element as RR: %s\n", r.Rr[i])
	}
	fmt.Printf("result has %d data element(s)\n", len(r.Data))
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("usage: <hostname>")
		os.Exit(1)
	}

	u := unbound.New()
	defer u.Destroy()

	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		fmt.Printf("error reading resolv.conf %s\n", err.Error())
		os.Exit(1)
	}

	if err := u.Hosts("/etc/hosts"); err != nil {
		fmt.Printf("error reading hosts: %s\n", err.Error())
		os.Exit(1)
	}

	r, err := u.Resolve(flag.Arg(0), dns.TypeA, dns.ClassINET)
	if err != nil {
		fmt.Printf("resolve error %s\n", err.Error())
		os.Exit(1)
	}
	examineResult(flag.Arg(0), r)
}
