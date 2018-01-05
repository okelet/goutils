package goutils

import "fmt"
import "github.com/gosexy/gettext"

func Gettext(message string, params ...interface{}) string {
	return fmt.Sprintf(gettext.Gettext(message), params)
}
