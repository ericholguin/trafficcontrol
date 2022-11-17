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
	"time"

	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/assert"
	"github.com/apache/trafficcontrol/traffic_ops/testing/api/utils"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
	client "github.com/apache/trafficcontrol/traffic_ops/v4-client"
)

func TestDeliveryServices(t *testing.T) {
	WithObjs(t, []TCObj{CDNs, Types, Tenants, Users, Parameters, Profiles, Statuses, Divisions, Regions, PhysLocations, CacheGroups, Servers, Topologies, ServerCapabilities, ServiceCategories, DeliveryServices, ServerServerCapabilities, DeliveryServicesRequiredCapabilities, DeliveryServiceServerAssignments}, func() {

		currentTime := time.Now().UTC().Add(-15 * time.Second)
		currentTimeRFC := currentTime.Format(time.RFC1123)
		tomorrow := currentTime.AddDate(0, 0, 1).Format(time.RFC1123)

		tenant1UserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "tenant1user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant2UserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "tenant2user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant3UserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "tenant3user", "pa$$word", Config.Default.Session.TimeoutInSecs)
		tenant4UserSession := utils.CreateV4Session(t, Config.TrafficOps.URL, "tenant4user", "pa$$word", Config.Default.Session.TimeoutInSecs)

		methodTests := utils.TestCase[client.Session, client.RequestOptions, tc.DeliveryServiceV4]{
			"GET": {
				"NOT MODIFIED when NO CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {tomorrow}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusNotModified)),
				},
				"OK when VALID request": {
					ClientSession: TOSession,
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when VALID ACCESSIBLETO parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"accessibleTo": {strconv.Itoa(GetTenantID(t, "root")())}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1)),
				},
				"OK when ACTIVE=TRUE": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"active": {"true"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Active": true})),
				},
				"OK when ACTIVE=FALSE": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"active": {"false"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Active": false})),
				},
				"OK when VALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {strconv.Itoa(GetCDNID(t, "cdn1")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"CDNName": "cdn1"})),
				},
				"OK when VALID LOGSENABLED parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"logsEnabled": {"false"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"LogsEnabled": false})),
				},
				"OK when VALID PROFILE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"profile": {strconv.Itoa(GetProfileID(t, "ATS_EDGE_TIER_CACHE")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"ProfileName": "ATS_EDGE_TIER_CACHE"})),
				},
				"OK when VALID SERVICECATEGORY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"serviceCategory": {"serviceCategory1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"ServiceCategory": "serviceCategory1"})),
				},
				"OK when VALID TENANT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"tenant": {strconv.Itoa(GetTenantID(t, "tenant1")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Tenant": "tenant1"})),
				},
				"OK when VALID TOPOLOGY parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"topology": {"mso-topology"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Topology": "mso-topology"})),
				},
				"OK when VALID TYPE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"type": {strconv.Itoa(GetTypeID(t, "HTTP")())}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseLengthGreaterOrEqual(1),
						validateDSExpectedFields(map[string]interface{}{"Type": tc.DSTypeHTTP})),
				},
				"OK when VALID XMLID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds1"}}},
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1),
						validateDSExpectedFields(map[string]interface{}{"XMLID": "ds1"})),
				},
				"EMPTY RESPONSE when INVALID ACCESSIBLETO parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"accessibleTo": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID CDN parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"cdn": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID PROFILE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"profile": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID TENANT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"tenant": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID TYPE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"type": {"10000"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INVALID XMLID parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"invalid_xml_id"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"FIRST RESULT when LIMIT=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("limit")),
				},
				"SECOND RESULT when LIMIT=1 OFFSET=1": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "offset": {"1"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("offset")),
				},
				"SECOND RESULT when LIMIT=1 PAGE=2": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"orderby": {"id"}, "limit": {"1"}, "page": {"2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validatePagination("page")),
				},
				"BAD REQUEST when INVALID LIMIT parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"-2"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID OFFSET parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "offset": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when INVALID PAGE parameter": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"limit": {"1"}, "page": {"0"}}},
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"VALID when SORTORDER param is DESC": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"sortOrder": {"desc"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), validateDescSort()),
				},
				"OK when PARENT TENANT reads DS of INACTIVE CHILD TENANT": {
					ClientSession: tenant1UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(1)),
				},
				"EMPTY RESPONSE when DS BELONGS to TENANT but PARENT TENANT is INACTIVE": {
					ClientSession: tenant3UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds3"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when INACTIVE TENANT reads DS of SAME TENANCY": {
					ClientSession: tenant2UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when TENANT reads DS OUTSIDE TENANCY": {
					ClientSession: tenant4UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds3"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"EMPTY RESPONSE when CHILD TENANT reads DS of PARENT TENANT": {
					ClientSession: tenant3UserSession,
					RequestOpts:   client.RequestOptions{QueryParameters: url.Values{"xmlId": {"ds2"}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK), utils.ResponseHasLength(0)),
				},
				"OK when CHANGES made": {
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfModifiedSince: {currentTimeRFC}}},
					Expectations:  utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
			},
			"POST": {
				"CREATED when VALID request WITH GEO LIMIT COUNTRIES": {
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								GeoLimit:          util.Ptr(2),
								GeoLimitCountries: util.Ptr("US,CA"),
								XMLID:             util.Ptr("geolimit-test"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated), utils.ResponseHasLength(1),
						validateDSExpectedFields(map[string]interface{}{"GeoLimitCountries": tc.GeoLimitCountriesType{"US", "CA"}})),
				},
				"BAD REQUEST when using LONG DESCRIPTION 2 and 3 fields": {
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								LongDesc1: util.Ptr("long desc 1"),
								LongDesc2: util.Ptr("long desc 2"),
								XMLID:     util.Ptr("ld1-ld2-test"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when XMLID left EMPTY": {
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr(""),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when XMLID is NIL": {
					ClientSession: TOSession,
					RequestBody: tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								CDNID:               util.Ptr(GetCDNID(t, "cdn1")()),
								LongDesc:            util.Ptr("something different"),
								MaxDNSAnswers:       util.Ptr(164598),
								Active:              util.Ptr(false),
								DisplayName:         util.Ptr("newds2displayname"),
								DSCP:                util.Ptr(41),
								GeoProvider:         util.Ptr(1),
								GeoLimit:            util.Ptr(1),
								InitialDispersion:   util.Ptr(2),
								IPV6RoutingEnabled:  util.Ptr(false),
								LogsEnabled:         util.Ptr(false),
								MissLat:             util.Ptr(42.881944),
								MissLong:            util.Ptr(-88.627778),
								MultiSiteOrigin:     util.Ptr(true),
								OrgServerFQDN:       util.Ptr("http://origin.example.net"),
								Protocol:            util.Ptr(2),
								RoutingName:         util.Ptr("ccr-ds2"),
								QStringIgnore:       util.Ptr(0),
								RegionalGeoBlocking: util.Ptr(true),
								TypeID:              util.Ptr(GetTypeID(t, "HTTP")()),
								XMLID:               nil,
							},
							Regional: true,
						},
					},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								Topology: util.Ptr("topology-doesnt-exist"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when creating STEERING DS with TLS VERSIONS": {
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								TypeID: util.Ptr(GetTypeId(t, "STEERING")),
								XMLID:  util.Ptr("test-TLS-creation-steering"),
							},
							TLSVersions: []string{"1.1"},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when creating HTTP DS with TLS VERSIONS": {
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("test-TLS-creation-http"),
							},
							TLSVersions: []string{"1.1"},
						},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusCreated), utils.ResponseHasLength(1)),
				},
				"BAD REQUEST when creating DS with TENANCY NOT THE SAME AS CURRENT TENANT": {
					ClientSession: tenant4UserSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								TenantID: util.Ptr(GetTenantID(t, "tenant3")()),
								XMLID:    util.Ptr("test-tenancy"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden), utils.ResponseHasLength(0)),
				},
			},
			"PUT": {
				"BAD REQUEST when using LONG DESCRIPTION 2 and 3 fields": {
					EndpointId: GetDeliveryServiceID(t, "ds1"), ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								LongDesc1: util.Ptr("long desc 1"),
								LongDesc2: util.Ptr("long desc 2"),
								XMLID:     util.Ptr("ds1"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when VALID request": {
					EndpointId:    GetDeliveryServiceID(t, "ds2"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV31: tc.DeliveryServiceFieldsV31{
								MaxRequestHeaderBytes: util.Ptr(131080),
							},
							DeliveryServiceFieldsV14: tc.DeliveryServiceFieldsV14{
								MaxOriginConnections: util.Ptr(100),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								LongDesc:            util.Ptr("something different"),
								MaxDNSAnswers:       util.Ptr(164598),
								Active:              util.Ptr(false),
								DisplayName:         util.Ptr("newds2displayname"),
								DSCP:                util.Ptr(41),
								GeoLimit:            util.Ptr(1),
								InitialDispersion:   util.Ptr(2),
								IPV6RoutingEnabled:  util.Ptr(false),
								LogsEnabled:         util.Ptr(false),
								MissLat:             util.Ptr(42.881944),
								MissLong:            util.Ptr(-88.627778),
								MultiSiteOrigin:     util.Ptr(true),
								OrgServerFQDN:       util.Ptr("http://origin.example.net"),
								Protocol:            util.Ptr(2),
								RoutingName:         util.Ptr("ccr-ds2"),
								QStringIgnore:       util.Ptr(0),
								RegionalGeoBlocking: util.Ptr(true),
								XMLID:               util.Ptr("ds2"),
							},
							Regional: true,
						},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateDSExpectedFields(map[string]interface{}{"MaxRequestHeaderSize": 131080,
							"LongDesc": "something different", "MaxDNSAnswers": 164598, "MaxOriginConnections": 100,
							"Active": false, "DisplayName": "newds2displayname", "DSCP": 41, "GeoLimit": 1,
							"InitialDispersion": 2, "IPV6RoutingEnabled": false, "LogsEnabled": false, "MissLat": 42.881944,
							"MissLong": -88.627778, "MultiSiteOrigin": true, "OrgServerFQDN": "http://origin.example.net",
							"Protocol": 2, "QStringIgnore": 0, "RegionalGeoBlocking": true,
						})),
				},
				"BAD REQUEST when INVALID REMAP TEXT": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								RemapText: util.Ptr("@plugin=tslua.so @pparam=/opt/trafficserver/etc/trafficserver/remapPlugin1.lua\nline2"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when MISSING SLICE PLUGIN SIZE": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								RangeRequestHandling: util.Ptr(3),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE SET with INVALID RANGE REQUEST SETTING": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV15: tc.DeliveryServiceFieldsV15{
								RangeSliceBlockSize: util.Ptr(262144),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								RangeRequestHandling: util.Ptr(1),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO SMALL": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV15: tc.DeliveryServiceFieldsV15{
								RangeSliceBlockSize: util.Ptr(0),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								RangeRequestHandling: util.Ptr(3),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when SLICE PLUGIN SIZE TOO LARGE": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV15: tc.DeliveryServiceFieldsV15{
								RangeSliceBlockSize: util.Ptr(40000000),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								RangeRequestHandling: util.Ptr(3),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to CLIENT STEERING DS": {
					EndpointId:    GetDeliveryServiceID(t, "ds-client-steering"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								Topology: util.Ptr("mso-topology"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								TypeID: util.Ptr(GetTypeID(t, "CLIENT_STEERING")()),
								XMLID:  util.Ptr("ds-client-steering"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when TOPOLOGY DOESNT EXIST": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								Topology: util.Ptr("topology-doesnt-exist"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds1"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to DS with DS REQUIRED CAPABILITY": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								Topology: util.Ptr("top-for-ds-req"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds1"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when ADDING TOPOLOGY to DS when NO CACHES in SAME CDN as DS": {
					EndpointId:    GetDeliveryServiceID(t, "top-ds-in-cdn2"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								Topology: util.Ptr("top-with-caches-in-cdn1"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								CDNID: util.Ptr(GetCDNID(t, "cdn2")()),
								XMLID: util.Ptr("top-ds-in-cdn2"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when REMOVING TOPOLOGY": {
					EndpointId:    GetDeliveryServiceID(t, "ds-based-top-with-no-mids"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								Topology: nil,
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds-based-top-with-no-mids"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when DS with TOPOLOGY updates HEADER REWRITE FIELDS": {
					EndpointId:    GetDeliveryServiceID(t, "ds-top"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								FirstHeaderRewrite: util.Ptr("foo"),
								InnerHeaderRewrite: util.Ptr("bar"),
								LastHeaderRewrite:  util.Ptr("baz"),
								Topology:           util.Ptr("mso-topology"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds-top"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"BAD REQUEST when DS with NO TOPOLOGY updates HEADER REWRITE FIELDS": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								FirstHeaderRewrite: util.Ptr("foo"),
								InnerHeaderRewrite: util.Ptr("bar"),
								LastHeaderRewrite:  util.Ptr("baz"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds1"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when DS with TOPOLOGY updates LEGACY HEADER REWRITE FIELDS": {
					EndpointId:    GetDeliveryServiceID(t, "ds-top"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								Topology: util.Ptr("mso-topology"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								EdgeHeaderRewrite: util.Ptr("foo"),
								MidHeaderRewrite:  util.Ptr("bar"),
								XMLID:             util.Ptr("ds-top"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"OK when DS with NO TOPOLOGY updates LEGACY HEADER REWRITE FIELDS": {
					EndpointId:    GetDeliveryServiceID(t, "ds2"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								ProfileID:         util.Ptr(GetProfileID(t, "ATS_EDGE_TIER_CACHE")()),
								EdgeHeaderRewrite: util.Ptr("foo"),
								MidHeaderRewrite:  util.Ptr("bar"),
								RoutingName:       util.Ptr("ccr-ds2"),
								XMLID:             util.Ptr("ds2"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK)),
				},
				"OK when UPDATING MINOR VERSION FIELDS": {
					EndpointId:    GetDeliveryServiceID(t, "ds-test-minor-versions"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV14: tc.DeliveryServiceFieldsV14{
								ConsistentHashQueryParams: []string{"d", "e", "f"},
								ConsistentHashRegex:       util.Ptr("foo"),
								MaxOriginConnections:      util.Ptr(500),
							},
							DeliveryServiceFieldsV13: tc.DeliveryServiceFieldsV13{
								DeepCachingType:   util.Ptr(tc.DeepCachingTypeNever),
								FQPacingRate:      util.Ptr(41),
								SigningAlgorithm:  util.Ptr("uri_signing"),
								TRRequestHeaders:  util.Ptr("X-ooF\nX-raB"),
								TRResponseHeaders: util.Ptr("Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8"),
								Tenant:            util.Ptr("tenant1"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								RoutingName: util.Ptr("cdn"),
								TenantID:    util.Ptr(GetTenantID(t, "tenant1")()),
								XMLID:       util.Ptr("ds-test-minor-versions"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.NoError(), utils.HasStatus(http.StatusOK),
						validateDSExpectedFields(map[string]interface{}{"ConsistentHashQueryParams": []string{"d", "e", "f"},
							"ConsistentHashRegex": "foo", "DeepCachingType": tc.DeepCachingTypeNever, "FQPacingRate": 41, "MaxOriginConnections": 500,
							"SigningAlgorithm": "uri_signing", "Tenant": "tenant1", "TRRequestHeaders": "X-ooF\nX-raB",
							"TRResponseHeaders": "Access-Control-Max-Age: 600\nContent-Type: text/html; charset=utf-8",
						})),
				},
				"BAD REQUEST when INVALID COUNTRY CODE": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								GeoLimit:          util.Ptr(2),
								GeoLimitCountries: util.Ptr("US,CA,12"),
								XMLID:             util.Ptr("invalid-geolimit-test"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when CHANGING TOPOLOGY of DS with ORG SERVERS ASSIGNED": {
					EndpointId:    GetDeliveryServiceID(t, "ds-top"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceFieldsV30: tc.DeliveryServiceFieldsV30{
								Topology: util.Ptr("another-topology"),
							},
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds-top"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusBadRequest)),
				},
				"BAD REQUEST when UPDATING DS OUTSIDE TENANCY": {
					EndpointId:    GetDeliveryServiceID(t, "ds3"),
					ClientSession: tenant4UserSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds3"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
				"PRECONDITION FAILED when updating with IMS & IUS Headers": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestOpts:   client.RequestOptions{Header: http.Header{rfc.IfUnmodifiedSince: {currentTimeRFC}}},
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds1"),
							},
						},
					}),
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
				"PRECONDITION FAILED when updating with IFMATCH ETAG Header": {
					EndpointId:    GetDeliveryServiceID(t, "ds1"),
					ClientSession: TOSession,
					RequestBody: useDefaultDeliveryService(t, tc.DeliveryServiceV4{
						DeliveryServiceV40: tc.DeliveryServiceV40{
							DeliveryServiceNullableFieldsV11: tc.DeliveryServiceNullableFieldsV11{
								XMLID: util.Ptr("ds1"),
							},
						},
					}),
					RequestOpts:  client.RequestOptions{Header: http.Header{rfc.IfMatch: {rfc.ETag(currentTime)}}},
					Expectations: utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusPreconditionFailed)),
				},
			},
			"DELETE": {
				"BAD REQUEST when DELETING DS OUTSIDE TENANCY": {
					EndpointId:    GetDeliveryServiceID(t, "ds3"),
					ClientSession: tenant4UserSession,
					Expectations:  utils.CkRequest(utils.HasError(), utils.HasStatus(http.StatusForbidden)),
				},
			},
		}

		for method, testCases := range methodTests {
			t.Run(method, func(t *testing.T) {
				for name, testCase := range testCases {
					switch method {
					case "GET":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.GetDeliveryServices(testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "POST":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.CreateDeliveryService(testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "PUT":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.UpdateDeliveryService(testCase.EndpointId(), testCase.RequestBody, testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, resp.Response, resp.Alerts, err)
							}
						})
					case "DELETE":
						t.Run(name, func(t *testing.T) {
							resp, reqInf, err := testCase.ClientSession.DeleteDeliveryService(testCase.EndpointId(), testCase.RequestOpts)
							for _, check := range testCase.Expectations {
								check(t, reqInf, nil, resp.Alerts, err)
							}
						})
					}
				}
			})
		}
	})
}

func validateDSExpectedFields(expectedResp map[string]interface{}) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		dsResp := resp.([]tc.DeliveryServiceV4)
		for field, expected := range expectedResp {
			for _, ds := range dsResp {
				switch field {
				case "Active":
					assert.Equal(t, expected, *ds.Active, "Expected Active to be %v, but got %v", expected, *ds.Active)
				case "DeepCachingType":
					assert.Equal(t, expected, *ds.DeepCachingType, "Expected DeepCachingType to be %v, but got %v", expected, *ds.DeepCachingType)
				case "CDNName":
					assert.Equal(t, expected, *ds.CDNName, "Expected CDNName to be %v, but got %v", expected, *ds.CDNName)
				case "ConsistentHashRegex":
					assert.Equal(t, expected, *ds.ConsistentHashRegex, "Expected ConsistentHashRegex to be %v, but got %v", expected, *ds.ConsistentHashRegex)
				case "ConsistentHashQueryParams":
					assert.Exactly(t, expected, ds.ConsistentHashQueryParams, "Expected ConsistentHashQueryParams to be %v, but got %v", expected, ds.ConsistentHashQueryParams)
				case "DisplayName":
					assert.Equal(t, expected, *ds.DisplayName, "Expected DisplayName to be %v, but got %v", expected, *ds.DisplayName)
				case "DSCP":
					assert.Equal(t, expected, *ds.DSCP, "Expected DSCP to be %v, but got %v", expected, *ds.DSCP)
				case "FQPacingRate":
					assert.Equal(t, expected, *ds.FQPacingRate, "Expected FQPacingRate to be %v, but got %v", expected, *ds.FQPacingRate)
				case "GeoLimit":
					assert.Equal(t, expected, *ds.GeoLimit, "Expected GeoLimit to be %v, but got &v", expected, ds.GeoLimit)
				case "GeoLimitCountries":
					assert.Exactly(t, expected, ds.GeoLimitCountries, "Expected GeoLimitCountries to be %v, but got &v", expected, ds.GeoLimitCountries)
				case "InitialDispersion":
					assert.Equal(t, expected, *ds.InitialDispersion, "Expected InitialDispersion to be %v, but got &v", expected, ds.InitialDispersion)
				case "IPV6RoutingEnabled":
					assert.Equal(t, expected, *ds.IPV6RoutingEnabled, "Expected IPV6RoutingEnabled to be %v, but got &v", expected, ds.IPV6RoutingEnabled)
				case "LogsEnabled":
					assert.Equal(t, expected, *ds.LogsEnabled, "Expected LogsEnabled to be %v, but got %v", expected, *ds.LogsEnabled)
				case "LongDesc":
					assert.Equal(t, expected, *ds.LongDesc, "Expected LongDesc to be %v, but got %v", expected, *ds.LongDesc)
				case "MaxDNSAnswers":
					assert.Equal(t, expected, *ds.MaxDNSAnswers, "Expected MaxDNSAnswers to be %v, but got %v", expected, *ds.MaxDNSAnswers)
				case "MaxOriginConnections":
					assert.Equal(t, expected, *ds.MaxOriginConnections, "Expected MaxOriginConnections to be %v, but got %v", expected, *ds.MaxOriginConnections)
				case "MaxRequestHeaderSize":
					assert.Equal(t, expected, *ds.MaxRequestHeaderBytes, "Expected MaxRequestHeaderBytes to be %v, but got %v", expected, *ds.MaxRequestHeaderBytes)
				case "MissLat":
					assert.Equal(t, expected, *ds.MissLat, "Expected MissLat to be %v, but got %v", expected, *ds.MissLat)
				case "MissLong":
					assert.Equal(t, expected, *ds.MissLong, "Expected MissLong to be %v, but got %v", expected, *ds.MissLong)
				case "MultiSiteOrigin":
					assert.Equal(t, expected, *ds.MultiSiteOrigin, "Expected MultiSiteOrigin to be %v, but got %v", expected, *ds.MultiSiteOrigin)
				case "OrgServerFQDN":
					assert.Equal(t, expected, *ds.OrgServerFQDN, "Expected OrgServerFQDN to be %v, but got %v", expected, *ds.OrgServerFQDN)
				case "ProfileName":
					assert.Equal(t, expected, *ds.ProfileName, "Expected ProfileName to be %v, but got %v", expected, *ds.ProfileName)
				case "Protocol":
					assert.Equal(t, expected, *ds.Protocol, "Expected Protocol to be %v, but got %v", expected, *ds.Protocol)
				case "QStringIgnore":
					assert.Equal(t, expected, *ds.QStringIgnore, "Expected QStringIgnore to be %v, but got %v", expected, *ds.QStringIgnore)
				case "RegionalGeoBlocking":
					assert.Equal(t, expected, *ds.RegionalGeoBlocking, "Expected QStringIgnore to be %v, but got %v", expected, *ds.RegionalGeoBlocking)
				case "ServiceCategory":
					assert.Equal(t, expected, *ds.ServiceCategory, "Expected ServiceCategory to be %v, but got %v", expected, *ds.ServiceCategory)
				case "SigningAlgorithm":
					assert.Equal(t, expected, *ds.SigningAlgorithm, "Expected SigningAlgorithm to be %v, but got %v", expected, *ds.SigningAlgorithm)
				case "Tenant":
					assert.Equal(t, expected, *ds.Tenant, "Expected Tenant to be %v, but got %v", expected, *ds.Tenant)
				case "Topology":
					assert.Equal(t, expected, *ds.Topology, "Expected Topology to be %v, but got %v", expected, *ds.Topology)
				case "TRRequestHeaders":
					assert.Equal(t, expected, *ds.TRRequestHeaders, "Expected TRRequestHeaders to be %v, but got %v", expected, *ds.TRRequestHeaders)
				case "TRResponseHeaders":
					assert.Equal(t, expected, *ds.TRResponseHeaders, "Expected TRResponseHeaders to be %v, but got %v", expected, *ds.TRResponseHeaders)
				case "Type":
					assert.Equal(t, expected, *ds.Type, "Expected Type to be %v, but got %v", expected, *ds.Type)
				case "XMLID":
					assert.Equal(t, expected, *ds.XMLID, "Expected XMLID to be %v, but got %v", expected, *ds.XMLID)
				default:
					t.Errorf("Expected field: %v, does not exist in response", field)
				}
			}
		}
	}
}

func validatePagination(paginationParam string) utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, _ tc.Alerts, _ error) {
		paginationResp := resp.([]tc.DeliveryServiceV4)

		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("orderby", "id")
		respBase, _, err := TOSession.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Cannot get Delivery Services: %v - alerts: %+v", err, respBase.Alerts)

		ds := respBase.Response
		assert.RequireGreaterOrEqual(t, len(ds), 3, "Need at least 3 Delivery Services in Traffic Ops to test pagination support, found: %d", len(ds))
		switch paginationParam {
		case "limit:":
			assert.Exactly(t, ds[:1], paginationResp, "expected GET Delivery Services with limit = 1 to return first result")
		case "offset":
			assert.Exactly(t, ds[1:2], paginationResp, "expected GET Delivery Services with limit = 1, offset = 1 to return second result")
		case "page":
			assert.Exactly(t, ds[1:2], paginationResp, "expected GET Delivery Services with limit = 1, page = 2 to return second result")
		}
	}
}

func validateDescSort() utils.CkReqFunc {
	return func(t *testing.T, _ toclientlib.ReqInf, resp interface{}, alerts tc.Alerts, _ error) {
		dsDescResp := resp.([]tc.DeliveryServiceV4)
		var descSortedList []string
		var ascSortedList []string
		assert.GreaterOrEqual(t, len(dsDescResp), 2, "Need at least 2 XMLIDs in Traffic Ops to test desc sort, found: %d", len(dsDescResp))
		// Get delivery services in the default ascending order for comparison.
		dsAscResp, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
		assert.NoError(t, err, "Unexpected error getting Delivery Services with default sort order: %v - alerts: %+v", err, dsAscResp.Alerts)
		assert.GreaterOrEqual(t, len(dsAscResp.Response), 2, "Need at least 2 XMLIDs in Traffic Ops to test sort, found %d", len(dsAscResp.Response))
		// Verify the response match in length, i.e. equal amount of delivery services.
		assert.Equal(t, len(dsAscResp.Response), len(dsDescResp), "Expected descending order response length: %v, to match ascending order response length %v", len(dsAscResp.Response), len(dsDescResp))
		// Insert xmlIDs to the front of a new list, so they are now reversed to be in ascending order.
		for _, ds := range dsDescResp {
			descSortedList = append([]string{*ds.XMLID}, descSortedList...)
		}
		// Insert xmlIDs by appending to a new list, so they stay in ascending order.
		for _, ds := range dsAscResp.Response {
			ascSortedList = append(ascSortedList, *ds.XMLID)
		}
		assert.Exactly(t, ascSortedList, descSortedList, "Delivery Service responses are not equal after reversal: %v - %v", ascSortedList, descSortedList)
	}
}

func GetDeliveryServiceID(t *testing.T, xmlId string) func() int {
	return func() int {
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("xmlId", xmlId)

		resp, _, err := TOSession.GetDeliveryServices(opts)
		assert.RequireNoError(t, err, "Get Delivery Service Request failed with error: %v", err)
		assert.RequireEqual(t, 1, len(resp.Response), "Expected delivery service response object length 1, but got %d", len(resp.Response))
		assert.RequireNotNil(t, resp.Response[0].ID, "Expected id to not be nil")

		return *resp.Response[0].ID
	}
}

func useDefaultDeliveryService(t *testing.T, newDSWithFields tc.DeliveryServiceV4) tc.DeliveryServiceV4 {
	// A Generated HTTP Delivery Service
	if newDSWithFields.Active == nil {
		newDSWithFields.Active = util.Ptr(true)
	}
	if newDSWithFields.CDNID == nil {
		newDSWithFields.CDNID = util.Ptr(GetCDNID(t, "cdn1")())
	}
	if newDSWithFields.DisplayName == nil {
		newDSWithFields.DisplayName = util.Ptr("generated test ds")
	}
	if newDSWithFields.DSCP == nil {
		newDSWithFields.DSCP = util.Ptr(0)
	}
	if newDSWithFields.GeoLimit == nil {
		newDSWithFields.GeoLimit = util.Ptr(0)
	}
	if newDSWithFields.GeoProvider == nil {
		newDSWithFields.GeoProvider = util.Ptr(0)
	}
	if newDSWithFields.InitialDispersion == nil {
		newDSWithFields.InitialDispersion = util.Ptr(1)
	}
	if newDSWithFields.IPV6RoutingEnabled == nil {
		newDSWithFields.IPV6RoutingEnabled = util.Ptr(false)
	}
	if newDSWithFields.LogsEnabled == nil {
		newDSWithFields.LogsEnabled = util.Ptr(false)
	}
	if newDSWithFields.MissLat == nil {
		newDSWithFields.MissLat = util.Ptr(0.0)
	}
	if newDSWithFields.MissLong == nil {
		newDSWithFields.MissLong = util.Ptr(0.0)
	}
	if newDSWithFields.MultiSiteOrigin == nil {
		newDSWithFields.MultiSiteOrigin = util.Ptr(false)
	}
	if newDSWithFields.OrgServerFQDN == nil {
		newDSWithFields.OrgServerFQDN = util.Ptr("http://generated.ds.test")
	}
	if newDSWithFields.Protocol == nil {
		newDSWithFields.Protocol = util.Ptr(0)
	}
	if newDSWithFields.ProfileName == nil {
		newDSWithFields.ProfileName = util.Ptr("ATS_EDGE_TIER_CACHE")
	}
	if newDSWithFields.QStringIgnore == nil {
		newDSWithFields.QStringIgnore = util.Ptr(0)
	}
	if newDSWithFields.RangeRequestHandling == nil {
		newDSWithFields.RangeRequestHandling = util.Ptr(0)
	}
	if newDSWithFields.RegionalGeoBlocking == nil {
		newDSWithFields.RegionalGeoBlocking = util.Ptr(false)
	}
	if newDSWithFields.RoutingName == nil {
		newDSWithFields.RoutingName = util.Ptr("ccr-ds-test")
	}
	if newDSWithFields.TenantID == nil {
		newDSWithFields.TenantID = util.Ptr(GetTenantID(t, "tenant1")())
	}
	if newDSWithFields.TypeID == nil {
		newDSWithFields.TypeID = util.Ptr(GetTypeID(t, "HTTP")())
	}
	if newDSWithFields.XMLID == nil {
		newDSWithFields.XMLID = util.Ptr("generated-ds")
	}
	return newDSWithFields
}

func CreateTestDeliveryServices(t *testing.T) {
	for _, ds := range testData.DeliveryServices {
		ds = ds.RemoveLD1AndLD2()
		if ds.XMLID == nil {
			t.Error("Found a Delivery Service in testing data with null or undefined XMLID")
			continue
		}
		resp, _, err := TOSession.CreateDeliveryService(ds, client.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service '%s': %v - alerts: %+v", *ds.XMLID, err, resp.Alerts)
	}
}

func DeleteTestDeliveryServices(t *testing.T) {
	dses, _, err := TOSession.GetDeliveryServices(client.RequestOptions{})
	assert.NoError(t, err, "Cannot get Delivery Services: %v - alerts: %+v", err, dses.Alerts)

	for _, ds := range dses.Response {
		delResp, _, err := TOSession.DeleteDeliveryService(*ds.ID, client.RequestOptions{})
		assert.NoError(t, err, "Could not delete Delivery Service: %v - alerts: %+v", err, delResp.Alerts)
		// Retrieve Delivery Service to see if it got deleted
		opts := client.NewRequestOptions()
		opts.QueryParameters.Set("id", strconv.Itoa(*ds.ID))
		getDS, _, err := TOSession.GetDeliveryServices(opts)
		assert.NoError(t, err, "Error deleting Delivery Service for '%s' : %v - alerts: %+v", *ds.XMLID, err, getDS.Alerts)
		assert.Equal(t, 0, len(getDS.Response), "Expected Delivery Service '%s' to be deleted", *ds.XMLID)
	}
}
