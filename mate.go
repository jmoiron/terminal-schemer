package main

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

type MateTerminalSchema struct {
	profiles []string
}

func NewMateTerminalSchema() *MateTerminalSchema {

	return &MateTerminalSchema{}
}
