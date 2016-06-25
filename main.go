package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// mate-terminal-schemer is a program that can save and apply terminal
// colorschemes to new or existing mate profiles.

var opts struct {
	listProfiles bool
	show         string
	export       string
	import_      string
	backup       string
	restore      string
}

func main() {
	flag.BoolVar(&opts.listProfiles, "list", false, "list installed profiles")
	flag.StringVar(&opts.show, "show", "", "show an installed profile's attributes")
	flag.StringVar(&opts.export, "export", "", "export a named profile to stdout")
	flag.StringVar(&opts.import_, "import", "", "import a profile from a file")
	flag.StringVar(&opts.backup, "backup", "", "backup all current profiles to a directory")
	flag.StringVar(&opts.restore, "restore", "", "restore all profiles from a directory")
	flag.Parse()

	switch {
	case opts.listProfiles:
		listCmd()
	case len(opts.show) > 0:
		showCmd(opts.show)
	case len(opts.export) > 0:
		exportCmd(opts.export)
	case len(opts.import_) > 0:
		importCmd(opts.import_)
	case len(opts.backup) > 0:
		backupCmd(opts.backup)
	case len(opts.restore) > 0:
		restoreCmd(opts.restore)
	}
}

func listCmd() {
	profiles, err := listProfiles()
	if err != nil {
		fmt.Printf("Error listing profiles: %s\n", err)
		return
	}
	for _, p := range profiles {
		fmt.Printf("%s\n", p.Name)
	}
	return
}

func showCmd(name string) {
	p, err := findProfile(name)
	if err != nil {
		fmt.Printf("Error finding %s: %s\n", name, err)
		return
	}
	for _, key := range p.keys() {
		val := p.get(key)
		fmt.Printf("%s: %s\n", key, val)
	}
}

func exportCmd(name string) {
	p, err := findProfile(name)
	if err != nil {
		fmt.Printf("Error finding %s: %s\n", name, err)
		return
	}
	fmt.Println(p.export())
}

func importCmd(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	var p profileSchema
	err = json.NewDecoder(f).Decode(&p)
	if err != nil {
		fmt.Printf("Error in %s: %s\n", path, err)
	}
	fmt.Printf("Importing %#v\n", p)
}

func restoreCmd(dirpath string) {
	fis, err := ioutil.ReadDir(dirpath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	for _, fi := range fis {
		// skip if it's a directory
		if fi.IsDir() {
			continue
		}
		// skip if it's not a json file
		if !strings.HasSuffix(fi.Name(), ".json") {
			continue
		}
		fmt.Println(fi.Name())
	}
}

func backupCmd(dirpath string) {
	fi, err := os.Stat(dirpath)
	switch {
	case os.IsNotExist(err):
		// create the directory if it doesn't exist
		err = os.MkdirAll(dirpath, 0755)
	case err == nil:
		if !fi.IsDir() {
			fmt.Printf("Error: %s is not a directory\n", dirpath)
		}
	}
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

}

func findProfile(name string) (profile, error) {
	all, err := listProfiles()
	if err != nil {
		return profile{}, err
	}
	var p *profile
	for _, a := range all {
		if a.Name == name {
			p = &a
			break
		}
	}
	if p == nil {
		return profile{}, errors.New("not found")
	}
	return *p, nil
}
