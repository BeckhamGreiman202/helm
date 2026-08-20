/*
Copyright The Helm Authors.

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

package values

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
)

func TestMergeValues(t *testing.T) {
	p := getter.All(&cli.EnvSettings{})
	options := Options{
		ValueFiles: []string{"testdata/vals.yaml"},
	}
	vals, err := options.MergeValues(p)
	if err != nil {
		t.Fatal(err)
	}

	if vals["name"] != "value" {
		t.Errorf("Expected name to be value, got %q", vals["name"])
	}
}

func TestMergeValuesTruncation(t *testing.T) {
	p := getter.All(&cli.EnvSettings{})
	sizes := []int{4096, 8192, 12288}
	for _, size := range sizes {
		contentSuffix := "\nkey: val"
		paddingSize := size - len(contentSuffix)
		if paddingSize < 0 {
			t.Fatalf("size %d is too small", size)
		}
		padding := make([]byte, paddingSize)
		for i := range padding {
			padding[i] = '#'
		}
		content := append(padding, []byte(contentSuffix)...)

		tmpFile, err := os.CreateTemp("", "values-*.yaml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatal(err)
		}

		options := Options{
			ValueFiles: []string{tmpFile.Name()},
		}
		vals, err := options.MergeValues(p)
		if err != nil {
			t.Fatalf("failed to merge values for size %d: %v", size, err)
		}

		val, ok := vals["key"]
		if !ok {
			t.Fatalf("key not found in values for size %d", size)
		}
		if valStr, ok := val.(string); !ok || valStr != "val" {
			t.Fatalf("expected 'val', got %v for size %d", val, size)
		}
	}
}
