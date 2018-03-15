package goutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

func ShowZenityError(title string, message string) {
	err := exec.Command("zenity", "--error", "--title", title, "--text", message).Run()
	if err != nil {
		fmt.Printf("Error showing error message: %v", err)
	}
}

func RunCommandBackground(command string, arguments []string) (int, *exec.Cmd, error) {
	cmd := exec.Command(command, arguments...)
	err := cmd.Start()
	if err != nil {
		return 0, nil, errors.Wrapf(err, "Error running command")
	}
	return cmd.Process.Pid, cmd, nil
}

func RunCommandAndWait(initialPath string, stdin io.Reader, command string, arguments []string, env map[string]string) (error, int, int, string, string) {

	cmd := exec.Command(command, arguments...)

	if env != nil {
		envList := []string{}
		for key, val := range env {
			envList = append(envList, fmt.Sprintf("%v=%v", key, val))
		}
		cmd.Env = envList
	}

	if stdin != nil {
		cmd.Stdin = stdin
	}

	var outBuff bytes.Buffer
	cmd.Stdout = &outBuff

	var errBuff bytes.Buffer
	cmd.Stderr = &errBuff

	if initialPath != "" {
		cmd.Dir = initialPath
	}

	err := cmd.Start()
	if err != nil {
		return err, 0, 0, "", ""
	}

	state, err := cmd.Process.Wait()
	if err != nil {
		return err, 0, 0, "", ""
	}

	exitCode := state.Sys().(syscall.WaitStatus).ExitStatus()
	pid := cmd.Process.Pid
	return nil, pid, exitCode, outBuff.String(), errBuff.String()

}

func RunCommandProxified(pm *ProxyManager, destinationUrl string, destinationAddress string, initialPath string, stdin io.Reader, command string, arguments []string, env map[string]string, callback func(error, int, int, string, string)) (error, int, int, string, string) {

	var err error

	if env == nil {
		env = map[string]string{}
		for _, e := range os.Environ() {
			elems := strings.SplitN(e, "=", 2)
			val := ""
			if len(elems) >= 2 {
				val = elems[1]
			}
			env[elems[0]] = val
		}
	}

	var cmd *exec.Cmd
	var p *Proxy
	if pm != nil {

		if destinationUrl != "" {

			p, err = pm.GetProxyForUrl(destinationUrl)
			if err != nil {
				newErr := errors.Wrap(err, "Error checking if proxy is valid for URL")
				if callback != nil {
					callback(newErr, 0, 0, "", "")
				}
				return newErr, 0, 0, "", ""
			}

		} else if destinationAddress != "" {

			p, err = pm.GetProxyForAddress(destinationAddress)
			if err != nil {
				newErr := errors.Wrap(err, "Error checking if proxy is valid for address")
				if callback != nil {
					callback(newErr, 0, 0, "", "")
				}
				return newErr, 0, 0, "", ""
			}

		} else {
			p, err = pm.GetDefaultProxy()
			if err != nil {
				newErr := errors.Wrap(err, "Error getting default proxy")
				if callback != nil {
					callback(newErr, 0, 0, "", "")
				}
				return newErr, 0, 0, "", ""
			}
		}

	}

	if p != nil {

		proxychainsPath, err := Which("proxychains4")
		if err != nil {
			newErr := errors.Wrap(err, "Error checking if proxychains is installed")
			if callback != nil {
				callback(newErr, 0, 0, "", "")
			}
			return newErr, 0, 0, "", ""
		}

		if proxychainsPath == "" {
			if callback != nil {
				callback(ProxychainsNotFoundError, 0, 0, "", "")
			}
			return ProxychainsNotFoundError, 0, 0, "", ""
		}

		password, err := p.GetPassword()
		if err != nil {
			newErr := errors.Wrap(err, "Error getting proxy password")
			if callback != nil {
				callback(newErr, 0, 0, "", "")
			}
			return newErr, 0, 0, "", ""
		}

		proxychainsConfigFileContents := bytes.Buffer{}
		proxychainsConfigFileContents.WriteString("strict_chain\n")
		proxychainsConfigFileContents.WriteString("proxy_dns\n")
		proxychainsConfigFileContents.WriteString("[ProxyList]\n")
		if p.Username != "" && password != "" {
			proxychainsConfigFileContents.WriteString(fmt.Sprintf("%v %v %v %v %v\n", p.Protocol, p.Address, p.Port, p.Username, password))
		} else {
			proxychainsConfigFileContents.WriteString(fmt.Sprintf("%v %v %v\n", p.Protocol, p.Address, p.Port))
		}

		proxychainsConfigFile, err := ioutil.TempFile("", "")
		if err != nil {
			newErr := errors.Wrap(err, "Error generating temporary proxychains config file")
			if callback != nil {
				callback(newErr, 0, 0, "", "")
			}
			return newErr, 0, 0, "", ""
		}

		_, err = proxychainsConfigFileContents.WriteTo(proxychainsConfigFile)
		if err != nil {
			newErr := errors.Wrap(err, "Error writing temporary proxychains config file")
			if callback != nil {
				callback(newErr, 0, 0, "", "")
			}
			return newErr, 0, 0, "", ""
		}
		proxychainsConfigFile.Close()
		Log.Debugf("Proxychains config file generated in %v", proxychainsConfigFile.Name())
		defer os.Remove(proxychainsConfigFile.Name())

		env["PROXYCHAINS_CONF_FILE"] = proxychainsConfigFile.Name()
		env["PROXYCHAINS_QUIET_MODE"] = "1"

		newCommand := proxychainsPath
		newArguments := append([]string{command}, arguments...)
		cmd = exec.Command(newCommand, newArguments...)

	} else {
		cmd = exec.Command(command, arguments...)
	}

	envList := []string{}
	for key, val := range env {
		envList = append(envList, fmt.Sprintf("%v=%v", key, val))
	}
	cmd.Env = envList

	if stdin != nil {
		cmd.Stdin = stdin
	}

	var outBuff bytes.Buffer
	cmd.Stdout = &outBuff

	var errBuff bytes.Buffer
	cmd.Stderr = &errBuff

	if initialPath != "" {
		cmd.Dir = initialPath
	}

	err = cmd.Start()
	if err != nil {
		if callback != nil {
			callback(err, 0, 0, "", "")
		}
		return err, 0, 0, "", ""
	}

	state, err := cmd.Process.Wait()
	if err != nil {
		if callback != nil {
			callback(err, 0, 0, "", "")
		}
		return err, 0, 0, "", ""
	}

	exitCode := state.Sys().(syscall.WaitStatus).ExitStatus()
	pid := cmd.Process.Pid

	if callback != nil {
		callback(nil, pid, exitCode, outBuff.String(), errBuff.String())
	}

	return nil, pid, exitCode, outBuff.String(), errBuff.String()

}

