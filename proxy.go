package goutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gobwas/glob"
	"github.com/gotk3/gotk3/glib"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var PROXYCHAINS_PATH string

type ProxyPasswordManager interface {
	GetProxyPassword(p *Proxy) (string, error)
	SetProxyPassword(p *Proxy, password string) error
}

type SimpleProxyPasswordManager struct {
	Password string
}

func NewSimpleProxyPasswordManager(password string) *SimpleProxyPasswordManager {
	s := SimpleProxyPasswordManager{}
	s.Password = password
	return &s
}

func (s *SimpleProxyPasswordManager) GetProxyPassword(p *Proxy) (string, error) {
	return s.Password, nil
}

func (s *SimpleProxyPasswordManager) SetProxyPassword(p *Proxy, password string) error {
	s.Password = password
	return nil
}

func init() {
	var err error
	PROXYCHAINS_PATH, err = Which("proxychains4")
	if err != nil {
		panic(fmt.Sprintf("Error detecting proxychains4: %v", err.Error()))
	}
}

type Proxy struct {
	ProxyPasswordManager
	UUID       string
	Protocol   string
	Address    string
	Port       int
	Username   string
	Exceptions []string
}

func NewEmptyProxy(passwordManager ProxyPasswordManager) *Proxy {
	if passwordManager == nil {
		passwordManager = NewSimpleProxyPasswordManager("")
	}
	p := Proxy{ProxyPasswordManager: passwordManager}
	p.UUID = uuid.Must(uuid.NewV4()).String()
	p.Exceptions = []string{}
	return &p
}

func NewProxyFromMap(h *MapHelper, passwordManager ProxyPasswordManager, loadPasswordsFromMap bool) *Proxy {
	if passwordManager == nil {
		passwordManager = NewSimpleProxyPasswordManager("")
	}
	p := Proxy{ProxyPasswordManager: passwordManager}
	p.UUID = h.GetString("uuid", uuid.Must(uuid.NewV4()).String())
	p.Protocol = h.GetString("protocol", "http")
	p.Address = h.GetString("address", "127.0.0.1")
	p.Port = h.GetInt("port", 8080)
	p.Username = h.GetString("username", "")
	p.Exceptions = h.GetListOfStrings("exceptions", []string{})
	if loadPasswordsFromMap {
		p.SetProxyPassword(&p, h.GetString("password", ""))
	}
	return &p
}

func (p *Proxy) GetPassword() (string, error) {
	return p.ProxyPasswordManager.GetProxyPassword(p)
}

func (p *Proxy) SetPassword(password string) error {
	return p.ProxyPasswordManager.SetProxyPassword(p, password)
}

func (p *Proxy) Debug() {
	password, err := p.GetProxyPassword(p)
	if err != nil {
		password = fmt.Sprintf("ERROR: %v", err)
	}
	fmt.Printf("Protocol: %v\n", p.Protocol)
	fmt.Printf("Address: %v\n", p.Address)
	fmt.Printf("Port: %v\n", p.Port)
	fmt.Printf("Username: %v\n", p.Username)
	fmt.Printf("Password: %v\n", password)
	fmt.Printf("Exceptions: %v\n", p.Exceptions)
}

func (p *Proxy) ToUrl(includePassword bool) (string, error) {
	password, err := p.GetProxyPassword(p)
	if err != nil {
		return "", errors.Wrap(err, "Error getting the proxy password")
	}
	userPass := bytes.NewBufferString("")
	if p.Username != "" {
		userPass.WriteString(p.Username)
		if password != "" {
			userPass.WriteString(":")
			if includePassword {
				userPass.WriteString(password)
			} else {
				userPass.WriteString("*******")
			}
		}
		userPass.WriteString("@")
	}
	return fmt.Sprintf("%v://%v%v:%v", p.Protocol, userPass.String(), p.Address, p.Port), nil
}

func (p *Proxy) ToSimpleUrl() string {
	return fmt.Sprintf("%v://%v:%v", p.Protocol, p.Address, p.Port)
}

