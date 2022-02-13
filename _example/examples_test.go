package example

import (
	"bytes"
	"context"
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
		"mux":        "TODO",
		"options":    "TODO",
	}

	ents, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{}

	for _, ent := range ents {
		if !ent.IsDir() {
			continue
		}
		t.Run(ent.Name(), func(t *testing.T) {
			url, ok := examples[ent.Name()]
			if !ok {
				t.Fatalf("url not specified for example %q", ent.Name())
				return
			}
			if url == "TODO" {
				t.Skipf("skipping example %s for now (TODO)", ent.Name())
				return
			}

			installScript := filepath.Join(ent.Name(), "install.sh")

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

			deadline := time.Now().Add(5 * time.Second)
			ctx, cancel := context.WithDeadline(context.Background(), deadline)
			defer cancel()

			// TODO(arl) do NOT use a context here, in case we need more control
			// of when to kill the server app.
			err := startStatsviz(ctx, ent.Name())
			if err != nil {
				t.Fatal(err)
			}
			// Let the time to the example application under test to listen on
			// the network interface.
			time.Sleep(1 * time.Second)
			client.Timeout = 3 * time.Second

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
// statsviz server.
func startStatsviz(ctx context.Context, dir string) error {
	binname := "./" + dir + ".test"
	cmd := exec.CommandContext(ctx, "go", "build", "-o", binname, "./"+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("startStatsviz: go build output:\n%s\n", out)
		return fmt.Errorf("go build failed: %v", err)
	}

	cmd = exec.CommandContext(ctx, binname)
	errc := make(chan error)
	go func() {
		outb := &bytes.Buffer{}
		cmd.Stderr = outb
		cmd.Stdout = outb
		errc <- cmd.Start()
		err := cmd.Wait()
		if err != nil {
			out := "<no output>"
			if outb.Len() > 0 {
				out = outb.String()
			}
			fmt.Printf("cmd.Wait failed: startStatsviz: %s output:\n%s\n", binname, out)
			return
		}
		if testing.Verbose() {
			fmt.Printf("startStatsviz: %s output:\n%s\n", binname, out)
		}
	}()
	if err := <-errc; err != nil {
		return fmt.Errorf("startStatsviz: %v", err)
	}

	return nil
}
