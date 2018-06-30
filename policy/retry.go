package policy

import (
	"context"
	"net/http"
	"time"

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

var defaultRetryCodes = [...]int{
	http.StatusRequestTimeout,      // 408
	http.StatusTooManyRequests,     // 429
	http.StatusInternalServerError, // 500
	http.StatusBadGateway,          // 502
	http.StatusServiceUnavailable,  // 503
	http.StatusGatewayTimeout,      // 504
}

type SimpleRetryPolicyConfig struct {
	Attempts int
	Delay    time.Duration
}

func (cfg SimpleRetryPolicyConfig) getAttempts() int {
	if cfg.Attempts == 0 {
		return 3
	}
	return cfg.Attempts
}

func (cfg SimpleRetryPolicyConfig) getDelay() time.Duration {
	if cfg.Delay == 0 {
		return 1 * time.Second
	}
	return cfg.Delay
}

// returns true if statusCode is in the slice of defaultRetryCodes
func (cfg SimpleRetryPolicyConfig) shouldRetry(statusCode int) bool {
	for _, v := range defaultRetryCodes {
		if statusCode == v {
			return true
		}
	}
	return false
}

func NewSimpleRetryPolicyFactory(cfg SimpleRetryPolicyConfig) pipeline.FactoryFunc {
	return pipeline.FactoryFunc(func(next pipeline.Policy, po *pipeline.PolicyOptions) pipeline.PolicyFunc {
		return func(ctx context.Context, req pipeline.Request) (resp pipeline.Response, err error) {
			for try := 0; try < cfg.getAttempts(); try++ {
				reqCopy := req.Copy()
				if err = reqCopy.RewindBody(); err != nil {
					panic(err)
				}
				resp, err = next.Do(ctx, reqCopy)
				if err == nil && !cfg.shouldRetry(resp.Response().StatusCode) {
					return
				}
				time.Sleep(cfg.getDelay())
			}
			return
		}
	})
}
