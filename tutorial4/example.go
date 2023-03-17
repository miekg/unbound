package main

// https://unbound.docs.nlnetlabs.nl/en/latest/developer/libunbound-tutorial/async-lookup.html

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
	"github.com/miekg/unbound"
)

// This is called when resolution is completed.
func mycallback(m interface{}, e error, r *unbound.Result) {
	done := m.(*int)
	*done = 1
	if e != nil {
		fmt.Printf("resolve error: %s\n", e.Error())
		return
	}
	if r.HaveData {
		fmt.Printf("The address of %s is %v\n", r.Qname, r.Data[0])
	}

}

func main() {
	u := unbound.New()
	defer u.Destroy()
	done := make(chan *unbound.ResultError)

	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		log.Fatalf("error %s\n", err.Error())
	}

	if err := u.Hosts("/etc/hosts"); err != nil {
		log.Fatalf("error %s\n", err.Error())
	}

	u.ResolveAsync("www.nlnetlabs.nl.", dns.TypeA, dns.ClassINET, done)
For:
	for {
		select {
		case r := <-done:
			if r.Error != nil {
				fmt.Printf("resolve error: %s\n", r.Error.Error())
				break For
			}
			if r.Result.HaveData {
				fmt.Printf("The address of %s is %v\n", r.Result.Qname, r.Result.Data[0])
				break For
			}

		}
	}
	fmt.Println("done")
}
