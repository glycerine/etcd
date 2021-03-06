// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"bytes"
	"testing"
)

func isPermsEqual(a, b []*rangePerm) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if len(b) <= i {
			return false
		}

		if !bytes.Equal(a[i].begin, b[i].begin) || !bytes.Equal(a[i].end, b[i].end) {
			return false
		}
	}

	return true
}

func TestGetMergedPerms(t *testing.T) {
	tests := []struct {
		params []*rangePerm
		want   []*rangePerm
	}{
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}},
			[]*rangePerm{{[]byte("a"), []byte("b")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("b"), []byte("c")}},
			[]*rangePerm{{[]byte("a"), []byte("c")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("c")}, {[]byte("b"), []byte("d")}},
			[]*rangePerm{{[]byte("a"), []byte("d")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("b"), []byte("c")}, {[]byte("d"), []byte("e")}},
			[]*rangePerm{{[]byte("a"), []byte("c")}, {[]byte("d"), []byte("e")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("c"), []byte("d")}, {[]byte("e"), []byte("f")}},
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("c"), []byte("d")}, {[]byte("e"), []byte("f")}},
		},
		{
			[]*rangePerm{{[]byte("e"), []byte("f")}, {[]byte("c"), []byte("d")}, {[]byte("a"), []byte("b")}},
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("c"), []byte("d")}, {[]byte("e"), []byte("f")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("c"), []byte("d")}, {[]byte("a"), []byte("z")}},
			[]*rangePerm{{[]byte("a"), []byte("z")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("c"), []byte("d")}, {[]byte("a"), []byte("z")}, {[]byte("1"), []byte("9")}},
			[]*rangePerm{{[]byte("1"), []byte("9")}, {[]byte("a"), []byte("z")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("c"), []byte("d")}, {[]byte("a"), []byte("z")}, {[]byte("1"), []byte("a")}},
			[]*rangePerm{{[]byte("1"), []byte("z")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("a"), []byte("z")}, {[]byte("5"), []byte("6")}, {[]byte("1"), []byte("9")}},
			[]*rangePerm{{[]byte("1"), []byte("9")}, {[]byte("a"), []byte("z")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("b")}, {[]byte("b"), []byte("c")}, {[]byte("c"), []byte("d")}, {[]byte("d"), []byte("f")}, {[]byte("1"), []byte("9")}},
			[]*rangePerm{{[]byte("1"), []byte("9")}, {[]byte("a"), []byte("f")}},
		},
		// overlapping
		{
			[]*rangePerm{{[]byte("a"), []byte("f")}, {[]byte("b"), []byte("g")}},
			[]*rangePerm{{[]byte("a"), []byte("g")}},
		},
		// keys
		{
			[]*rangePerm{{[]byte("a"), []byte("")}, {[]byte("b"), []byte("")}},
			[]*rangePerm{{[]byte("a"), []byte("")}, {[]byte("b"), []byte("")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("")}, {[]byte("a"), []byte("c")}},
			[]*rangePerm{{[]byte("a"), []byte("c")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("")}, {[]byte("a"), []byte("c")}, {[]byte("b"), []byte("")}},
			[]*rangePerm{{[]byte("a"), []byte("c")}},
		},
		{
			[]*rangePerm{{[]byte("a"), []byte("")}, {[]byte("b"), []byte("c")}, {[]byte("b"), []byte("")}, {[]byte("c"), []byte("")}, {[]byte("d"), []byte("")}},
			[]*rangePerm{{[]byte("a"), []byte("")}, {[]byte("b"), []byte("c")}, {[]byte("c"), []byte("")}, {[]byte("d"), []byte("")}},
		},
		// duplicate ranges
		{
			[]*rangePerm{{[]byte("a"), []byte("f")}, {[]byte("a"), []byte("f")}},
			[]*rangePerm{{[]byte("a"), []byte("f")}},
		},
		// duplicate keys
		{
			[]*rangePerm{{[]byte("a"), []byte("")}, {[]byte("a"), []byte("")}, {[]byte("a"), []byte("")}},
			[]*rangePerm{{[]byte("a"), []byte("")}},
		},
	}

	for i, tt := range tests {
		result := mergeRangePerms(tt.params)
		if !isPermsEqual(result, tt.want) {
			t.Fatalf("#%d: result=%q, want=%q", i, result, tt.want)
		}
	}
}
