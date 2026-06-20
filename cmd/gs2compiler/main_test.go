package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"gs2compiler", "--help"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(stdout.String(), "GS2 Script Compiler") {
		t.Fatalf("missing help text: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr: %q", stderr.String())
	}
}

func TestRunNoInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"gs2compiler"}, &stdout, &stderr)
	if code == 0 {
		t.Fatal("expected failure")
	}
	if !strings.Contains(stderr.String(), "No input file specified") {
		t.Fatalf("stderr: %q", stderr.String())
	}
}

func TestRunDefaultOutput(t *testing.T) {
	dir := localTestDir(t)
	src := filepath.Join(dir, "sample.gs2")
	if err := os.WriteFile(src, []byte(`function onCreated() { temp.a = 1; }`), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"gs2compiler", src}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit %d stderr %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Compilation successful") {
		t.Fatalf("stdout: %q", stdout.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "sample.gs2bc")); err != nil {
		t.Fatal(err)
	}
}

func TestRunPositionalOutput(t *testing.T) {
	dir := localTestDir(t)
	src := filepath.Join(dir, "sample.gs2")
	out := filepath.Join(dir, "out.gs2bc")
	if err := os.WriteFile(src, []byte(`function onCreated() { temp.a = 1; }`), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"gs2compiler", src, out}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit %d stderr %q", code, stderr.String())
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
}

func TestRunLongOutput(t *testing.T) {
	dir := localTestDir(t)
	src := filepath.Join(dir, "sample.gs2")
	out := filepath.Join(dir, "out.gs2bc")
	if err := os.WriteFile(src, []byte(`function onCreated() { temp.a = 1; }`), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"gs2compiler", "--output", out, src}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit %d stderr %q", code, stderr.String())
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatal(err)
	}
}

func TestRunDirectoryMode(t *testing.T) {
	dir := localTestDir(t)
	if err := os.WriteFile(filepath.Join(dir, "a.gs2"), []byte(`function onCreated() { temp.a = 1; }`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte(`function onCreated() { temp.b = 2; }`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "c.md"), []byte(`function onCreated() { temp.c = 3; }`), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"gs2compiler", dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit %d stderr %q", code, stderr.String())
	}
	for _, name := range []string{"a.gs2bc", "b.gs2bc"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "c.gs2bc")); !os.IsNotExist(err) {
		t.Fatalf("unexpected c.gs2bc err=%v", err)
	}
}

func TestRunMultiFileMode(t *testing.T) {
	dir := localTestDir(t)
	a := filepath.Join(dir, "a.gs2")
	b := filepath.Join(dir, "b.gs2")
	c := filepath.Join(dir, "c.gs2")
	if err := os.WriteFile(a, []byte(`function onCreated() { temp.a = 1; }`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte(`function onCreated() { temp.b = 2; }`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(c, []byte(`function onCreated() { temp.c = 3; }`), 0644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{"gs2compiler", a, b, c}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit %d stderr %q", code, stderr.String())
	}
	for _, name := range []string{"a.gs2bc", "b.gs2bc", "c.gs2bc"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
}

func localTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp(".", ".gs2compiler-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}
