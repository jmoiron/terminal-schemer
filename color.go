package main

import (
	"fmt"
	"strconv"
	"strings"
)

type color string

func (c color) toHex24() color {
	s := strings.Trim(string(c), "#")
	switch len(s) {
	case 6:
		return color("#" + s)
	case 12:
		r, _ := strconv.ParseInt(s[:4], 16, 64)
		g, _ := strconv.ParseInt(s[4:8], 16, 64)
		b, _ := strconv.ParseInt(s[8:], 16, 64)
		return color(fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8))
	}
	return c
}

func (c color) toHex48() color {
	return c
}
