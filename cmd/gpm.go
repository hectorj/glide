package cmd

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/glide/cfg"
)

// This file contains commands for working with GPM/GVP.

// HasGPMGodeps indicates whether a Godeps file exists.
func HasGPMGodeps(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	dir := cookoo.GetString("dir", "", p)
	path := filepath.Join(dir, "Godeps")
	_, err := os.Stat(path)
	return err == nil, nil
}

// GPMGodeps parses a GPM-flavored Godeps file.
//
// Params
// 	- dir (string): Directory root.
//
// Returns an []*cfg.Dependency
func GPMGodeps(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	dir := cookoo.GetString("dir", "", p)
	return parseGPMGodeps(dir)
}
func parseGPMGodeps(dir string) ([]*cfg.Dependency, error) {
	path := filepath.Join(dir, "Godeps")
	if i, err := os.Stat(path); err != nil {
		return []*cfg.Dependency{}, nil
	} else if i.IsDir() {
		Info("Godeps is a directory. This is probably a Godep project.\n")
		return []*cfg.Dependency{}, nil
	}
	Info("Found Godeps file.\n")

	buf := []*cfg.Dependency{}

	file, err := os.Open(path)
	if err != nil {
		return buf, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts, ok := parseGodepsLine(scanner.Text())
		if ok {
			dep := &cfg.Dependency{Name: parts[0]}
			if len(parts) > 1 {
				dep.Reference = parts[1]
			}
			buf = append(buf, dep)
		}
	}
	if err := scanner.Err(); err != nil {
		Warn("Scan failed: %s\n", err)
		return buf, err
	}

	return buf, nil
}

// GPMGodepsGit reads a Godeps-Git file for gpm-git.
func GPMGodepsGit(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	dir := cookoo.GetString("dir", "", p)
	path := filepath.Join(dir, "Godeps-Git")
	if _, err := os.Stat(path); err != nil {
		return []*cfg.Dependency{}, nil
	}
	Info("Found Godeps-Git file.\n")

	buf := []*cfg.Dependency{}

	file, err := os.Open(path)
	if err != nil {
		return buf, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts, ok := parseGodepsLine(scanner.Text())
		if ok {
			dep := &cfg.Dependency{Name: parts[1], Repository: parts[0]}
			if len(parts) > 2 {
				dep.Reference = parts[2]
			}
			buf = append(buf, dep)
		}
	}
	if err := scanner.Err(); err != nil {
		Warn("Scan failed: %s\n", err)
		return buf, err
	}

	return buf, nil
}

func parseGodepsLine(line string) ([]string, bool) {
	line = strings.TrimSpace(line)

	if len(line) == 0 || strings.HasPrefix(line, "#") {
		return []string{}, false
	}

	return strings.Fields(line), true
}
