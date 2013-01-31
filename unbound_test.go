package unbound

import (
	"fmt"
	"testing"
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

func TestDotLess(t *testing.T) {
	u := New()
	defer u.Destroy()
	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		return
	}
	a, err := u.LookupTXT("gmail.com")
	if err != nil {
		return
	}
	for _, r := range a {
		if len(r) == 0 {
			t.Log("Failure to get the TXT from gmail.com")
			t.Fail()
		}
	}
}

func TestUnicode(t *testing.T) {
	u  := New()
	defer u.Destroy()
	if err := u.ResolvConf("/etc/resolv.conf"); err != nil {
		return
	}
	a, err := u.LookupHost("☁→❄→☃→☀→☺→☂→☹→✝.ws.")
	if err != nil {
		t.Logf("Failed to lookup host %s\n", err.Error())
		t.Fail()
	}
	for _, r := range a {
		if len(r) == 0 {
			t.Log("Failure to get the A for ☁→❄→☃→☀→☺→☂→☹→✝.ws.")
			t.Fail()
			continue
		}
		t.Logf("Found %s\n", r)
	}
}
