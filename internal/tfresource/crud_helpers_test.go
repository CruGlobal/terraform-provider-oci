// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package tfresource

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oci_common "github.com/oracle/oci-go-sdk/v60/common"
	oci_vault "github.com/oracle/oci-go-sdk/v60/vault"
	oci_work_requests "github.com/oracle/oci-go-sdk/v60/workrequests"
)

type extraWaitPostCreateDelete interface {
	ExtraWaitPostCreateDelete() time.Duration
}

type statefullyCreatedResource interface {
	State() string
	CreatedPending() []string
	CreatedTarget() []string
	setState(StatefulResource) error
	Get() error
}

type SyncResource interface {
	GetMutex() *sync.Mutex
}
type ResourceCrud struct {
	id                     string
	D                      *mockResourceData
	Mutex                  *sync.Mutex
	DisableNotFoundRetries bool
	SyncResource
	statefullyCreatedResource
	extraWaitPostCreateDelete
}

func (b ResourceCrud) Create() error {
	if b.id == "4" {
		return errors.New("")
	}
	return nil
}
func (b ResourceCrud) ID() string {
	return ""
}
func (b ResourceCrud) SetData() error {
	return nil
}
func (b ResourceCrud) VoidState() {}
func (b ResourceCrud) GetMutex() *sync.Mutex {
	return &sync.Mutex{}
}
func (b ResourceCrud) State() string {
	if b.id == "1" {
		return "SUCCEEDED"
	}
	return "FAILED"
}
func (b ResourceCrud) Get() error {
	if b.id == "2" {
		return errors.New("")
	}
	return nil
}
func (b ResourceCrud) CreatedPending() []string {
	return []string{"FAILED", "b"}
}
func (b ResourceCrud) CreatedTarget() []string {
	return []string{"FAILED", "b"}
}
func (b ResourceCrud) setState(s StatefulResource) error {
	if b.id == "3" {
		return errors.New("")
	}
	return nil
}
func (b ResourceCrud) ExtraWaitPostCreateDelete() time.Duration {
	timeoutDuration, _ := time.ParseDuration("1s")
	return timeoutDuration
}

type statefullyDeletedResource interface {
	State() string
	DeletedPending() []string
	DeletedTarget() []string
	setState(StatefulResource) error
	Get() error
}
type readResourceCrud struct {
	D                      *mockResourceData
	Mutex                  *sync.Mutex
	Client                 *oci_vault.VaultsClient
	Res                    *oci_vault.Secret
	DisableNotFoundRetries bool
	SyncResource
	statefullyDeletedResource
}

func (b readResourceCrud) Create() error {
	return nil
}
func (b readResourceCrud) ID() string {
	return ""
}
func (b readResourceCrud) SetData() error {
	return nil
}
func (b readResourceCrud) VoidState() {}
func (b readResourceCrud) State() string {
	return "FAILED"
}
func (b readResourceCrud) Get() error {
	return nil
}
func (b readResourceCrud) DeletedPending() []string {
	return []string{"a", "b"}
}
func (b readResourceCrud) DeletedTarget() []string {
	return []string{"a", "b"}
}

type statefullyUpdatedResource interface {
	State() string
	DeletedPending() []string
	DeletedTarget() []string
	setState(StatefulResource) error
	Get() error
}
type updateResourceCrud struct {
	D                      *mockResourceData
	Mutex                  *sync.Mutex
	Client                 *oci_vault.VaultsClient
	Res                    *oci_vault.Secret
	DisableNotFoundRetries bool
	SyncResource
	statefullyUpdatedResource
}

func (b updateResourceCrud) Update() error {
	return nil
}
func (b updateResourceCrud) ID() string {
	return ""
}
func (b updateResourceCrud) SetData() error {
	return nil
}
func (b updateResourceCrud) VoidState() {}
func (b updateResourceCrud) State() string {
	return "FAILED"
}
func (b updateResourceCrud) Get() error {
	return nil
}
func (b updateResourceCrud) GetMutex() *sync.Mutex {
	return &sync.Mutex{}
}
func (b updateResourceCrud) UpdatedPending() []string {
	return []string{"a", "b"}
}
func (b updateResourceCrud) UpdatedTarget() []string {
	return []string{"a", "b"}
}

type extraWaitPostDelete interface {
	ExtraWaitPostDelete() time.Duration
}
type deleteResourceCrud struct {
	D                      *mockResourceData
	Mutex                  *sync.Mutex
	DisableNotFoundRetries bool
	SyncResource
	statefullyDeletedResource
	extraWaitPostCreateDelete
	extraWaitPostDelete
}

func (b deleteResourceCrud) Delete() error {
	return nil
}
func (b deleteResourceCrud) ID() string {
	return ""
}
func (b deleteResourceCrud) SetData() error {
	return nil
}
func (b deleteResourceCrud) VoidState() {}
func (b deleteResourceCrud) State() string {
	return "FAILED"
}
func (b deleteResourceCrud) Get() error {
	return nil
}
func (b deleteResourceCrud) GetMutex() *sync.Mutex {
	return &sync.Mutex{}
}
func (b deleteResourceCrud) DeletedPending() []string {
	return []string{"a", "b"}
}
func (b deleteResourceCrud) DeletedTarget() []string {
	return []string{"a", "b"}
}
func (b deleteResourceCrud) ExtraWaitPostCreateDelete() time.Duration {
	timeoutDuration, _ := time.ParseDuration("1s")
	return timeoutDuration
}
func (b deleteResourceCrud) ExtraWaitPostDelete() time.Duration {
	timeoutDuration, _ := time.ParseDuration("1s")
	return timeoutDuration
}

type mockResourceData struct {
	state string
}

func (d *mockResourceData) GetOkExists(_ string) (interface{}, bool) {
	return "test", true
}
func (d *mockResourceData) SetId(_ string) {
}
func (d *mockResourceData) Timeout(_ string) time.Duration {
	timeoutDuration, _ := time.ParseDuration("10m")
	return timeoutDuration
}
func (d *mockResourceData) Partial(_ bool) {}
func (d *mockResourceData) HasChange(_ string) bool {
	if d.state == "1" {
		return false
	}
	return true
}
func (d *mockResourceData) GetChange(_ string) (interface{}, interface{}) {
	if d.state == "2" {
		return []interface{}{map[string]interface{}{"load_balancer_id": "1"}, map[string]interface{}{"load_balancer_id": "2"}}, []interface{}{map[string]interface{}{"load_balancer_id": "1"}}
	}
	if d.state == "3" {
		return []interface{}{map[string]interface{}{"load_balancer_id": "1"}}, []interface{}{map[string]interface{}{"load_balancer_id": "2"}}
	}
	return []interface{}{map[string]interface{}{"load_balancer_id": "1"}}, []interface{}{map[string]interface{}{"load_balancer_id": "1"}}
}

