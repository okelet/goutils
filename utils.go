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

func RunCommandAndWait(initialPath string, stdin io.Reader, command string, arguments []string) (error, int, int, string, string) {

	cmd := exec.Command(command, arguments...)

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

// Returns the full path of the specified command, if it is in the path.
// If the command is not found, no error is returned, but an empty string as path.
func Which(command string) (string, error) {
	err, _, exitCode, stdOut, stdErr := RunCommandAndWait("", nil, "bash", []string{"-c", "which " + command})
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
