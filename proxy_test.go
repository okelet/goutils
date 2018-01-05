package goutils

import (
	"testing"
	"fmt"
)

func TestProxy(t *testing.T) {

	p, err := LoadProxyFromGnome()
	if err != nil {
		fmt.Printf("Error loading proxy from Gnome: %v", err)
	} else if p != nil {
		if p != nil {
			p.Debug()
		}
	}

	p, err = LoadProxyFromEnvironment()
	if err != nil {
		fmt.Printf("Error loading proxy from environment: %v", err)
	} else if p != nil {
		p.Debug()
	}

	p = &Proxy{
		Protocol:   "http",
		Address:    "127.0.0.2",
		Port:       8080,
		Username:   "user",
		Password:   "pass",
		Exceptions: []string{"*.test.com"},
	}
	p.Debug()
	for _, v := range []string{"sub.test.com", "test.com", "hello.net"} {
		matches, err := p.MatchesAddress(v)
		if err != nil {
			fmt.Printf("Error checking if %v matches: %v\n", v, err)
		} else {
			fmt.Printf("Matches %v: %v\n", v, matches)
		}
	}

}
