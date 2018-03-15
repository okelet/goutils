package goutils

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/pkg/errors"
)

var ProxychainsNotFoundError error

func init() {
	ProxychainsNotFoundError = errors.New("proxychains4 not found")
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
	if parsedUrl.User != nil {
		p.Username = parsedUrl.User.Username()
		password, hasPassword := parsedUrl.User.Password()
		if hasPassword {
			err := p.SetPassword(password)
			if err != nil {
				return nil, errors.Wrap(err, "Error setting the password for the proxy")
			}
		}
	}

	p.Exceptions = []string{}
	if proxyExceptions != "" {
		p.Exceptions = FilterEmptyStrings(strings.Split(proxyExceptions, ","))
	}

	return p, nil

}

func SetEnvironmentProxy(p *Proxy) error {
	var err error
	if p != nil {
		var url string
		url, err = p.ToUrl(true)
		if err != nil {
			return errors.Wrap(err, "Error generating proxy URL")
		}
		err = os.Setenv("http_proxy", url)
		if err != nil {
			return errors.Wrapf(err, "Error unsetting environment variable %v", "http_proxy")
		}
		err = os.Setenv("HTTP_PROXY", url)
		if err != nil {
			return errors.Wrapf(err, "Error unsetting environment variable %v", "HTTP_PROXY")
		}
		err = os.Setenv("https_proxy", url)
		if err != nil {
			return errors.Wrapf(err, "Error unsetting environment variable %v", "https_proxy")
		}
		err = os.Setenv("HTTPS_PROXY", url)
		if err != nil {
			return errors.Wrapf(err, "Error unsetting environment variable %v", "HTTPS_PROXY")
		}
		err = os.Setenv("ftp_proxy", url)
		if err != nil {
			return errors.Wrapf(err, "Error unsetting environment variable %v", "ftp_proxy")
		}
		err = os.Setenv("FTP_PROXY", url)
		if err != nil {
			return errors.Wrapf(err, "Error unsetting environment variable %v", "FTP_PROXY")
		}
		if len(p.Exceptions) > 0 {
			err = os.Setenv("no_proxy", strings.Join(p.Exceptions, ","))
			if err != nil {
				return errors.Wrapf(err, "Error setting environment variable %v", "no_proxy")
			}
			err = os.Setenv("NO_PROXY", strings.Join(p.Exceptions, ","))
			if err != nil {
				return errors.Wrapf(err, "Error setting environment variable %v", "NO_PROXY")
			}
		} else {
			err = os.Unsetenv("no_proxy")
			if err != nil {
				return errors.Wrapf(err, "Error unsetting environment variable %v", "no_proxy")
			}
			err = os.Unsetenv("NO_PROXY")
			if err != nil {
				return errors.Wrapf(err, "Error unsetting environment variable %v", "NO_PROXY")
			}
		}
	} else {
		for _, k := range []string{"http_proxy", "https_proxy", "ftp_proxy", "no_proxy"} {
			err = os.Unsetenv(k)
			if err != nil {
				return errors.Wrapf(err, "Error unsetting environment variable %v", k)
			}
			err = os.Unsetenv(strings.ToUpper(k))
			if err != nil {
				errors.Wrapf(err, "Error unsetting environment variable %v", strings.ToUpper(k))
			}
		}
	}
	return nil
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