type mockWorkRequestClient struct{}

func (client *mockWorkRequestClient) GetWorkRequest(_ context.Context, wreq oci_work_requests.GetWorkRequestRequest) (response oci_work_requests.GetWorkRequestResponse, err error) {
	var v oci_work_requests.GetWorkRequestResponse
	switch *wreq.WorkRequestId {
	case "1":
		s := ""
		et := "default"
		id := "oci"
		r := []oci_work_requests.WorkRequestResource{{EntityType: &et, Identifier: &id, ActionType: "CREATED"}}
		wr := oci_work_requests.WorkRequest{Status: "SUCCEEDED", Resources: r}
		v = oci_work_requests.GetWorkRequestResponse{RawResponse: nil, WorkRequest: wr, OpcRequestId: &s}
	case "2":
		s := ""
		et := "default"
		id := ""
		r := []oci_work_requests.WorkRequestResource{{EntityType: &et, Identifier: &id, ActionType: "CREATED"}}
		wr := oci_work_requests.WorkRequest{Status: "ABC", Resources: r}
		v = oci_work_requests.GetWorkRequestResponse{RawResponse: nil, WorkRequest: wr, OpcRequestId: &s}
	case "3":
		s := ""
		et := "default"
		id := "oci"
		r := []oci_work_requests.WorkRequestResource{{EntityType: &et, Identifier: &id, ActionType: "CREATED"}}
		wr := oci_work_requests.WorkRequest{Status: "CANCELED", Resources: r}
		v = oci_work_requests.GetWorkRequestResponse{RawResponse: nil, WorkRequest: wr, OpcRequestId: &s}
	default:
		s := ""
		et := "default"
		id := "oci"
		r := []oci_work_requests.WorkRequestResource{{EntityType: &et, Identifier: &id, ActionType: "CREATED"}}
		wr := oci_work_requests.WorkRequest{Status: "SUCCEEDED", Resources: r}
		v = oci_work_requests.GetWorkRequestResponse{RawResponse: nil, WorkRequest: wr, OpcRequestId: &s}
	}
	return v, nil
}
func (client *mockWorkRequestClient) ListWorkRequestErrors(_ context.Context, _ oci_work_requests.ListWorkRequestErrorsRequest) (response oci_work_requests.ListWorkRequestErrorsResponse, err error) {
	s := "default"
	v := oci_work_requests.ListWorkRequestErrorsResponse{RawResponse: nil, Items: nil, OpcNextPage: &s, OpcRequestId: &s}
	return v, nil
}

type TestResource struct {
	GetError          error
	GetAttempts       int
	ActualGetAttempts int
}

func (t *TestResource) Get() error {
	t.ActualGetAttempts++
	t.GetAttempts--
	if t.GetAttempts <= 0 {
		return t.GetError
	}
	return nil
}

func TestUnitNormalizeBoolString(t *testing.T) {
	type args struct {
		value string
	}
	type expected struct {
		res      string
		gotError bool
	}
	type testFormat struct {
		name     string
		args     args
		expected expected
	}
	tests := []testFormat{
		{
			name:     "Test valid bool string '1'",
			args:     args{value: "1"},
			expected: expected{res: "true", gotError: false},
		},
		{
			name:     "Test valid bool string '0'",
			args:     args{value: "0"},
			expected: expected{res: "false", gotError: false},
		},
		{
			name:     "Test valid bool string 'false'",
			args:     args{value: "false"},
			expected: expected{res: "false", gotError: false},
		},
		{
			name:     "Test valid bool string 'True'",
			args:     args{value: "True"},
			expected: expected{res: "true", gotError: false},
		},
		{
			name:     "Test invalid bool string 'hello'",
			args:     args{value: "hello"},
			expected: expected{res: "", gotError: true},
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		res, err := NormalizeBoolString(test.args.value)
		if res != test.expected.res {
			t.Errorf("Output %q not equal to expected %q", res, test.expected.res)
		}
		if (err != nil) != test.expected.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", err, test.expected.gotError)
		}
	}
}

func TestUnitGetDataSourceItemSchema(t *testing.T) {
	reqSchema := &schema.Resource{
		Schema: nil,
		Create: func(d *schema.ResourceData, m interface{}) error { return nil },
		Read:   func(d *schema.ResourceData, m interface{}) error { return nil },
	}
	resSchema := &schema.Resource{
		Schema: nil,
		Create: nil,
		Read:   nil,
	}

	changeReqSchema := func(s *schema.Resource, v bool) *schema.Resource {
		if v {
			s.Schema = make(map[string]*schema.Schema)
			s.Schema["Id"] = &schema.Schema{
				Computed: false,
				Type:     schema.TypeInt,
			}
		}
		return s
	}
	changeResSchema := func(s *schema.Resource, v bool) *schema.Resource {
		if v {
			s.Schema = make(map[string]*schema.Schema)
			s.Schema["Id"] = &schema.Schema{
				Computed: true,
				Type:     schema.TypeString,
			}
		}
		return s
	}
	type args struct {
		resourceSchema *schema.Resource
	}
	type testFormat struct {
		name   string
		args   args
		result *schema.Resource
	}
	tests := []testFormat{
		{
			name:   "Test schema contains 'Id' field",
			args:   args{resourceSchema: changeReqSchema(reqSchema, true)},
			result: changeResSchema(resSchema, true),
		},
		{
			name:   "Test schema does not contain 'Id' field",
			args:   args{resourceSchema: changeReqSchema(reqSchema, false)},
			result: changeResSchema(resSchema, false),
		},
	}
	convertResFieldsToDSFields = func(resSchema *schema.Resource) *schema.Resource {
		return resSchema
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		t.Logf("%s", fmt.Sprint(test.args.resourceSchema))
		if res := GetDataSourceItemSchema(test.args.resourceSchema); reflect.DeepEqual(res, test.result) {
			t.Logf("output schema %q", fmt.Sprint(res))
			t.Logf("expected schema %q", fmt.Sprint(test.result))
			t.Errorf("Output schema and expected schema are not same")
		}
	}

}

