package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type profileSchema struct {
	Name            string  `json:"name"`
	BackgroundColor color   `json:"background-color"`
	ForegroundColor color   `json:"foreground-color"`
	Palette         []color `json:"palette"`
	Font            string  `json:"font"`
	AllowBold       bool    `json:"allow-bold"`
}

func (p profile) export() string {
	s := profileSchema{
		Name:            p.Name,
		Font:            p.get("font"),
		BackgroundColor: color(p.get("background-color")).toHex24(),
		ForegroundColor: color(p.get("foreground-color")).toHex24(),
		AllowBold:       p.get("allow-bold") == "true",
	}
	// color conversions
	palette := p.get("palette")
	for _, c := range strings.Split(palette, ":") {
		s.Palette = append(s.Palette, color(c).toHex24())
	}

	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		return ""
	}
	return string(b)
}

// load creates the profile under the name `name` if given, otherwise it
// uses the Name from the profileSchema as the name of the profile.
func (p profileSchema) load(name string) error {
	return nil
}
