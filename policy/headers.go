package policy

import (
	"context"
	"net/http"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

type CtxKeyReqHeaders struct{}

func NewReqHeadersFactory() pipeline.Factory {
	return pipeline.FactoryFunc(func(next pipeline.Policy, po *pipeline.PolicyOptions) pipeline.PolicyFunc {
		return func(ctx context.Context, req pipeline.Request) (pipeline.Response, error) {
			headers := ctx.Value(CtxKeyReqHeaders{})
			if headers != nil {
				httpHeader := headers.(http.Header)
				for k, vs := range httpHeader {
					for _, v := range vs {
						req.Header.Add(k, v)
					}
				}
			}
			return next.Do(ctx, req)
		}
	})
}