func TestUnitGetSingularDataSourceItemSchema(t *testing.T) {
	reqSchema := &schema.Resource{
		Schema: nil,
		Create: func(d *schema.ResourceData, m interface{}) error { return nil },
		Update: func(d *schema.ResourceData, m interface{}) error { return nil },
		Delete: func(d *schema.ResourceData, m interface{}) error { return nil },
		Read:   nil,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: nil,
			Read:   nil,
		},
		CustomizeDiff: func(d *schema.ResourceDiff, m interface{}) error { return nil },
	}
	reqReadFunc := func(d *schema.ResourceData, m interface{}) error { return nil }
	reqAddFieldMap := map[string]*schema.Schema{
		"foo": {
			Computed: true,
		},
	}
	resSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"foo": {
				Computed: true,
			},
		},
		Create:        nil,
		Read:          reqReadFunc,
		Update:        nil,
		Delete:        nil,
		Importer:      nil,
		Timeouts:      nil,
		CustomizeDiff: nil,
	}
	convertResFieldsToDSFields = func(resSchema *schema.Resource) *schema.Resource {
		return resSchema
	}
	type args struct {
		resourceSchema *schema.Resource
		addFieldMap    map[string]*schema.Schema
		readFunc       schema.ReadFunc
	}
	type testFormat struct {
		name   string
		args   args
		result *schema.Resource
	}
	changeReqSchema := func(s *schema.Resource, v bool) *schema.Resource {
		if v {
			s.Schema = make(map[string]*schema.Schema)
			s.Schema["Id"] = &schema.Schema{
				Computed: false,
				Type:     schema.TypeInt,
			}
		}
		return s
	}
	changeResSchema := func(s *schema.Resource, v bool) *schema.Resource {
		if v {
			s.Schema["Id"] = &schema.Schema{
				Computed: true,
				Type:     schema.TypeString,
			}
		}
		return s
	}

	tests := []testFormat{
		{
			name:   "Test schema contains 'Id' field",
			args:   args{resourceSchema: changeReqSchema(reqSchema, true), addFieldMap: reqAddFieldMap, readFunc: reqReadFunc},
			result: changeResSchema(resSchema, true),
		},
		{
			name:   "Test schema does not contain 'Id' field",
			args:   args{resourceSchema: changeReqSchema(reqSchema, false), addFieldMap: reqAddFieldMap, readFunc: reqReadFunc},
			result: changeResSchema(resSchema, false),
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := GetSingularDataSourceItemSchema(test.args.resourceSchema, test.args.addFieldMap, test.args.readFunc); reflect.DeepEqual(res, test.result) {
			t.Logf("output schema %q", fmt.Sprint(res))
			t.Logf("expected schema %q", fmt.Sprint(test.result))
			t.Errorf("Output schema and expected schema are not same")
		}
	}
}