func CombineStdErrOutput(stdOut string, errOut string) string {
	stdOut = strings.TrimSuffix(stdOut, "\n")
	errOut = strings.TrimSuffix(errOut, "\n")
	if stdOut != "" && errOut != "" {
		return fmt.Sprintf("%v\n%v", stdOut, errOut)
	} else if stdOut != "" {
		return stdOut
	} else if errOut != "" {
		return errOut
	} else {
		return ""
	}
}

// Returns the full path of the specified command, if it is in the path.
// If the command is not found, no error is returned, but an empty string as path.
func Which(command string) (string, error) {
	err, _, exitCode, stdOut, stdErr := RunCommandAndWait("", nil, "bash", []string{"-c", "which " + command}, map[string]string{})
	if err != nil {
		return "", errors.Wrapf(err, "Error locating command %s", command)
	} else if exitCode == 1 {
		// Not found
		return "", nil
	} else if exitCode != 0 {
		return "", errors.Errorf("Error detecting command %s: %s/%s", command, stdOut, stdErr)
	}
	parts := FilterEmptyStrings(strings.Split(stdOut, "\n"))
	if len(parts) == 0 {
		return "", nil
	} else if len(parts) > 1 {
		return "", errors.Errorf("Command %s not found (%d lines)", command, len(parts))
	} else {
		return parts[0], nil
	}
}

func HomeDir() (string, error) {
	dir := os.Getenv("HOME")
	if dir == "" {
		return "", errors.New("Empty HOME environment variable")
	}
	return dir, nil
}

func ExpandHomeDir(path string) (string, error) {

	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	dir, err := HomeDir()
	if err != nil {
		return "", errors.Wrap(err, "Empty HOME environment variable")
	}

	re, err := regexp.Compile("^~")
	if err != nil {
		return "", errors.Wrapf(err, "Error creatng regular expression to expand home dir")
	}

	return re.ReplaceAllString(path, dir), nil

}

func ContractHomeDir(path string) (string, error) {

	if len(path) == 0 {
		return path, nil
	}

	dir, err := HomeDir()
	if err != nil {
		return "", errors.Wrap(err, "Empty HOME environment variable")
	}

	re, err := regexp.Compile("^" + regexp.QuoteMeta(dir) + "(/|$)")
	if err != nil {
		return "", err
	}

	if re.MatchString(path) {
		return re.ReplaceAllString(path, "~${1}"), nil
	}

	return path, nil

}

