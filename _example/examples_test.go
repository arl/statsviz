package example

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestExamples(t *testing.T) {
	p := testscript.Params{
		Dir: "testdata",
		Setup: func(env *testscript.Env) error {
			// We want to run scripts with the local version of Statsviz.
			// Provide scripts with statsviz root dir so we can use a
			// 'go mod -edit replace' directive.
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			env.Setenv("STATSVIZ_ROOT", filepath.Dir(wd))
			return nil
		},
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"checkui": checkui,
		},
	}

	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	testscript.Run(t, p)
}

// checkui requests statsviz url from a script.
// In a script, run it with:
//
//	checkui url [basic_auth_user basic_auth_pwd]
func checkui(ts *testscript.TestScript, neg bool, args []string) {
	if len(args) != 1 && len(args) != 3 {
		ts.Fatalf(`checkui: wrong number of arguments. Call with "checkui URL [BASIC_USER BASIC_PWD]`)
	}
	u := args[0]
	ts.Logf("checkui: loading web page %s", args[0])

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		ts.Fatalf("checkui: bad request: %v", err)
	}
	if len(args) == 3 {
		ts.Logf("checkui: setting basic auth")
		req.SetBasicAuth(args[1], args[2])
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 1 * time.Second,
	}

	// Let 1 second for the server to start and listen.
	time.Sleep(1 * time.Second)

	resp, err := client.Do(req)
	ts.Check(err)

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	ts.Check(err)

	want := []byte(`id="plots"`)
	if !bytes.Contains(body, want) {
		ts.Fatalf("checkui: response body doesn't contain %s\n\nbody;\n\n%s", want, body)
	}
}