func TestUnitconvertResourceFieldsToDatasourceFields(t *testing.T) {
	changeReqSchema := func(k string) *schema.Resource {
		reqSchema := &schema.Resource{
			Schema: map[string]*schema.Schema{
				"oci": {
					Computed: false,
					Required: true,
					Optional: true,
				},
			},
		}
		if k == "1" {
			reqSchema.Schema["oci"].Type = schema.TypeSet
			reqSchema.Schema["oci"].Elem = nil
		}
		return reqSchema
	}
	type args struct {
		resourceSchema *schema.Resource
	}
	type testFormat struct {
		name   string
		args   args
		result *schema.Resource
	}
	changeResSchema := func(k string) *schema.Resource {
		resSchema := &schema.Resource{
			Schema: map[string]*schema.Schema{
				"oci": {
					Computed:         true,
					Required:         false,
					Optional:         false,
					DiffSuppressFunc: nil,
					ValidateFunc:     nil,
					ConflictsWith:    nil,
					Default:          nil,
					DefaultFunc:      nil,
				},
			},
		}
		if k == "1" {
			resSchema.Schema["oci"].Type = schema.TypeList
			resSchema.Schema["oci"].Set = nil
		}
		return resSchema
	}

	tests := []testFormat{
		{
			name:   "Test schema contains TypeSet as Type and nil as Set",
			args:   args{resourceSchema: changeReqSchema("1")},
			result: changeResSchema("1"),
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		convertResourceFieldsToDatasourceFields(test.args.resourceSchema)
	}
}

func TestUnitGetTimeoutDuration(t *testing.T) {
	type args struct {
		timeout string
	}
	type testFormat struct {
		name   string
		args   args
		output *time.Duration
	}
	timeout := func(timeout string) *time.Duration {
		timeoutDuration, _ := time.ParseDuration(timeout)
		return &timeoutDuration
	}
	tests := []testFormat{
		{
			name:   "Test valid timeout string '20m'",
			args:   args{timeout: "20m"},
			output: timeout("20m"),
		},
		{
			name:   "Test invalid timeout string 'hello'",
			args:   args{timeout: "hello"},
			output: &TwentyMinutes,
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := GetTimeoutDuration(test.args.timeout); *res != *test.output {
			t.Errorf("Output %q not equal to expected %q", *res, *test.output)
		}
	}
}

func TestUnitGenericMapToJsonMap(t *testing.T) {
	type args struct {
		genericMap map[string]interface{}
	}
	type testFormat struct {
		name     string
		args     args
		output   map[string]interface{}
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:   "Test bool value for genericMap bool",
			args:   args{genericMap: map[string]interface{}{"a": true}},
			output: map[string]interface{}{"a": true},
			mockFunc: func() {
				jsonMarshal = func(object interface{}) ([]byte, error) {
					return []byte{116, 114, 117, 101}, nil
				}
			},
		},
		{
			name:   "Test string array value for genericMap string",
			args:   args{genericMap: map[string]interface{}{"a": []string{"foo", "bar"}}},
			output: map[string]interface{}{"a": []string{"foo", "bar"}},
			mockFunc: func() {
				jsonMarshal = func(object interface{}) ([]byte, error) {
					return []byte{91, 34, 102, 111, 111, 34, 44, 34, 98, 97, 114, 34, 93}, nil
				}
			},
		},
		//{
		//	name: "Test string value for genericMap string",
		//	args: args{genericMap: map[string]interface{}{"a": "a"}},
		//	output: map[string]interface{}{"a": "a"},
		//	mockFunc: func() {
		//		jsonMarshal = func(object interface{}) ([]byte, error) {
		//			return []byte{34, 97, 34}, nil
		//		}
		//	},
		//},
		//{
		//	name: "Test invalid value for genericMap string",
		//	args: args{genericMap: map[string]interface{}{"a": true}},
		//	output: map[string]interface{}{},
		//	mockFunc: func() {
		//		jsonMarshal = func(object interface{}) ([]byte, error) {
		//			return []byte{34, 97, 34}, errors.New("")
		//		}
		//	},
		//},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		test.mockFunc()
		if res := GenericMapToJsonMap(test.args.genericMap); reflect.DeepEqual(res, test.output) {
			t.Logf("output map %q", fmt.Sprint(res))
			t.Logf("expected map %q", fmt.Sprint(test.output))
			t.Errorf("Output map and expected map are not same")
		}
	}
}

func TestUnitDbVersionDiffSuppress(t *testing.T) {
	reqResourceData := &schema.ResourceData{}
	type args struct {
		key string
		old string
		new string
		d   *schema.ResourceData
	}
	type testFormat struct {
		name   string
		args   args
		output bool
	}
	tests := []testFormat{
		{
			name:   "Test valid empty values for old and new",
			args:   args{key: "", old: "", new: "", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid same values for old and new ",
			args:   args{key: "", old: "18.0.0.0", new: "18.0.0.0", d: reqResourceData},
			output: true,
		},
		{
			name:   "Test valid different values for old and new ",
			args:   args{key: "", old: "16.0.0.0", new: "19.0.0.0", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid prefix values",
			args:   args{key: "", old: "16.1.0.0", new: "16.1", d: reqResourceData},
			output: true,
		},
		{
			name:   "Test valid different values for old and new ",
			args:   args{key: "", old: "15.0.0.0", new: "15.1.0.0", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid different values for old and new ",
			args:   args{key: "", old: "test", new: "15.1.0.0", d: reqResourceData},
			output: false,
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := DbVersionDiffSuppress(test.args.key, test.args.old, test.args.new, test.args.d); res != test.output {
			t.Errorf("Scenario failed %s", test.name)
			t.Errorf("Output %t not equal to expected %t", res, test.output)
		}
	}
}

func TestUnitAdDiffSuppress(t *testing.T) {
	reqResourceData := &schema.ResourceData{}
	type args struct {
		key string
		old string
		new string
		d   *schema.ResourceData
	}
	type testFormat struct {
		name   string
		args   args
		output bool
	}
	tests := []testFormat{
		{
			name:   "Test valid values for old and new",
			args:   args{key: "", old: "20.10", new: "22.10", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test invalid value for old",
			args:   args{key: "", old: "test", new: "10", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test invalid value for new",
			args:   args{key: "", old: "16", new: "test", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid values for old and new as 'inf' and 'nan'",
			args:   args{key: "", old: "inf", new: "nan", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid same values for old and new ",
			args:   args{key: "", old: "21.1212121", new: "21.12121211", d: reqResourceData},
			output: true,
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := AdDiffSuppress(test.args.key, test.args.old, test.args.new, test.args.d); res != test.output {
			t.Errorf("Output %t not equal to expected %t", res, test.output)
		}
	}
}

func TestUnitGiVersionDiffSuppress(t *testing.T) {
	reqResourceData := &schema.ResourceData{}
	type args struct {
		key string
		old string
		new string
		d   *schema.ResourceData
	}
	type testFormat struct {
		name   string
		args   args
		output bool
	}
	tests := []testFormat{
		{
			name:   "Test valid empty values for old and new",
			args:   args{key: "", old: "", new: "", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid different values for old and new ",
			args:   args{key: "", old: "18.0.0.0", new: "17.0.0.0", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid same values for old and new ",
			args:   args{key: "", old: "18.0.0.0", new: "18.0.0.0", d: reqResourceData},
			output: true,
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := GiVersionDiffSuppress(test.args.key, test.args.old, test.args.new, test.args.d); res != test.output {
			t.Errorf("Output %t not equal to expected %t", res, test.output)
		}
	}
}

func TestUnitMySqlVersionDiffSuppress(t *testing.T) {
	reqResourceData := &schema.ResourceData{}
	type args struct {
		key string
		old string
		new string
		d   *schema.ResourceData
	}
	type testFormat struct {
		name   string
		args   args
		output bool
	}
	tests := []testFormat{
		{
			name:   "Test valid empty values for old and new",
			args:   args{key: "", old: "", new: "", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid different values for old and new ",
			args:   args{key: "", old: "18.0.0.0", new: "17.0.0.0", d: reqResourceData},
			output: false,
		},
		{
			name:   "Test valid same values for old and new ",
			args:   args{key: "", old: "18.1.0.0", new: "18.1.0.0", d: reqResourceData},
			output: true,
		},
		{
			name:   "Test valid different values for old and new ",
			args:   args{key: "", old: "15.0.0.0", new: "15.1.0.0", d: reqResourceData},
			output: false,
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := MySqlVersionDiffSuppress(test.args.key, test.args.old, test.args.new, test.args.d); res != test.output {
			t.Errorf("Output %t not equal to expected %t", res, test.output)
		}
	}
}

func TestUnitloadBalancersSuppressDiff(t *testing.T) {

	type args struct {
		d *mockResourceData
	}
	type testFormat struct {
		name   string
		args   args
		output bool
	}
	changeReqResourceData := func(k string) *mockResourceData {
		reqResourceData := &mockResourceData{
			state: k,
		}
		return reqResourceData
	}
	tests := []testFormat{
		{
			name:   "Test HasChange() is false",
			args:   args{d: changeReqResourceData("1")},
			output: true,
		},
		{
			name:   "Test HasChange() is true and len of oldBalances not equal to new",
			args:   args{d: changeReqResourceData("2")},
			output: false,
		},
		{
			name:   "Test HasChange() is true and value of load_balancer_id differs",
			args:   args{d: changeReqResourceData("3")},
			output: false,
		},
		{
			name:   "Test HasChange() is true and output is true",
			args:   args{d: changeReqResourceData("4")},
			output: true,
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		t.Logf("Running %s", fmt.Sprint(test.args.d))
		if res := loadBalancersSuppressDiff(test.args.d); res != test.output {
			t.Errorf("Output %t not equal to expected %t", res, test.output)
		}
	}
}

func TestUnitGenerateDataSourceHashID(t *testing.T) {
	reqSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"foo": {
				Computed: false,
				Optional: true,
				Type:     schema.TypeList,
			},
			"id": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"bar": {
				Computed: false,
				Required: true,
				Type:     schema.TypeInt,
			},
		},
		Create: nil,
		Read:   nil,
	}
	reqResourceData := &mockResourceData{}
	type args struct {
		idPrefix       string
		resourceSchema *schema.Resource
		resourceData   *mockResourceData
	}
	type testFormat struct {
		name   string
		args   args
		output string
	}
	tests := []testFormat{
		{
			name:   "Test valid empty value for resourceSchema",
			args:   args{idPrefix: "", resourceSchema: nil, resourceData: reqResourceData},
			output: "",
		},
		{
			name:   "Test valid values for resourceSchema and resourceData",
			args:   args{idPrefix: "", resourceSchema: reqSchema, resourceData: reqResourceData},
			output: "2662564013",
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := GenerateDataSourceHashID(test.args.idPrefix, test.args.resourceSchema, test.args.resourceData); res != test.output {
			t.Errorf("Output %s not equal to expected %s", res, test.output)
		}
	}
}

func TestUnitCreateResource(t *testing.T) {
	s := &ResourceCrud{}
	reqResourceData := &mockResourceData{}
	s.D = reqResourceData

	type args struct {
		d    *mockResourceData
		sync ResourceCreator
	}
	type testFormat struct {
		name     string
		args     args
		gotError bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:     "Test error is returned",
			args:     args{sync: s, d: reqResourceData},
			gotError: true,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return errors.New("default")
					//return nil
				}
			},
		},
		{
			name:     "Test no error is returned",
			args:     args{sync: s, d: reqResourceData},
			gotError: false,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					//return errors.New("default")
					return nil
				}
			},
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		test.mockFunc()
		if res := CreateResource(test.args.d, test.args.sync); (res != nil) != test.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.gotError)
		}
	}
}

func TestUnitReadResource(t *testing.T) {
	s := &readResourceCrud{}
	reqResourceData := &mockResourceData{}
	s.D = reqResourceData

	type args struct {
		sync ResourceReader
	}
	type testFormat struct {
		name     string
		args     args
		gotError bool
	}
	tests := []testFormat{
		{
			name:     "Test",
			args:     args{sync: s},
			gotError: false,
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := ReadResource(test.args.sync); (res != nil) != test.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.gotError)
		}
	}
}

func TestUnitUpdateResource(t *testing.T) {
	s := &updateResourceCrud{}
	reqResourceData := &mockResourceData{}
	s.D = reqResourceData

	type args struct {
		d    *mockResourceData
		sync ResourceUpdater
	}
	type testFormat struct {
		name     string
		args     args
		gotError bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:     "Test",
			args:     args{sync: s, d: reqResourceData},
			gotError: true,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return errors.New("default")
				}
			},
		},
		{
			name:     "Test",
			args:     args{sync: s, d: reqResourceData},
			gotError: false,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return nil
				}
			},
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		test.mockFunc()
		if res := UpdateResource(test.args.d, test.args.sync); (res != nil) != test.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.gotError)
		}
	}
}

func TestUnitDeleteResource(t *testing.T) {
	s := &deleteResourceCrud{}
	reqResourceData := &mockResourceData{}
	s.D = reqResourceData

	type args struct {
		d    *mockResourceData
		sync ResourceDeleter
	}
	type testFormat struct {
		name     string
		args     args
		gotError bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:     "Test",
			args:     args{sync: s, d: reqResourceData},
			gotError: true,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return errors.New("default")
				}
			},
		},
		{
			name:     "Test",
			args:     args{sync: s, d: reqResourceData},
			gotError: false,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return nil
				}
			},
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		test.mockFunc()
		if res := DeleteResource(test.args.d, test.args.sync); (res != nil) != test.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.gotError)
		}
	}
}

