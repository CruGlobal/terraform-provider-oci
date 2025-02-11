// Copyright (c) 2016, 2018, 2023, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Database Management API
//
// Use the Database Management API to perform tasks such as obtaining performance and resource usage metrics
// for a fleet of Managed Databases or a specific Managed Database, creating Managed Database Groups, and
// running a SQL job on a Managed Database or Managed Database Group.
//

package databasemanagement

import (
	"encoding/json"
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/common"
	"strings"
)

// ConfigureAutomaticCaptureFiltersDetails The details required to configure automatic capture filters.
type ConfigureAutomaticCaptureFiltersDetails struct {

	// The filters used in automatic initial plan capture.
	AutoCaptureFilters []AutomaticCaptureFilterDetails `mandatory:"true" json:"autoCaptureFilters"`

	Credentials ManagedDatabaseCredential `mandatory:"true" json:"credentials"`
}

func (m ConfigureAutomaticCaptureFiltersDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ConfigureAutomaticCaptureFiltersDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// UnmarshalJSON unmarshals from json
func (m *ConfigureAutomaticCaptureFiltersDetails) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		AutoCaptureFilters []AutomaticCaptureFilterDetails `json:"autoCaptureFilters"`
		Credentials        manageddatabasecredential       `json:"credentials"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	var nn interface{}
	m.AutoCaptureFilters = make([]AutomaticCaptureFilterDetails, len(model.AutoCaptureFilters))
	copy(m.AutoCaptureFilters, model.AutoCaptureFilters)
	nn, e = model.Credentials.UnmarshalPolymorphicJSON(model.Credentials.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.Credentials = nn.(ManagedDatabaseCredential)
	} else {
		m.Credentials = nil
	}

	return
}
