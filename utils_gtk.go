package goutils

import (
	"os/exec"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

func XdgOpenFromMenuItem(item *gtk.MenuItem, url string) {
	exec.Command("xdg-open", url).Start()
}

type AssetProvider interface {
	GetAsset(assetName string) ([]byte, error)
}

func LoadPixbufFromAsset(assetProvider AssetProvider, assetPath string) (*gdk.Pixbuf, error) {

	data, err := assetProvider.GetAsset(assetPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error loading icon from %v", assetPath)
	}

	pixBufLoader, err := gdk.PixbufLoaderNew()
	if err != nil {
		return nil, errors.Wrap(err, "Error creating PixBufLoader")
	}

	// Set size
	pixBufLoader.SetSize(16, 16)

	_, err = pixBufLoader.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating PixBufLoader from data")
	}

	err = pixBufLoader.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Error closing PixBufLoader")
	}

	pixbuf, err := pixBufLoader.GetPixbuf()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting PixBufLoader")
	}

	return pixbuf, nil

}

func ShowMessageWrapped(w *gtk.Window, messageType gtk.MessageType, title string, message string) {
	glib.IdleAdd(func() {
		ShowMessage(w, messageType, title, message)
	})
}

func ShowMessage(w *gtk.Window, messageType gtk.MessageType, title string, message string) {
	d := gtk.MessageDialogNew(w, gtk.DIALOG_MODAL, messageType, gtk.BUTTONS_OK, title)
	d.FormatSecondaryMarkup(message)
	d.Run()
	d.Destroy()
}

func ConfirmMessage(w *gtk.Window, title string, message string) bool {
	d := gtk.MessageDialogNew(w, gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, title)
	d.FormatSecondaryMarkup(message)
	res := d.Run()
	d.Destroy()
	if res == int(gtk.RESPONSE_YES) {
		return true
	} else {
		return false
	}
}

func GetBuilder(assetProvider AssetProvider, assetPath string) (*gtk.Builder, error) {

	data, err := assetProvider.GetAsset(assetPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error loading asset "+assetPath)
	}
	builder, err := gtk.BuilderNew()
	if err != nil {
		return nil, errors.Wrap(err, "Error creating gtk.BuilderNew")
	}
	builder.AddFromString(string(data[:]))
	return builder, nil

}