func TestUnitWaitForUpdatedState(t *testing.T) {
	s := &updateResourceCrud{}
	reqResourceData := &mockResourceData{}
	s.D = reqResourceData

	type args struct {
		d    *mockResourceData
		sync ResourceUpdater
	}
	type testFormat struct {
		name     string
		args     args
		gotError bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:     "Test error is returned",
			args:     args{sync: s, d: reqResourceData},
			gotError: true,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return errors.New("default")
				}
			},
		},
		{
			name:     "Test no error is returned",
			args:     args{sync: s, d: reqResourceData},
			gotError: false,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return nil
				}
			},
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		test.mockFunc()
		if res := WaitForUpdatedState(test.args.d, test.args.sync); (res != nil) != test.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.gotError)
		}
	}
}

func TestUnitWaitForCreatedState(t *testing.T) {
	s := &ResourceCrud{}
	reqResourceData := &mockResourceData{}
	s.D = reqResourceData

	type args struct {
		d    *mockResourceData
		sync ResourceCreator
	}
	type testFormat struct {
		name     string
		args     args
		gotError bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:     "Test error is returned",
			args:     args{sync: s, d: reqResourceData},
			gotError: true,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return errors.New("default")
				}
			},
		},
		{
			name:     "Test no error is returned",
			args:     args{sync: s, d: reqResourceData},
			gotError: false,
			mockFunc: func() {
				waitForStateRefreshVar = func(sr StatefulResource, timeout time.Duration, operationName string, pending []string, target []string) error {
					return nil
				}
			},
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		test.mockFunc()
		if res := WaitForCreatedState(test.args.d, test.args.sync); (res != nil) != test.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.gotError)
		}
	}
}

func TestUnitWaitForStateRefresh(t *testing.T) {
	type args struct {
		sync          StatefulResource
		timeout       time.Duration
		operationName string
		pending       []string
		target        []string
	}
	type testFormat struct {
		name     string
		args     args
		gotError bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:     "Test no error is returned",
			args:     args{sync: &ResourceCrud{D: &mockResourceData{}, id: "1"}, timeout: time.Second, operationName: "", pending: []string{}, target: []string{"SUCCEEDED"}},
			gotError: false,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "SUCCEEDED", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
			},
		},
		{
			name:     "Test resource failed error is returned",
			args:     args{sync: &ResourceCrud{D: &mockResourceData{}}, timeout: time.Second, operationName: "", pending: []string{}, target: []string{"FAILED"}},
			gotError: true,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "FAILED", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
			},
		},
		{
			name:     "Test *resource.UnexpectedStateError error is returned",
			args:     args{sync: &ResourceCrud{D: &mockResourceData{}}, timeout: time.Second, operationName: "", pending: []string{"A"}, target: []string{"A"}},
			gotError: true,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "ABC", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
			},
		},
		{
			name:     "Test *resource.UnexpectedStateError error is returned",
			args:     args{sync: &ResourceCrud{D: &mockResourceData{}}, timeout: time.Second, operationName: "", pending: []string{"A"}, target: []string{}},
			gotError: true,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "ABC", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
			},
		},
		{
			name:     "Test *resource.TimeoutError error is returned",
			args:     args{sync: &ResourceCrud{D: &mockResourceData{}}, timeout: 0, operationName: "", pending: []string{}, target: []string{}},
			gotError: true,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "ABC", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
			},
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		test.mockFunc()
		if res := WaitForStateRefresh(test.args.sync, test.args.timeout, test.args.operationName, test.args.pending, test.args.target); (res != nil) != test.gotError {
			t.Log(res)
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.gotError)
		}
	}
}

func TestUnitWaitForWorkRequestWithErrorHandling(t *testing.T) {
	type output struct {
		identifier string
		gotError   bool
	}
	type args struct {
		workRequestClient   *mockWorkRequestClient
		workRequestIds      *string
		entityType          string
		action              oci_work_requests.WorkRequestResourceActionTypeEnum
		timeout             time.Duration
		disableFoundRetries bool
	}
	type testFormat struct {
		name     string
		args     args
		output   output
		mockFunc func()
	}
	timeoutDuration, _ := time.ParseDuration("1s")
	workReqIds := "1 2 3"
	tests := []testFormat{
		{
			name:   "Test error is not returned",
			args:   args{workRequestClient: nil, workRequestIds: &workReqIds, entityType: "", action: "", timeout: timeoutDuration, disableFoundRetries: false},
			output: output{identifier: "test", gotError: false},
			mockFunc: func() {
				WaitForWorkRequestVar = func(wrc workReqClient, wId *string, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum, tt time.Duration, dfr bool, ei bool) (*string, error) {
					id := "test"
					return &id, nil
				}
			},
		},
		{
			name:   "Test error is returned",
			args:   args{workRequestClient: nil, workRequestIds: &workReqIds, entityType: "", action: "", timeout: timeoutDuration, disableFoundRetries: false},
			output: output{identifier: "test", gotError: true},
			mockFunc: func() {
				WaitForWorkRequestVar = func(wrc workReqClient, wId *string, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum, tt time.Duration, dfr bool, ei bool) (*string, error) {
					id := "test"
					return &id, errors.New("default")
				}
			},
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		test.mockFunc()
		id, err := WaitForWorkRequestWithErrorHandling(test.args.workRequestClient, test.args.workRequestIds, test.args.entityType, test.args.action, test.args.timeout, test.args.disableFoundRetries)
		if (err != nil) != test.output.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", err, test.output.gotError)
		}
		if *id != test.output.identifier {
			t.Errorf("Output string - %s is not equal to expected string - %s", *id, test.output.identifier)
		}
	}
}

