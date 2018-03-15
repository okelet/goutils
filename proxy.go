package goutils

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

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
