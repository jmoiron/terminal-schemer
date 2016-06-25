package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"sort"
	"strings"
)

// a profile as it exists in dconf
type profile struct {
	Key     string
	FullKey string
	Name    string
}

func newProfile(key string) profile {
	p := profile{
		Key:     key,
		FullKey: fmt.Sprintf("/org/mate/terminal/profiles/%s", key),
	}
	p.Name = p.get("visible-name")
	return p
}

// keys returns all keys for the given profile
func (p profile) keys() []string {
	out, err := dconf("list", p.FullKey+"/")
	if err != nil {
		fmt.Printf("Error reading keys for profile %s: %s", p.Key, err)
		return nil
	}
	sout := strings.Trim(string(out), "\n")
	ks := strings.Split(sout, "\n")
	sort.Strings(ks)
	return ks
}

// get a value for a given key from this profile
func (p profile) get(key string) string {
	path := path.Join(p.FullKey, key)
	out, err := dconf("read", path)
	if err != nil {
		fmt.Printf("Error reading %s: %s", path, err)
		return ""
	}
	return strings.Trim(string(out), "'\n")
}

// dconf runs a dconf command, returning its combined output and any error
func dconf(cmd ...string) ([]byte, error) {
	c := exec.Command("dconf", cmd...)
	return c.CombinedOutput()
}

// listProfiles returns a list of installed profiles
func listProfiles() ([]profile, error) {
	// XXX: only supporting Mate cus that's what I got, but it wouldn't be
	// too bad an idea to generalize it to gnome;  this tool is sorely missing.
	out, err := dconf("read", "/org/mate/terminal/global/profile-list")
	if err != nil {
		return nil, err
	}
	var profiles []profile
	var names []string
	// output is `['profile1', 'profile2', ..]` which is a valid.. python list?
	// turn it into a json list by replacing ' -> ".
	// this makes a copy but there probably aren't a lot of these.
	out = bytes.Replace(out, []byte(`'`), []byte(`"`), -1)
	err = json.Unmarshal(out, &names)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		profiles = append(profiles, newProfile(name))
	}
	return profiles, nil
}
