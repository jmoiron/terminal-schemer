package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jmoiron/gotk3/glib"
	"github.com/stretchr/testify/assert"
)

func TestGSettingsPOC(t *testing.T) {
	assert := assert.New(t)
	source := glib.SettingsSchemaSourceGetDefault()
	assert.NotNil(source)
	source.Ref()
	defer source.Unref()

	nonrel, rel := source.ListSchemas(false)
	assert.NotEmpty(nonrel)
	assert.NotEmpty(rel)

	schema := source.Lookup(nonrel[0], true)
	assert.NotNil(schema)

	schema = source.Lookup(rel[0], true)
	assert.NotNil(schema)

	var terminalSchemas []string
	for _, key := range rel {
		switch key {
		case "org.mate.terminal.profile":
			terminalSchemas = append(terminalSchemas, key)
		case "org.gnome.Terminal.Legacy.Profile":
			terminalSchemas = append(terminalSchemas, key)
		}
	}
	if len(terminalSchemas) == 0 {
		t.Skip("no terminal schemas")
	}
	for _, ts := range terminalSchemas {
		schema = source.Lookup(ts, true)
		assert.NotNil(schema)
		fmt.Println(ts, schema)
		fmt.Println(schema.ListChildren())
		fmt.Println(schema.ListKeys())
	}

	schema = glib.SettingsSchemaSourceGetDefault().Lookup(mateProfileList.schema, true)
	settings := glib.SettingsNew(mateProfileList.schema)
	fmt.Println(settings.ListChildren())
	fmt.Println(schema.ListKeys())

	val := settings.GetValue(mateProfileList.key)
	fmt.Println("variant val:", val.GetStrv())
	fmt.Println("variant container?:", val.IsContainer())
	fmt.Println(val)
	fmt.Println(val.AnnotatedString())
	fmt.Println(val.Type())

	settings = glib.SettingsNew("org.mate.interface")
	fmt.Println("cursor-blink-time", settings.GetInt("cursor-blink-time"))
	fmt.Println(settings.GetValue("cursor-blink-time").GetInt())

	source = glib.SettingsSchemaSourceGetDefault()
	profiles := val.GetStrv()
	for _, prof := range profiles {
		path := filepath.Join("/org/mate/terminal/profiles", prof) + "/"
		schema := "org.mate.terminal.profile"
		fmt.Println(path, schema)
		settings := glib.SettingsNewWithPath("org.mate.terminal.profile", path)
		keys := source.Lookup(schema, true).ListKeys()
		fmt.Println(schema, path, settings, keys)
		for _, key := range keys {
			variant := settings.GetValue(key)
			fmt.Println(key, variant)
		}
	}
}
