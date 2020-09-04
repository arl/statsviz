// +build dev

package statsviz

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// assets contains project assets located in current directory.
var assets http.FileSystem = dir("static")

// The rest of the code in this file is extracted from go/src/net/http/fs.go
// and slightly tweaked so that the livereload <script> tag gets injected in
// index.html.
// This is useful for development only, the generated assets won't contain
// the livereload <script> tag.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type dir string

// mapDirOpenError maps the provided non-nil error from opening name
// to a possibly better non-nil error. In particular, it turns OS-specific errors
// about opening files in non-directories into os.ErrNotExist. See Issue 18984.
func mapDirOpenError(originalErr error, name string) error {
	if os.IsNotExist(originalErr) || os.IsPermission(originalErr) {
		return originalErr
	}

	parts := strings.Split(name, string(filepath.Separator))
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		fi, err := os.Stat(strings.Join(parts[:i+1], string(filepath.Separator)))
		if err != nil {
			return originalErr
		}
		if !fi.IsDir() {
			return os.ErrNotExist
		}
	}
	return originalErr
}

// Open implements FileSystem using os.Open, opening files for reading rooted
// and relative to the directory d.
func (d dir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))

	f, err := os.Open(fullName)
	if err != nil {
		return nil, mapDirOpenError(err, fullName)
	}

	if name == "/index.html" {
		f, err := injectLiveReload(f)
		return f, err
	}

	return f, nil
}

func injectLiveReload(f http.File) (http.File, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	const livereload = `<script src="http://localhost:35729/livereload.js"></script>`
	const marker = `<!--DO NOT REMOVE: this comment gets replaced in development-->`
	b = bytes.Replace(b, []byte(marker), []byte(livereload), 1)

	return inmemFileFromFile(b, "index.html"), nil
}

// inmemFile is a in-memory implementatio of http.File. Its sole use in this
// package is to provide an http.File in which a livereload <script> has been
// injected to ease development. It hasn't been tested for other use cases and
// would probably won't work.
type inmemFile struct {
	b    []byte
	pos  int
	name string
}

func inmemFileFromFile(b []byte, name string) *inmemFile {
	return &inmemFile{b: b, name: name}
}

func (f *inmemFile) Read(p []byte) (n int, err error) {
	n = len(f.b) - f.pos
	if len(p) < n {
		n = len(p)
	}

	n = copy(p, f.b[f.pos:f.pos+n])
	f.pos += n
	return n, nil
}

func (f *inmemFile) Seek(offset int64, whence int) (int64, error) {
	off := int(offset)
	switch whence {
	case 0:
		f.pos = off
	case 1:
		f.pos += off
	case 2:
		f.pos = len(f.b) + off
	}
	return int64(f.pos), nil
}

func (f *inmemFile) Close() error                             { return nil }
func (f *inmemFile) Readdir(count int) ([]os.FileInfo, error) { return nil, nil }
func (f *inmemFile) Stat() (os.FileInfo, error)               { return inmemFileInfo{f}, nil }

type inmemFileInfo struct{ f *inmemFile }

// Implements os.FileInfo
func (s inmemFileInfo) Name() string       { return s.f.name }
func (s inmemFileInfo) Size() int64        { return int64(len(s.f.b)) }
func (s inmemFileInfo) Mode() os.FileMode  { return os.ModeTemporary }
func (s inmemFileInfo) ModTime() time.Time { return time.Time{} }
func (s inmemFileInfo) IsDir() bool        { return false }
func (s inmemFileInfo) Sys() interface{}   { return nil }
