//go:build !windows

package main

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

// TERMOTE_BIND=unix:<path> listens on a socket only this user may open,
// replaces one a dead server left, and never takes a live one or a file.
func TestListenUnixSocket(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "t.sock")
	ln, err := listenServer("unix:"+path, "7680")
	if err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil || fi.Mode()&os.ModeSocket == 0 || fi.Mode().Perm() != 0o600 {
		t.Fatalf("socket %v, %v", fi, err)
	}
	go func() {
		if c, err := ln.Accept(); err == nil {
			c.Close()
		}
	}()
	if _, err := listenServer("unix:"+path, "7680"); err == nil {
		t.Fatal("took a live server's socket")
	}
	ln.Close()

	// What a killed server leaves: a socket file nobody answers on.
	stale := filepath.Join(dir, "stale.sock")
	l, err := net.Listen("unix", stale)
	if err != nil {
		t.Fatal(err)
	}
	l.(*net.UnixListener).SetUnlinkOnClose(false)
	l.Close()
	if ln, err := listenServer("unix:"+stale, ""); err != nil {
		t.Fatalf("stale socket: %v", err)
	} else {
		ln.Close()
	}

	file := filepath.Join(dir, "file")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, bind := range []string{"unix:" + file, "unix:", "unix:rel.sock"} {
		if ln, err := listenServer(bind, ""); err == nil {
			ln.Close()
			t.Errorf("%s: listened", bind)
		}
	}
	if !isUnixBind("unix:/x") || isUnixBind("127.0.0.1") {
		t.Error("isUnixBind")
	}
}
