/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package challenger

import (
	context "context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicLifecycle(t *testing.T) {
	tests := []struct {
		name  string
		paths map[string]string
	}{{
		name: "single",
		paths: map[string]string{
			"/foo/bar": "baz",
		},
	}, {
		name: "multiple",
		paths: map[string]string{
			"/foo/bar":                      "baz",
			"/.well-known/acme/dsaflkjhsdf": "ugh",
			"/wtf":                          "is this",
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := New(context.Background())
			if err != nil {
				t.Fatalf("New() = %v", err)
			}

			// Check before registering
			for path := range test.paths {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rec := httptest.NewRecorder()
				c.ServeHTTP(rec, req)

				if got, want := rec.Result().StatusCode, http.StatusNotFound; got != want {
					t.Errorf("SeverHTTP(before register) = %d, wanted %d", got, want)
				}
			}

			// Check after registering
			for path, payload := range test.paths {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rec := httptest.NewRecorder()

				c.RegisterChallenge(path, payload)
				c.ServeHTTP(rec, req)
				if got, want := rec.Result().StatusCode, http.StatusOK; got != want {
					t.Errorf("SeverHTTP(after register) = %d, wanted %d", got, want)
				}

				body, err := io.ReadAll(rec.Result().Body)
				if err != nil {
					t.Errorf("ReadAll() = %v", err)
				} else if got, want := string(body), payload; got != want {
					t.Errorf("ReadAll() = %s, wanted %s", got, want)
				}
			}

			// Check after unregistering
			for path := range test.paths {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rec := httptest.NewRecorder()

				c.UnregisterChallenge(path)
				c.ServeHTTP(rec, req)

				if got, want := rec.Result().StatusCode, http.StatusNotFound; got != want {
					t.Errorf("SeverHTTP(after unregister) = %d, wanted %d", got, want)
				}
			}

		})
	}
}
