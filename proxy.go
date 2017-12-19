package goutils

import (
	"github.com/gotk3/gotk3/glib"
	"os"
	"net/url"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"github.com/gobwas/glob"
	"fmt"
)

type Proxy struct{
	Protocol string
	Address string
	Port int
	Username string
	Password string
	Exceptions []string
}

func (p *Proxy) Debug() {
	fmt.Printf("Protocol: %v\n", p.Protocol)
	fmt.Printf("Address: %v\n", p.Address)
	fmt.Printf("Port: %v\n", p.Port)
	fmt.Printf("Username: %v\n", p.Username)
	fmt.Printf("Password: %v\n", p.Password)
	fmt.Printf("Exceptions: %v\n", p.Exceptions)
}

func (p *Proxy) MatchesUrl(address string) (bool, error) {
	parsedUrl, err := url.Parse(address)
	if err != nil {
		return false, errors.Wrapf(err, "Error parsing address %v", address)
	}
	return p.MatchesAddress(parsedUrl.Hostname())
}

func (p *Proxy) MatchesAddress(host string) (bool, error) {
	for _, v := range p.Exceptions {
		g, err := glob.Compile(v)
		if err != nil {
			return false, errors.Wrapf(err, "Error converting %v to glob", v)
		}
		if g.Match(host) {
			return true, nil
		}
	}
	return false, nil
}


func LoadProxyFromEnvironment() (*Proxy, error) {

	httpProxy := os.Getenv("http_proxy")
	if httpProxy == "" {
		httpProxy = os.Getenv("HTTP_PROXY")
	}
	if httpProxy == "" {
		httpProxy = os.Getenv("https_proxy")
	}
	if httpProxy == "" {
		httpProxy = os.Getenv("HTTPS_PROXY")
	}
	if httpProxy == "" {
		httpProxy = os.Getenv("ftp_proxy")
	}
	if httpProxy == "" {
		httpProxy = os.Getenv("FTP_PROXY")
	}
	if httpProxy == "" {
		return nil, nil
	}

	proxyExceptions := os.Getenv("no_proxy")
	if proxyExceptions == "" {
		proxyExceptions = os.Getenv("NO_PROXY")
	}

	parsedUrl, err := url.Parse(httpProxy)
	if err != nil {
		return nil, errors.Wrapf(err, "Error parsing proxy URL %v", httpProxy)
	}

	p := Proxy{}
	p.Protocol = parsedUrl.Scheme
	p.Address = parsedUrl.Hostname()
	port := parsedUrl.Port()
	if port == "" {
		if p.Protocol == "socks" || p.Protocol == "socks4" || p.Protocol == "socks5" {
			p.Port = 1080
		} else {
			p.Port = 8080
		}
	} else {
		i, err := strconv.Atoi(port)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing port %v as integer", port)
		}
		p.Port = i
	}
	p.Username = parsedUrl.User.Username()
	password, hasPassword := parsedUrl.User.Password()
	if hasPassword {
		p.Password = password
	}

	p.Exceptions = []string{}
	if proxyExceptions != "" {
		p.Exceptions = strings.Split(proxyExceptions, ",")
	}

	return &p, nil

}

func LoadProxyFromGnome() (*Proxy, error) {

	s := glib.SettingsNew("org.gnome.system.proxy")
	proxyMode := s.GetString("mode")
	if proxyMode != "manual" {
		return nil, nil
	}

	httpProxySettings := glib.SettingsNew("org.gnome.system.proxy.http")
	// httpsProxySettings := glib.SettingsNew("org.gnome.system.proxy.https")
	// ftpProxySettings := glib.SettingsNew("org.gnome.system.proxy.ftp")
	p := Proxy{}
	p.Address = httpProxySettings.GetString("host")
	p.Port = httpProxySettings.GetInt("port")
	return &p, nil

}