func TestUnitWaitForWorkRequest(t *testing.T) {
	type output struct {
		identifier string
		gotError   bool
	}
	type args struct {
		workRequestClient   *mockWorkRequestClient
		workRequestId       *string
		entityType          string
		action              oci_work_requests.WorkRequestResourceActionTypeEnum
		timeout             time.Duration
		disableFoundRetries bool
		expectIdentifier    bool
	}
	type testFormat struct {
		name      string
		workReqId string
		args      args
		output    output
		mockFunc  func()
	}
	timeoutDuration, _ := time.ParseDuration("1s")
	tests := []testFormat{
		{
			name:      "Test error is not returned",
			workReqId: "1",
			args:      args{workRequestClient: nil, entityType: "", action: "CREATED", timeout: timeoutDuration, disableFoundRetries: false, expectIdentifier: true},
			output:    output{identifier: "oci", gotError: false},
			mockFunc: func() {
				getWorkRequestErrorsVar = func(wrc workReqClient, wId *string, rp *oci_common.RetryPolicy, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum) error {
					return nil
				}
			},
		},
		{
			name:      "Test error is returned",
			workReqId: "2",
			args:      args{workRequestClient: nil, entityType: "default", action: "CREATED", timeout: timeoutDuration, disableFoundRetries: false, expectIdentifier: true},
			output:    output{identifier: "", gotError: true},
			mockFunc: func() {
				getWorkRequestErrorsVar = func(wrc workReqClient, wId *string, rp *oci_common.RetryPolicy, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum) error {
					return errors.New("")
				}
			},
		},
		{
			name:      "Test error is returned",
			workReqId: "3",
			args:      args{workRequestClient: nil, entityType: "default", action: "", timeout: timeoutDuration, disableFoundRetries: false, expectIdentifier: true},
			output:    output{identifier: "", gotError: true},
			mockFunc: func() {
				getWorkRequestErrorsVar = func(wrc workReqClient, wId *string, rp *oci_common.RetryPolicy, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum) error {
					return errors.New("")
				}
			},
		},
		{
			name:      "Test error is returned",
			workReqId: "default",
			args:      args{workRequestClient: nil, entityType: "default", action: "", timeout: timeoutDuration, disableFoundRetries: false, expectIdentifier: true},
			output:    output{identifier: "", gotError: true},
			mockFunc: func() {
				getWorkRequestErrorsVar = func(wrc workReqClient, wId *string, rp *oci_common.RetryPolicy, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum) error {
					return errors.New("")
				}
			},
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		test.mockFunc()
		id, err := WaitForWorkRequest(test.args.workRequestClient, &test.workReqId, test.args.entityType, test.args.action, test.args.timeout, test.args.disableFoundRetries, test.args.expectIdentifier)
		if id != nil && *id != test.output.identifier {
			t.Log(*id)
			t.Errorf("Output identifier - %s not equal to expected identifier - %s", *id, test.output.identifier)
		}
		if (err != nil) != test.output.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", err, test.output.gotError)
		}
	}
}

func TestUnitgetWorkRequestErrors(t *testing.T) {
	type args struct {
		workRequestClient *mockWorkRequestClient
		workRequestId     *string
		retryPolicy       *oci_common.RetryPolicy
		entityType        string
		action            oci_work_requests.WorkRequestResourceActionTypeEnum
	}
	type testFormat struct {
		name   string
		args   args
		output bool
	}
	workReqId := "1"
	tests := []testFormat{
		{
			name:   "Test error is not returned",
			args:   args{workRequestClient: nil, workRequestId: &workReqId, entityType: "", action: "CREATED"},
			output: true,
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		if res := getWorkRequestErrors(test.args.workRequestClient, test.args.workRequestId, test.args.retryPolicy, test.args.entityType, test.args.action); (res != nil) != test.output {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.output)
		}
	}
}

func TestUnitwaitForStateRefreshForHybridPolling(t *testing.T) {
	workReqIds := "1 2 3"
	type args struct {
		workRequestClient   *mockWorkRequestClient
		workRequestIds      *string
		entityType          string
		action              oci_work_requests.WorkRequestResourceActionTypeEnum
		disableFoundRetries bool
		sync                StatefulResource
		timeout             time.Duration
		operationName       string
		pending             []string
		target              []string
	}
	type testFormat struct {
		name     string
		args     args
		output   bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:   "Test no error is returned",
			args:   args{workRequestClient: nil, workRequestIds: &workReqIds, entityType: "", action: "CREATED", disableFoundRetries: false, sync: &ResourceCrud{id: "1"}, timeout: time.Second, operationName: "default", pending: []string{}, target: []string{"SUCCEEDED"}},
			output: false,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "SUCCEEDED", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
				getWorkRequestErrorsVar = func(wrc workReqClient, wId *string, rp *oci_common.RetryPolicy, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum) error {
					return errors.New("")
				}
			},
		},
		{
			name:   "Test resource failed error is returned",
			args:   args{workRequestClient: nil, workRequestIds: &workReqIds, entityType: "", action: "CREATED", disableFoundRetries: false, sync: &ResourceCrud{}, timeout: time.Second, operationName: "default", pending: []string{}, target: []string{"FAILED"}},
			output: true,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "FAILED", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
				getWorkRequestErrorsVar = func(wrc workReqClient, wId *string, rp *oci_common.RetryPolicy, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum) error {
					return errors.New("")
				}
			},
		},
		{
			name:   "Test *resource.UnexpectedStateError error is returned",
			args:   args{workRequestClient: nil, workRequestIds: &workReqIds, entityType: "", action: "CREATED", disableFoundRetries: false, sync: &ResourceCrud{}, timeout: time.Second, operationName: "default", pending: []string{"A"}, target: []string{"A"}},
			output: true,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "ABC", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
				getWorkRequestErrorsVar = func(wrc workReqClient, wId *string, rp *oci_common.RetryPolicy, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum) error {
					return errors.New("")
				}
			},
		},
		{
			name:   "Test *resource.TimeoutError error is returned",
			args:   args{workRequestClient: nil, workRequestIds: &workReqIds, entityType: "", action: "CREATED", disableFoundRetries: false, sync: &ResourceCrud{}, timeout: 0, operationName: "default", pending: []string{}, target: []string{}},
			output: true,
			mockFunc: func() {
				stateRefreshFuncVar = func(sync StatefulResource) resource.StateRefreshFunc {
					return func() (res interface{}, s string, e error) {
						wr := oci_work_requests.WorkRequest{Status: "ABC", Resources: []oci_work_requests.WorkRequestResource{{EntityType: nil, Identifier: nil, ActionType: "CREATED"}}}
						return wr, string(wr.Status), nil
					}
				}
				getWorkRequestErrorsVar = func(wrc workReqClient, wId *string, rp *oci_common.RetryPolicy, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum) error {
					return errors.New("")
				}
			},
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		test.mockFunc()
		if res := waitForStateRefreshForHybridPolling(test.args.workRequestClient, test.args.workRequestIds, test.args.entityType, test.args.action, test.args.disableFoundRetries, test.args.sync, test.args.timeout, test.args.operationName, test.args.pending, test.args.target); (res != nil) != test.output {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.output)
		}
	}
}

