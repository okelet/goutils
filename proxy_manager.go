package goutils

import (
	"net/url"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"
)

const PROXY_METHOD_PAC = "pac"
const PROXY_METHOD_SIMPLE = "simple"
const PROXY_METHOD_DIRECT = "direct"

type ProxyManager struct {
	UUID        string
	Method      string
	PACURL      string
	SimpleProxy *Proxy
}

func (pm *ProxyManager) SetDirectMethod(pacUrl string) {
	pm.Method = PROXY_METHOD_DIRECT
}

func (pm *ProxyManager) SetPACMethod(pacUrl string) {
	pm.Method = PROXY_METHOD_PAC
	pm.PACURL = pacUrl
}

func (pm *ProxyManager) SetSimpleMethod(proxy *Proxy) {
	pm.Method = PROXY_METHOD_SIMPLE
	pm.SimpleProxy = proxy
}

func (pm *ProxyManager) GetProxyForUrl(destinationUrl string) (*Proxy, error) {
	if pm.Method == PROXY_METHOD_DIRECT {
		// Direct
		return nil, nil
	} else if pm.Method == PROXY_METHOD_PAC {
		// TODO
		return nil, errors.New("Not implemented")
	} else if pm.Method == PROXY_METHOD_PAC {
		// Simple mode
		parsedUrl, err := url.Parse(destinationUrl)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing address %v", destinationUrl)
		}
		return pm.GetProxyForAddress(parsedUrl.Hostname())
	} else {
		// Unknown
		return nil, errors.New("Unknown proxy method discovery")
	}
}

func (pm *ProxyManager) GetProxyForAddress(destinationAddress string) (*Proxy, error) {
	if pm.Method == PROXY_METHOD_DIRECT {
		// Direct
		return nil, nil
	} else if pm.Method == PROXY_METHOD_PAC {
		// TODO
		return nil, errors.New("Not implemented")
	} else if pm.Method == PROXY_METHOD_SIMPLE {
		// Simple mode
		for _, v := range pm.SimpleProxy.Exceptions {
			g, err := glob.Compile(v)
			if err != nil {
				return nil, errors.Wrapf(err, "Error converting %v to glob", v)
			}
			// If the exception matches the address parameter, the proxy is not valid
			if g.Match(destinationAddress) {
				return nil, nil
			}
		}
		return pm.SimpleProxy, nil
	} else {
		// Unknown
		return nil, errors.New("Unknown proxy method discovery")
	}
}

func (pm *ProxyManager) GetDefaultProxy() (*Proxy, error) {
	if pm.Method == PROXY_METHOD_DIRECT {
		// Direct
		return nil, nil
	} else if pm.Method == PROXY_METHOD_PAC {
		// TODO
		return nil, errors.New("Not implemented")
	} else if pm.Method == PROXY_METHOD_SIMPLE {
		return pm.SimpleProxy, nil
	} else {
		// Unknown
		return nil, errors.New("Unknown proxy method discovery")
	}
}

func (pm *ProxyManager) ToMap(includePassword bool) (*MapHelper, error) {

	if pm.Method == PROXY_METHOD_DIRECT {
		// Direct
		h := NewEmptyMapHelper()
		h.SetString("method", PROXY_METHOD_DIRECT)
		return h, nil
	} else if pm.Method == PROXY_METHOD_PAC {
		h := NewEmptyMapHelper()
		h.SetString("method", PROXY_METHOD_PAC)
		h.SetString("pac", pm.PACURL)
		return h, nil
	} else if pm.Method == PROXY_METHOD_SIMPLE {
		h := NewEmptyMapHelper()
		h.SetString("method", PROXY_METHOD_SIMPLE)
		proxyData, err := pm.SimpleProxy.ToMap(includePassword)
		if err != nil {
			return nil, errors.Wrap(err, "Error getting proxy data")
		}
		h.SetHelper("proxy", proxyData)
		return h, nil
	} else {
		// Unknown
		return nil, errors.New("Unknown proxy method discovery")
	}

}
