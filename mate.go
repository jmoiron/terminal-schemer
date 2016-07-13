package main

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/glib"
)

// mate.go contains a gsettings terminal profile schema editor implementation
// for MATE, which these days is a bit of a hybrid of the final versions of
// Gnome2 and Gnome3 when it comes to config.

type setting struct {
	schema string
	key    string
}

var (
	mateProfileList = setting{"org.mate.terminal.global", "profile-list"}
)

const (
	mateProfileSchema      = "org.mate.terminal.profile"
	mateProfilePath        = "/org/mate/terminal/profiles/%s/"
	mateNameKey            = "visible-name"
	mateBackgroundColorKey = "background-color"
	mateForegroundColorKey = "foreground-color"
	matePaletteKey         = "palette"
	mateFontKey            = "font"
	mateAllowBoldKey       = "allow-bold"
)

type mateProfile struct {
	key  string
	path string
	// parts of the profile considered for export and import
	name            string
	backgroundColor string
	foregroundCOlor string
	palette         []string
	font            string
	allowBold       bool
}

func newMateProfile(key string) (*mateProfile, error) {
	p := &mateProfile{
		key:  key,
		path: fmt.Sprintf(mateProfilePath, key),
	}
	return p, p.load()
}

// load reads the profile from gsettings
func (p *mateProfile) load() error {
	source := glib.SettingsSchemaSourceGetDefault()
	settings := glib.SettingsNewWithPath(mateProfileSchema, p.path)
	keys := source.Lookup(mateProfileSchema, true).ListKeys()
	keySet := map[string]struct{}{}
	for _, key := range keys {
		keySet[key] = struct{}{}
	}
	fmt.Println("Lengths", len(keys), len(keySet))
	if _, ok := keySet[mateNameKey]; ok {
		p.name = settings.GetString(mateNameKey)
	}
	if _, ok := keySet[mateBackgroundColorKey]; ok {
		p.backgroundColor = settings.GetString(mateBackgroundColorKey)
	}
	if _, ok := keySet[mateForegroundColorKey]; ok {
		p.foregroundCOlor = settings.GetString(mateForegroundColorKey)
	}
	if _, ok := keySet[matePaletteKey]; ok {
		// XXX is this ever a strv?
		paletteString := settings.GetString(matePaletteKey)
		p.palette = strings.Split(paletteString, ":")
	}
	if _, ok := keySet[mateFontKey]; ok {
		p.font = settings.GetString(mateFontKey)
	}
	if _, ok := keySet[mateAllowBoldKey]; ok {
		p.allowBold = settings.GetBoolean(mateAllowBoldKey)
	}
	return nil
}

func listMateProfileNames() []string {
	return glib.SettingsNew(mateProfileList.schema).GetValue(mateProfileList.key).GetStrv()
}
