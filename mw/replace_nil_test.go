// Copyright © 2024 OpenIM open source community. All rights reserved.
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

package mw

import (
	"encoding/json"
	"fmt"
	"testing"
)

type A struct {
	B  *B
	BB B
	BS []*B
	C  []int
	D  map[string]string
	E  interface{}
	F  *int
}

type B struct {
	D *C
	E []int
}

type C struct {
}

type D struct {
	sb  string
	nt  []C
	ssb *A
}

func TestReplaceNil(t *testing.T) {
	a := &A{}
	k := any(a)
	ReplaceNil(&k)
	//printJson(k)
	//printJson(repl(k))
	// {"B":null,"BB":{"D":null,"E":[]},"C":[],"D":{},"E":null,"F":null}

	var b *A
	k = any(b)
	ReplaceNil(&k)
	//printJson(repl(k))
	// {}

	i := 5
	c := &A{
		B: nil,
		BB: B{
			D: &C{},
			E: []int{1, 2, 5, 3, 6},
		},
		C: []int{1, 1, 1},
		D: map[string]string{
			"a": "A",
			"b": "B",
		},
		E: map[int]int{
			1: 11,
			2: 22,
		},
		F: &i,
	}
	k = any(c)
	ReplaceNil(&k)
	printJson(k)
	// {"B":null,"BB":{"D":{},"E":[1,2,5,3,6]},"C":[1,1,1],"D":{"a":"A","b":"B"},"E":{"1":11,"2":22},"F":5}

	dd := &D{
		sb:  "fhldsa",
		nt:  []C{},
		ssb: &A{},
	}
	k = any(dd)
	ReplaceNil(&k)
	printJson(k)
	// {}
}

func printJson(data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
