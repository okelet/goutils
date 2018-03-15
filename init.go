package goutils

import "github.com/juju/loggo"

var Log loggo.Logger

var PACKAGE_NAME string

func init() {

	PACKAGE_NAME = "com.github.okelet.goutils"

	Log = loggo.GetLogger(PACKAGE_NAME)

}
