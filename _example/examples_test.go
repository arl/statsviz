package example

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestExamples(t *testing.T) {
	examples := map[string]string{
		"chi":        "TODO",
		"default":    "http://localhost:8080/debug/statsviz/",
		"echo":       "TODO",
		"fasthttp":   "TODO",
		"fiber":      "TODO",
		"gin":        "TODO",
		"gorilla":    "TODO",
		"https":      "https://localhost:8080/debug/statsviz/",
		"iris":       "TODO",
		"middleware": "TODO",
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

	client := http.Client{}
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

			installScript := filepath.Join("install.sh")

			_, err = os.Stat(installScript)
			if err != nil && !os.IsNotExist(err) {
				t.Fatalf("stat %s: %v", installScript, err)
			}
			if os.IsExist(err) {
				cmd := exec.Command(installScript)
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Logf("command output: %s", out)
					t.Fatalf("exec %s: %v", installScript, err)
				}
			}

			stop, err := startStatsviz(ent.Name())
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

			resp, err := client.Get(url)
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

// startStatsviz runs go run 'dir', which starts the application opening a
// statsviz server. The returned function stops (kills) it.
func startStatsviz(dir string) (func() error, error) {
	binname := "." + string(os.PathSeparator) + dir + ".test"
	cmd := exec.Command("go", "build", "-o", binname)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("startStatsviz: go build output:\n%s\n", out)
		return nil, fmt.Errorf("go build failed: %v", err)
	}

	cmd = exec.Command(binname)
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
			fmt.Printf("startStatsviz: %s output: %s\n", binname, out)
		}
	}()
	if err := <-errc; err != nil {
		return nil, fmt.Errorf("startStatsviz: %v", err)
	}

	return cmd.Process.Kill, nil
}
