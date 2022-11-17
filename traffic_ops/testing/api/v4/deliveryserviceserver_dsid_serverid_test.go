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
	"net/url"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestDeliveryServicesDSIDServerID(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServiceCategories, DeliveryServices, DeliveryServiceServerAssignments}, func() {

		methodTests := utils.TestCase[client.Session, client.RequestOptions, struct{}]{
			"DELETE": {
				"OK when VALID REQUEST": {
					EndpointId:    GetDeliveryServiceID(t, "ds-top"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"server": {strconv.Itoa(GetServerID(t, "denver-mso-org-01")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when REMOVING ONLY EDGE SERVER ASSIGNMENT": {
					EndpointId:    GetDeliveryServiceID(t, "test-ds-server-assignments"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"server": {strconv.Itoa(GetServerID(t, "test-ds-server-assignments")())}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
				"BAD REQUEST when REMOVING ONLY ORIGIN SERVER ASSIGNMENT": {
					EndpointId:    GetDeliveryServiceID(t, "test-ds-server-assignments"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"server": {strconv.Itoa(GetServerID(t, "test-mso-org-01")())}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusConflict)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					var serverId int
					var err error
					if val, ok := testCase.RequestOpts.QueryParameters["server"]; ok {
						serverId, err = strconv.Atoi(val[0])
						assert.RequireNoError(t, err, "Expected no error converting string to integer.")
					}

					switch method {
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryServiceServer(testCase.EndpointId(), serverId, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp, err)
							}
						})
					}
				}
			})
		}
	})
}
