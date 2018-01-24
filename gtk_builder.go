package goutils

import (
	"reflect"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

type BuilderBase struct {
	Builder *gtk.Builder
}

func NewBuilderBase(assetProvider AssetProvider, assetPath string) (*BuilderBase, error) {
	b := BuilderBase{}
	var err error
	b.Builder, err = GetBuilder(assetProvider, assetPath)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating Builder")
	}
	return &b, nil
}

func (b *BuilderBase) GetDialog(name string) (*gtk.Dialog, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Dialog)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Dialog.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetWindow(name string) (*gtk.Window, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Window)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Window.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetEntry(name string) (*gtk.Entry, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Entry)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Entry.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetLabel(name string) (*gtk.Label, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Label)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Label.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetComboBox(name string) (*gtk.ComboBox, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.ComboBox)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.ComboBox.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetTextView(name string) (*gtk.TextView, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.TextView)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.TextView.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetTextBuffer(name string) (*gtk.TextBuffer, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.TextBuffer)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.TextBuffer.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetStatusbar(name string) (*gtk.Statusbar, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Statusbar)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Statusbar.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetPaned(name string) (*gtk.Paned, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Paned)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Paned.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetTreeStore(name string) (*gtk.TreeStore, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.TreeStore)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.TreeStore.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetTreeView(name string) (*gtk.TreeView, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.TreeView)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.TreeView.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetNotebook(name string) (*gtk.Notebook, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Notebook)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Notebook.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetButton(name string) (*gtk.Button, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Button)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Button.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetListStore(name string) (*gtk.ListStore, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.ListStore)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.ListStore.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetInfoBar(name string) (*gtk.InfoBar, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.InfoBar)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Switch.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetSwitch(name string) (*gtk.Switch, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.Switch)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.Switch.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetSpinButton(name string) (*gtk.SpinButton, error) {
	obj, err := b.Builder.GetObject(name)
	if err != nil {
		return nil, errors.Errorf("Error getting %s", name)
	}
	widget, ok := obj.(*gtk.SpinButton)
	if !ok {
		return nil, errors.Errorf("Can't cast %s to gtk.SpinButton.", reflect.TypeOf(obj).String())
	}
	return widget, nil
}

func (b *BuilderBase) GetTextViewText(tv *gtk.TextView) (string, error) {
	buffer, err := tv.GetBuffer()
	if err != nil {
		return "", err
	}
	startIter := buffer.GetStartIter()
	endIter := buffer.GetEndIter()
	text, err := buffer.GetText(startIter, endIter, true)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (b *BuilderBase) SetTextViewText(tv *gtk.TextView, text string) error {
	buffer, err := tv.GetBuffer()
	if err != nil {
		return err
	}
	buffer.SetText(text)
	return nil
}
