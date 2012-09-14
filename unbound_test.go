package unbound

import (
	"fmt"

//	"testing"
)

func ExampleLookupCNAME() {
	u := New()
	defer u.Destroy()
	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		return
	}
	s, err := u.LookupCNAME("www.miek.nl.")
	// A en AAAA lookup get canoncal name
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", s)
}

func ExampleLookupIP() {
	u := New()
	defer u.Destroy()
	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		return
	}
	a, err := u.LookupIP("nlnetlabs.nl.")
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", a)
}
