/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugins_test

import (
	"net/http"

	"github.com/apache/apisix-go-plugin-runner/tests/e2e/tools"
	"github.com/gavv/httpexpect/v2"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
)

var _ = ginkgo.Describe("ResponseRewrite Plugin", func() {
	table.DescribeTable("tries to test ResponseRewrite feature.",
		func(tc tools.HttpTestCase) {
			tools.RunTestCase(tc)
		},
		table.Entry("Config APISIX.", tools.HttpTestCase{
			Object: tools.GetA6CPExpect(),
			Method: http.MethodPut,
			Path:   "/apisix/admin/routes/1",
			Body: `{
				"uri":"/test/go/runner/say",
				"plugins":{
					"ext-plugin-post-resp":{
						"conf":[
							{
								"name":"response-rewrite",
								"value":"{\"headers\":{\"X-Server-Id\":\"9527\"},\"body\":\"response rewrite\"}"
							}
						]
					}
				},
				"upstream":{
					"nodes":{
						"web:8888":1
					},
					"type":"roundrobin"
				}
			}`,
			Headers:           map[string]string{"X-API-KEY": tools.GetAdminToken()},
			ExpectStatusRange: httpexpect.Status2xx,
		}),
		table.Entry("Should rewrite response.", tools.HttpTestCase{
			Object:       tools.GetA6DPExpect(),
			Method:       http.MethodGet,
			Path:         "/test/go/runner/say",
			ExpectBody:   []string{"response rewrite"},
			ExpectStatus: http.StatusOK,
			ExpectHeaders: map[string]string{
				"X-Resp-A6-Runner": "Go",
				"X-Server-Id":      "9527",
			},
		}),
	)
})