func EnsureDirectoryExists(path string, mode os.FileMode) error {
	pathStat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, mode)
			if err != nil {
				return errors.Errorf("Error creating application directory %v: %v", path, err)
			}
		} else {
			return errors.Errorf("Error checking if application directory %v exists: %v", path, err)
		}
	} else if !pathStat.IsDir() {
		return errors.Errorf("Application directory %v exists, but is not a directory", path)
	}
	return nil
}

func CreateUserXdgShortcut(applicationId string, appPath string, appName string, appComment string, icon string, categories []string) error {

	// TODO Add keywords
	content := bytes.NewBufferString("[Desktop Entry]\n")
	content.WriteString("Type=Application\n")
	content.WriteString(fmt.Sprintf("Name=%v\n", appName))
	if appComment != "" {
		content.WriteString(fmt.Sprintf("Comment=%v\n", appComment))
	}
	content.WriteString(fmt.Sprintf("Exec=%v\n", appPath))
	content.WriteString(fmt.Sprintf("Icon=%v\n", icon))
	if len(categories) > 0 {
		content.WriteString(fmt.Sprintf("Categories=%v\n", strings.Join(categories, ";")+";"))
	}

	homeDir, err := ExpandHomeDir("~")
	if err != nil {
		return errors.Wrapf(err, "Error getting user home dir")
	}

	homeLocalDir := path.Join(homeDir, ".local")
	err = EnsureDirectoryExists(homeLocalDir, 0600)
	if err != nil {
		return errors.Wrapf(err, "Error checking/creating directory %v", homeLocalDir)
	}

	homeLocalShareDir := path.Join(homeLocalDir, "share")
	err = EnsureDirectoryExists(homeLocalShareDir, 0600)
	if err != nil {
		return errors.Wrapf(err, "Error checking/creating directory %v", homeLocalShareDir)
	}

	homeLocalShareApplicationsDir := path.Join(homeLocalShareDir, "applications")
	err = EnsureDirectoryExists(homeLocalShareApplicationsDir, 0600)
	if err != nil {
		return errors.Wrapf(err, "Error checking/creating directory %v", homeLocalShareApplicationsDir)
	}

	shortcutFile := path.Join(homeLocalShareApplicationsDir, applicationId+".desktop")
	err = ioutil.WriteFile(shortcutFile, content.Bytes(), 0644)
	if err != nil {
		return errors.Wrapf(err, "Error creating shortcut %v", shortcutFile)
	}

	return nil

}

func SetXdgAutostart(applicationId string, appPath string, appName string, icon string, enabled bool) error {

	content := "[Desktop Entry]\n"
	content += fmt.Sprintf("Type=Application\n")
	content += fmt.Sprintf("Name=%v\n", appName)
	content += fmt.Sprintf("Exec=%v\n", appPath)
	content += fmt.Sprintf("Icon=%v\n", icon)
	content += fmt.Sprintf("X-GNOME-Autostart-enabled=%v\n", enabled)

	homeDir, err := ExpandHomeDir("~")
	if err != nil {
		return errors.Wrapf(err, "Error getting user home dir")
	}

	homeConfigDir := path.Join(homeDir, ".config")
	err = EnsureDirectoryExists(homeConfigDir, 0700)
	if err != nil {
		return errors.Wrapf(err, "Error checking/creating directory %v", homeConfigDir)
	}

	homeConfigAutostartDir := path.Join(homeConfigDir, "autostart")
	err = EnsureDirectoryExists(homeConfigAutostartDir, 0700)
	if err != nil {
		return errors.Wrapf(err, "Error checking/creating directory %v", homeConfigAutostartDir)
	}

	shortcutFile := path.Join(homeConfigAutostartDir, applicationId+".desktop")
	err = ioutil.WriteFile(shortcutFile, []byte(content), 0644)
	if err != nil {
		return errors.Wrapf(err, "Error creating shortcut %v", shortcutFile)
	}

	return nil

}

func LoadJsonFileAsMap(path string, failIfNotFound bool) (map[string]interface{}, error) {

	exists, err := FileExists(path)
	if err != nil {
		return nil, err
	}

	if !exists {
		if failIfNotFound {
			return nil, errors.Errorf("File %v doesn't exist.", path)
		} else {
			return map[string]interface{}{}, nil
		}
	}

	reader, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Error opening file %v.", path)
	}

	decoder := json.NewDecoder(reader)

	var x map[string]interface{}
	err = decoder.Decode(&x)
	if err != nil {
		return nil, errors.Wrapf(err, "Error loading JSON from file %v.", path)
	}

	return x, nil

}

func SaveMapAsJsonFile(path string, data map[string]interface{}) error {
	byteData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "Error marshalling map")
	}
	// TODO: usar algo como os.ModeFile, igual que hay os.ModeDir
	err = ioutil.WriteFile(path, byteData, 0664)
	if err != nil {
		return errors.Wrapf(err, "Error saving JSON file")
	}
	return nil
}