func TestUnitstateRefreshFunc(t *testing.T) {
	type args struct {
		sync StatefulResource
	}
	type testFormat struct {
		name   string
		args   args
		output bool
	}
	tests := []testFormat{
		{
			name:   "Test no error is returned",
			args:   args{sync: &ResourceCrud{id: "1"}},
			output: false,
		},
		{
			name:   "Test sync.Get() error is returned",
			args:   args{sync: &ResourceCrud{id: "2"}},
			output: true,
		},
		{
			name:   "Test sync.setState() error is returned",
			args:   args{sync: &ResourceCrud{id: "3"}},
			output: true,
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		if _, _, err := stateRefreshFunc(test.args.sync)(); (err != nil) != test.output {
			t.Errorf("Output error - %q which is not equal to expected error - %t", err, test.output)
		}
	}
}

func TestUnitResourceRefreshForHybridPolling(t *testing.T) {
	s := &ResourceCrud{}
	workReqIds := "1 2 3"
	reqResourceData := &mockResourceData{}
	type args struct {
		workRequestClient   *mockWorkRequestClient
		workRequestIds      *string
		entityType          string
		action              oci_work_requests.WorkRequestResourceActionTypeEnum
		disableFoundRetries bool
		d                   *mockResourceData
		sync                ResourceCreator
	}
	type testFormat struct {
		name     string
		args     args
		output   bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:   "Test error is returned",
			args:   args{workRequestClient: nil, workRequestIds: &workReqIds, entityType: "", action: "CREATED", disableFoundRetries: false, d: reqResourceData, sync: s},
			output: true,
			mockFunc: func() {
				waitForStateRefreshForHybridPollingVar = func(wrc workReqClient, wIds *string, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum, dfr bool, s StatefulResource, tt time.Duration, on string, p []string, t []string) error {
					return errors.New("")
				}
			},
		},
		{
			name:   "Test error is not returned",
			args:   args{workRequestClient: nil, workRequestIds: &workReqIds, entityType: "", action: "CREATED", disableFoundRetries: false, d: reqResourceData, sync: s},
			output: false,
			mockFunc: func() {
				waitForStateRefreshForHybridPollingVar = func(wrc workReqClient, wIds *string, et string, a oci_work_requests.WorkRequestResourceActionTypeEnum, dfr bool, s StatefulResource, tt time.Duration, on string, p []string, t []string) error {
					return nil
				}
			},
		},
	}
	for _, test := range tests {
		t.Logf("Running %s", test.name)
		test.mockFunc()
		if res := ResourceRefreshForHybridPolling(test.args.workRequestClient, test.args.workRequestIds, test.args.entityType, test.args.action, test.args.disableFoundRetries, test.args.d, test.args.sync); (res != nil) != test.output {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.output)
		}
	}
}

func TestUnitCreateResourceUsingHybridPolling(t *testing.T) {
	type args struct {
		sync ResourceCreator
	}
	type testFormat struct {
		name     string
		args     args
		output   bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:   "Test sync.Create() error is returned",
			args:   args{sync: &ResourceCrud{id: "4"}},
			output: true,
			mockFunc: func() {
				HandleErrorVar = func(sync interface{}, e error) error {
					return errors.New("")
				}
			},
		},
		{
			name:   "Test no error is returned",
			args:   args{sync: &ResourceCrud{}},
			output: false,
			mockFunc: func() {
				HandleErrorVar = func(sync interface{}, e error) error {
					return errors.New("")
				}
			},
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		test.mockFunc()
		if res := CreateResourceUsingHybridPolling(test.args.sync); (res != nil) != test.output {
			t.Errorf("Output error - %q which is not equal to expected error - %t", res, test.output)
		}
	}
}

