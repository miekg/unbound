package main

// https://www.unbound.net/documentation/libunbound-tutorial-1.html

import (
	"fmt"
	"github.com/miekg/unbound"
	"log"
)

func main() {
	u := unbound.New()
	defer u.Destroy()

	addr, err := u.LookupHost("wwww.nlnetlabs.nl")
	if err != nil {
		log.Fatalf("error %s\n", err.Error())
	}

	for _, a := range addr {
		fmt.Printf("The address is %s\n", a)
	}
}
