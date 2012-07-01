package main

// https://www.unbound.net/documentation/libunbound-tutorial-5.html

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
		fmt.Printf("error %s\n", err.Error())
		os.Exit(1)
	}

	if err := u.Hosts("/etc/hosts"); err != nil {
		fmt.Printf("error %s\n", err.Error())
		os.Exit(1)
	}
	q := make(chan bool)
	// start two goroutine
	go goroutineOne(u, q)
	go goroutineTwo(u, q)
	// wait for both routines to complete
	<-q
	<-q
}

func goroutineOne(u *unbound.Unbound, quit chan bool) {
	r, err := u.Resolve("www.nlnetlabs.nl.", dns.TypeA, dns.ClassINET)
	if err != nil {
		fmt.Printf("resolve error %s\n", err.Error())
		quit <- true
		return
	}
	if r.HaveData {
		fmt.Printf("Routine1: The address of %s is %v\n", r.Qname, r.Data[0])
	}
	quit <- true
	return
}

func goroutineTwo(u *unbound.Unbound, quit chan bool) {
	r, err := u.Resolve("www.google.nl.", dns.TypeA, dns.ClassINET)
	if err != nil {
		fmt.Printf("resolve error %s\n", err.Error())
		quit <- true
		return
	}
	if r.HaveData {
		fmt.Printf("Routine2: The address of %s is %v\n", r.Qname, r.Data[0])
	}
	quit <- true
	return
}
