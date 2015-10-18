package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, fi := range files {
		name := fi.Name()
		if !strings.HasSuffix(name, ".png") {
			continue
		}

		f, err := os.Open("testdata/" + name)
		if err != nil {
			t.Errorf("failed to open file %q", "testdata/"+name)
			continue
		}

		want, err := ioutil.ReadFile("testdata/" + name + ".if")
		if err != nil {
			t.Errorf("failed to open file %q", "testdata/"+name+".if")
			continue
		}

		var buf bytes.Buffer
		cmd := exec.Command("go", "run", "any2if.go")
		cmd.Stdin = f
		cmd.Stdout = &buf
		err = cmd.Run()
		f.Close()
		if err != nil {
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				t.Error("any2if exited with non-zero exit status")
			} else {
				t.Errorf("failed to start any2if: %v", err)
			}
			continue
		}

		if !bytes.Equal(buf.Bytes(), want) {
			t.Errorf("invalid output for input file %q", "testdata/"+name)
		}
	}
}
