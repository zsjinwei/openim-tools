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

package a2r

import (
	"context"
	"github.com/zsjinwei/openim-tools/mw"
	"google.golang.org/grpc"
)

func NewNilReplaceOption[A, B, C any](_ func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error)) *Option[A, B] {
	return &Option[A, B]{
		RespAfter: respNilReplace[B],
	}
}

// respNilReplace replaces nil maps and slices in the resp object and initializing them.
func respNilReplace[T any](data *T) error {
	mw.ReplaceNil(data)
	return nil
}

// ------------------------------------------------------------------------------------------------------------------
