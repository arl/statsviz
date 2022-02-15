//go:build linux
// +build linux

package example

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestExamples(t *testing.T) {
	examples := map[string]string{
		"chi":        "http://localhost:8080/debug/statsviz/",
		"default":    "http://localhost:8080/debug/statsviz/",
		"echo":       "http://localhost:8080/debug/statsviz/",
		"fasthttp":   "http://localhost:8080/debug/statsviz/",
		"fiber":      "http://localhost:8080/debug/statsviz/",
		"gin":        "http://localhost:8080/debug/statsviz/",
		"gorilla":    "http://localhost:8080/debug/statsviz/",
		"https":      "https://localhost:8080/debug/statsviz/",
		"iris":       "http://localhost:8080/debug/statsviz/",
		"middleware": "http://localhost:8080/debug/statsviz/",
		"mux":        "http://localhost:8080/debug/statsviz/",
		"options":    "http://localhost:8080/foo/bar/",
	}

	ents, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	for _, ent := range ents {
		if !ent.IsDir() {
			continue
		}
		t.Run(ent.Name(), func(t *testing.T) {
			// Some examples require files in their specific package, so by
			// convention, we'll always be using the package as working
			// directory.
			if err := os.Chdir(filepath.Join(wd, ent.Name())); err != nil {
				t.Fatal(err)
			}

			url, ok := examples[ent.Name()]
			if !ok {
				t.Fatalf("url not specified for example %q", ent.Name())
				return
			}
			if url == "TODO" {
				t.Skipf("skipping example %s for now (TODO)", ent.Name())
				return
			}

			stop, err := gorun()
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err := stop()
				if err != nil {

				}
			}()

			// Let the application we're testing the time to start listening for
			// HTTP connections.
			time.Sleep(1 * time.Second)
			client.Timeout = 1 * time.Second

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Fatalf("bad requets: %v", err)
			}
			if strings.Contains(t.Name(), "middleware") {
				req.SetBasicAuth("hello", "world")
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("HTTP get %s: %v", url, err)
			}

			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				t.Fatal(err)
			}

			for _, s := range []string{"Heap", "Goroutines", "GC / CPU fraction"} {
				if !bytes.Contains(body, []byte(s)) {
					t.Errorf("body doesn't contain %s", s)
				}
			}
			if t.Failed() {
				fmt.Printf("body:%s\n", body)
			}
		})
	}
}

// gorun runs 'go run .' and either returns an error immediately if there was
// one, or returns a function that can be called at any moment in order to stop
// both 'go run' and the started process.
func gorun() (stop func() error, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("go run: %v", err)
	}

	cmd := exec.Command("go", "run", ".")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	errc := make(chan error)
	go func() {
		outb := &bytes.Buffer{}
		cmd.Stderr = outb
		cmd.Stdout = outb
		errc <- cmd.Start()
		// Ignore error since we kill the process ourselves.
		_ = cmd.Wait()

		if testing.Verbose() {
			out := "<no output>"
			if outb.Len() > 0 {
				out = outb.String()
			}
			fmt.Printf("go run %s, output: %s\n", wd, out)
		}
	}()

	if err := <-errc; err != nil {
		return nil, fmt.Errorf("go run %s: %v", wd, err)
	}

	stop = func() error {
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			return fmt.Errorf("go run %s: can't kill pid=%v %v", wd, cmd.Process.Pid, err)
		}
		return nil
	}
	return stop, nil
}
