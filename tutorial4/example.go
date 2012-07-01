package main

// https://www.unbound.net/documentation/libunbound-tutorial-4.html

import (
	"dns"
	"fmt"
	"os"
	"time"
	"unbound"
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
	done := 0

	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		fmt.Printf("error %s\n", err.Error())
		os.Exit(1)
	}

	if err := u.Hosts("/etc/hosts"); err != nil {
		fmt.Printf("error %s\n", err.Error())
		os.Exit(1)
	}

	err := u.ResolveAsync("www.nlnetlabs.nl.", dns.TypeA, dns.ClassINET, &done, mycallback)
	if err != nil { // Will not happen in Go's case, as the return code is always nil
		fmt.Printf("error %s\n", err.Error())
		os.Exit(1)
	}
	i := 0
	for done == 0 {
		time.Sleep(1e8) // wait 1/10 of a second
		fmt.Printf("time passed (%d) ..\n", i)
		i++
	}
	fmt.Println("done")
}
