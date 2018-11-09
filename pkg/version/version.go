// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build go1.8

package version

import (
	"log"
	"os"
	"strings"

	"rsc.io/goversion/version"
)

type Version struct {
	AppPath        string // /path/to/exe
	AppModPath     string // app go import path
	AppModVersion  string // app version, e.g. v0.1.0
	GoVersion      string // Go version (runtime.Version in the program)
	ModuleInfo     string // program's module information
	BoringCrypto   bool   // program uses BoringCrypto
	StandardCrypto bool   // program uses standard crypto (replaced by BoringCrypto)
	FIPSOnly       bool   // program imports "crypto/tls/fipsonly"
}

func GetVersion() *Version {
	var q = *pkgVersion
	return &q
}

func GetVersionString() string {
	return pkgVersion.AppModVersion
}

var pkgVersion = func() *Version {
	apppath, err := os.Executable()
	if err != nil {
		log.Fatalf("getAppPath failed: %v", err)
	}
	v, err := ReadVersion(apppath)
	if err != nil {
		log.Fatalf("ReadVersion failed: apppath = %s, err = %v", apppath, err)
	}
	return v
}()

func ReadVersion(apppath string) (*Version, error) {
	exeVersion, err := version.ReadExe(apppath)
	if err != nil {
		return nil, err
	}

	modPath, modVersion := parseVersionInfo(exeVersion)

	v := &Version{
		AppPath:        apppath,
		AppModPath:     modPath,
		AppModVersion:  modVersion,
		GoVersion:      exeVersion.Release,
		ModuleInfo:     exeVersion.ModuleInfo,
		BoringCrypto:   exeVersion.BoringCrypto,
		StandardCrypto: exeVersion.StandardCrypto,
		FIPSOnly:       exeVersion.FIPSOnly,
	}

	return v, nil
}

func parseVersionInfo(v version.Version) (modPath, modVersion string) {
	for _, line := range strings.Split(strings.TrimSpace(v.ModuleInfo), "\n") {
		row := strings.Split(line, "\t")
		if len(row) >= 2 && row[0] == "path" {
			modPath = row[1]
		}
		if len(row) >= 3 && row[0] == "mod" {
			modVersion = row[2]
		}
	}
	return
}