func (p *Proxy) ToMap(includePassword bool) (*MapHelper, error) {
	password, err := p.GetProxyPassword(p)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting the proxy password")
	}
	h := NewEmptyMapHelper()
	h.SetString("protocol", p.Protocol)
	h.SetString("address", p.Address)
	h.SetInt("port", p.Port)
	if p.Username != "" {
		h.SetString("username", p.Username)
	}
	if len(p.Exceptions) > 0 {
		h.SetListOfStrings("exceptions", p.Exceptions)
	}
	if includePassword && password != "" {
		h.SetString("password", password)
	}
	return h, nil
}

func (p *Proxy) IsValidForUrl(address string) (bool, error) {
	parsedUrl, err := url.Parse(address)
	if err != nil {
		return false, errors.Wrapf(err, "Error parsing address %v", address)
	}
	return p.IsValidForAddress(parsedUrl.Hostname())
}

func (p *Proxy) IsValidForAddress(host string) (bool, error) {
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

func GetEnvironmentProxy(passwordManager ProxyPasswordManager) (*Proxy, error) {

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

	p := NewEmptyProxy(passwordManager)
	p.Protocol = "http"
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
		err := passwordManager.SetProxyPassword(p, password)
		if err != nil {
			return nil, errors.Wrap(err, "Error setting the password for the proxy")
		}
	}

	p.Exceptions = []string{}
	if proxyExceptions != "" {
		p.Exceptions = strings.Split(proxyExceptions, ",")
	}

	return p, nil

}

func GetGnomeProxy(passwordManager ProxyPasswordManager) (*Proxy, error) {

	proxySettings := glib.SettingsNew("org.gnome.system.proxy")
	proxyMode := proxySettings.GetString("mode")
	exceptions := proxySettings.GetStrv("ignore-hosts")
	if proxyMode == "none" {
		return nil, nil
	} else if proxyMode == "manual" {

		p := NewEmptyProxy(passwordManager)
		p.Protocol = "http"
		httpProxySettings := glib.SettingsNew("org.gnome.system.proxy.http")
		httpsProxySettings := glib.SettingsNew("org.gnome.system.proxy.https")
		ftpProxySettings := glib.SettingsNew("org.gnome.system.proxy.ftp")

		p.Address = httpProxySettings.GetString("host")
		if p.Address == "" {
			p.Address = httpsProxySettings.GetString("host")
		}
		if p.Address == "" {
			p.Address = ftpProxySettings.GetString("host")
		}
		if p.Address == "" {
			return nil, errors.New("No host set for http, https and ftp proxy")
		}

		p.Port = httpProxySettings.GetInt("port")
		if p.Port <= 0 {
			p.Port = httpsProxySettings.GetInt("port")
		}
		if p.Port <= 0 {
			p.Port = ftpProxySettings.GetInt("port")
		}
		if p.Port <= 0 {
			return nil, errors.New("No port set for http, https and ftp proxy")
		}

		p.Exceptions = exceptions
		return p, nil

	} else {
		return nil, errors.Errorf("Unsupported proxy mode %v", proxyMode)
	}
}

func SetGnomeProxy(p *Proxy) error {

	proxySettings := glib.SettingsNew("org.gnome.system.proxy")
	if p != nil {
		proxySettings.SetString("mode", "manual")
		proxySettings.SetStrv("ignore-hosts", p.Exceptions)
		ftpProxySettings := glib.SettingsNew("org.gnome.system.proxy.ftp")
		httpProxySettings := glib.SettingsNew("org.gnome.system.proxy.http")
		httpsProxySettings := glib.SettingsNew("org.gnome.system.proxy.https")
		for _, s := range []*glib.Settings{ftpProxySettings, httpProxySettings, httpsProxySettings} {
			s.SetString("host", p.Address)
			s.SetInt("port", p.Port)
		}
	} else {
		proxySettings.SetString("mode", "none")
	}

	return nil

}

// TODO: Implement
func RunProxifiedBackgroundCommand(p *Proxy, command string, arguments []string) (int, *exec.Cmd, error) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		return 0, nil, errors.Wrapf(err, "Error creating temporary configuration file")
	}
	defer os.Remove(tempFile.Name())
	return 0, nil, nil
}
