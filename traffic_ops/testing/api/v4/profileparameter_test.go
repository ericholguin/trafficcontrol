package v4

/*

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

import (
	"net/http"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestProfileParameter(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Parameters, Profiles}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.PostProfileParam]{
			"POST": {
				"OK when VALID request": {
					ClientSession: TOSession,
					RequestBody: tc.PostProfileParam{
						ProfileID: util.Ptr(int64(GetProfileID(t, "ATS_EDGE_TIER_CACHE")())),
						ParamIDs: util.Ptr([]int64{
							int64(GetParameterID(t, "CONFIG proxy.config.allocator.enable_reclaim", "records.config", "INT 0")()),
							int64(GetParameterID(t, "CONFIG proxy.config.allocator.max_overage", "records.config", "INT 3")()),
						}),
					},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "POST":
						t.Run(name, func(t *testing.T) {
							alerts, reqInf, err := testCase.ClientSession.CreateProfileWithMultipleParameters(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, alerts, err)
							}
						})
					}
				}
			})
		}
	})
}
