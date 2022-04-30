// Adapted from https://github.com/uber-go/atomic
// Original copyright below (MIT license):
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package statsviz

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestFloat64(t *testing.T) {
	atom := newFloat64(4.2)
	if got := atom.load(); got != 4.2 {
		t.Errorf("load() = %v, want %v", got, 4.2)
	}

	atom.Store(0.5)
	if got := atom.load(); got != 0.5 {
		t.Errorf("store failed. load() = %v, want %v", got, 0.5)
	}

	atom.Store(42.0)
	if got := atom.load(); got != 42.0 {
		t.Errorf("store failed. load() = %v, want %v", got, 42.0)
	}

	t.Run("JSON/Marshal", func(t *testing.T) {
		atom.Store(42.5)
		buf, err := json.Marshal(atom)
		if err != nil {
			t.Fatal(err)
		}
		if want := []byte("42.5"); !bytes.Equal(buf, want) {
			t.Errorf("json.Marshal = %q, want %q", buf, want)
		}
	})

	t.Run("JSON/Unmarshal", func(t *testing.T) {
		err := json.Unmarshal([]byte("40.5"), &atom)
		if err != nil {
			t.Fatal(err)
		}
		if got := atom.load(); got != 40.5 {
			t.Errorf("json.Unmarshal failed. atom = %v, want %v", got, 42.0)
		}
	})

	t.Run("JSON/Unmarshal/Error", func(t *testing.T) {
		err := json.Unmarshal([]byte("\"40.5\""), &atom)
		if err == nil {
			t.Errorf("want error, got nil")
		}
	})
}
