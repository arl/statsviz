package example

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestExamples(t *testing.T) {
	ents, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	for _, ent := range ents {
		if !ent.IsDir() {
			continue
		}

		// TODO(arl) just for now
		if ent.Name() != "default" {
			continue
		}

		t.Run(ent.Name(), func(t *testing.T) {
			popd, err := pushd(ent.Name())
			if err != nil {
				t.Fatal(err)
			}

			_, err = os.Stat("install.sh")
			if err != nil && !os.IsNotExist(err) {
				popd()
				t.Fatalf("install.sh: %v", err)
			}
			if os.IsExist(err) {
				cmd := exec.Command("./install.sh")
				out, err := cmd.CombinedOutput()
				if err != nil {
					popd()
					t.Logf("command output: %s", out)
					t.Fatalf("install.sh: %v", err)
				}
			}

			popd()

			// go build -o "./bin/$(basename "${example}")" ./"${example}"
			cmd := exec.Command("go", "run", "./"+ent.Name())
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("command output: %s", out)
				t.Fatalf("go run %s: %v", ent.Name(), err)
			}
		})
	}
}

func pushd(dir string) (popd func() error, err error) {
	pwd := ""
	if pwd, err = os.Getwd(); err != nil {
		return nil, fmt.Errorf("pushd: %v", err)
	}

	if err = os.Chdir(dir); err != nil {
		return nil, fmt.Errorf("pushd: %v", err)
	}

	return func() error { return os.Chdir(pwd) }, nil
}