func TestUnitFilterMissingResourceError(t *testing.T) {
	type args struct {
		sync ResourceVoider
		err  error
	}
	type testFormat struct {
		name string
		args args
	}
	tests := []testFormat{
		{
			name: "Test error is returned",
			args: args{sync: &ResourceCrud{}, err: errors.New("does not exist")},
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		FilterMissingResourceError(test.args.sync, &test.args.err)
	}
}

// TBI
func TestUnitGenerateDataSourceID(t *testing.T) {
	GenerateDataSourceID()
}

// TBI
func TestUnitStringsToSet(t *testing.T) {
	type args struct {
		ss []string
	}
	type testFormat struct {
		name string
		args args
	}
	tests := []testFormat{
		{
			name: "Test no error is returned",
			args: args{ss: []string{"oci"}},
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		st := StringsToSet(test.args.ss)
		t.Log(st)
	}
}

func TestUnitGetRetryPolicyWithAdditionalRetryCondition(t *testing.T) {
	type args struct {
		timeout                time.Duration
		retryConditionFunction func(oci_common.OCIOperationResponse) bool
		service                string
	}
	type output struct {
		ShouldRetryOperation  bool
		NextDuration          time.Duration
		MaximumNumberAttempts int
	}
	type testFormat struct {
		name     string
		args     args
		output   output
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:   "Test ShouldRetryOperation returns true",
			args:   args{timeout: time.Second, retryConditionFunction: func(o oci_common.OCIOperationResponse) bool { return true }, service: ""},
			output: output{ShouldRetryOperation: true, MaximumNumberAttempts: 0},
			mockFunc: func() {
				ShouldRetryVar = func(response oci_common.OCIOperationResponse, disableNotFoundRetries bool, service string, startTime time.Time, optionals ...interface {
				}) bool {
					return true
				}
			},
		},
		{
			name:   "Test ShouldRetryOperation returns true",
			args:   args{timeout: time.Second, retryConditionFunction: func(o oci_common.OCIOperationResponse) bool { return true }, service: ""},
			output: output{ShouldRetryOperation: true, MaximumNumberAttempts: 0},
			mockFunc: func() {
				ShouldRetryVar = func(response oci_common.OCIOperationResponse, disableNotFoundRetries bool, service string, startTime time.Time, optionals ...interface {
				}) bool {
					return false
				}
			},
		},
		{
			name:   "Test ShouldRetryOperation returns false",
			args:   args{timeout: time.Second, retryConditionFunction: func(o oci_common.OCIOperationResponse) bool { return false }, service: ""},
			output: output{ShouldRetryOperation: false, MaximumNumberAttempts: 0},
			mockFunc: func() {
				ShouldRetryVar = func(response oci_common.OCIOperationResponse, disableNotFoundRetries bool, service string, startTime time.Time, optionals ...interface {
				}) bool {
					return false
				}
			},
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		test.mockFunc()
		var o oci_common.OCIOperationResponse
		x := GetRetryPolicyWithAdditionalRetryCondition(test.args.timeout, test.args.retryConditionFunction, test.args.service)
		y := x.ShouldRetryOperation(o)
		if y != test.output.ShouldRetryOperation {
			t.Errorf("ShouldRetryOperation - expected value %t is not equal to output value %t", y, test.output.ShouldRetryOperation)
		}
	}
}

func TestUnitworkRequestShouldRetryFunc(t *testing.T) {
	type args struct {
		timeout time.Duration
	}
	type testFormat struct {
		name     string
		args     args
		output   bool
		mockFunc func()
	}
	tests := []testFormat{
		{
			name:   "Test returns true",
			args:   args{timeout: time.Second},
			output: true,
			mockFunc: func() {
				ShouldRetryVar = func(response oci_common.OCIOperationResponse, disableNotFoundRetries bool, service string, startTime time.Time, optionals ...interface {
				}) bool {
					return true
				}
			},
		},
		{
			name:   "Test returns false",
			args:   args{timeout: time.Second},
			output: false,
			mockFunc: func() {
				ShouldRetryVar = func(response oci_common.OCIOperationResponse, disableNotFoundRetries bool, service string, startTime time.Time, optionals ...interface {
				}) bool {
					return false
				}
			},
		},
	}
	for _, test := range tests {
		t.Log("Running ", test.name)
		test.mockFunc()
		var response oci_common.OCIOperationResponse
		if res := workRequestShouldRetryFunc(test.args.timeout)(response); res != test.output {
			t.Errorf("expected value %t is not equal to output value %t", res, test.output)
		}
	}
}

func TestUnitConvertObjectToJsonString(t *testing.T) {
	type args struct {
		object interface{}
	}
	type output struct {
		res      string
		gotError bool
	}
	type testFormat struct {
		name     string
		args     args
		output   output
		mockFunc func()
	}

	tests := []testFormat{
		{
			name:   "Test valid slice",
			args:   args{object: []string{"foo", "bar"}},
			output: output{res: "[\"foo\",\"bar\"]", gotError: false},
			mockFunc: func() {
				jsonMarshal = func(object interface{}) ([]byte, error) {
					return []byte{91, 34, 102, 111, 111, 34, 44, 34, 98, 97, 114, 34, 93}, nil
				}
			},
		},
		{
			name:   "Test invalid slice",
			args:   args{object: []string{"foo", "bar"}},
			output: output{res: "", gotError: true},
			mockFunc: func() {
				jsonMarshal = func(object interface{}) ([]byte, error) {
					return []byte{91, 34, 102, 111, 111, 34, 44, 34, 98, 97, 114, 34, 93}, errors.New("")
				}
			},
		},
	}
	for _, test := range tests {
		test.mockFunc()
		res, err := ConvertObjectToJsonString(test.args.object)
		if res != test.output.res {
			t.Errorf("Output %s not equal to expected %s", res, test.output.res)
		}
		if (err != nil) != test.output.gotError {
			t.Errorf("Output error - %q which is not equal to expected error - %t", err, test.output.gotError)
		}
	}
}

func TestUnitWaitForResourceCondition_basic(t *testing.T) {
	//if httpreplay.ModeRecordReplay() {
	//	t.Skip("Skip TestWaitForResourceCondition_basic test in HttpReplay mode.")
	//}
	getAttempts := 1
	testResource := &TestResource{GetError: nil, GetAttempts: getAttempts}
	finalStateFunc := func() bool {
		return testResource.GetAttempts == 0
	}

	// Test normal case
	err := WaitForResourceCondition(testResource, finalStateFunc, 0)
	if err != nil {
		t.Errorf("Got unexpected error '%q' from single attempt", err)
		return
	}

	// Test normal case with multiple attempts
	testResource = &TestResource{GetError: nil, GetAttempts: 3}
	err = WaitForResourceCondition(testResource, finalStateFunc, time.Minute)
	if err != nil {
		t.Errorf("Got unexpected error '%q' from multiple attempts", err)
		return
	}

	// Test case where Get returns error after 1 attempt
	testResource = &TestResource{GetError: fmt.Errorf("GetError"), GetAttempts: 1}
	err = WaitForResourceCondition(testResource, finalStateFunc, 0)
	if err == nil || !strings.HasPrefix(err.Error(), "GetError") {
		t.Errorf("Got unexpected error '%q' after single attempt, expected a GetError", err)
		return
	}

	// Test case where Get returns error after multiple attempts
	testResource = &TestResource{GetError: fmt.Errorf("GetError"), GetAttempts: 3}
	err = WaitForResourceCondition(testResource, finalStateFunc, time.Minute)
	if err == nil || !strings.HasPrefix(err.Error(), "GetError") {
		t.Errorf("Got unexpected error '%q' after multiple attempts, expected a GetError", err)
		return
	}

	// Test timing out with zero timeout duration
	testResource = &TestResource{GetError: nil, GetAttempts: 10}
	err = WaitForResourceCondition(testResource, finalStateFunc, 0)
	if err == nil || !strings.HasPrefix(err.Error(), "Timed out") {
		t.Errorf("Got unexpected error '%q' after a single attempt, expected a timeout error", err)
		return
	}

	// Test timing out with non-zero timeout duration, also validate that we got expected number of Get attempts due to exponential backoff
	testResource = &TestResource{GetError: nil, GetAttempts: 10}
	err = WaitForResourceCondition(testResource, finalStateFunc, 20*time.Second)
	if err == nil || !strings.HasPrefix(err.Error(), "Timed out") {
		t.Errorf("Got unexpected error '%q' after a single attempt, expected a timeout error", err)
		return
	}

	// Expected Get attempts at: 0, 2, 6, 14, 20 seconds
	if testResource.ActualGetAttempts != 5 {
		t.Errorf("Expected 5 Get attempts, got %d instead", testResource.ActualGetAttempts)
		return
	}
}
