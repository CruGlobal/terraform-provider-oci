// Copyright (c) 2016, 2018, 2023, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// OS Management Hub API
//
// Use the OS Management Hub API to manage and monitor updates and patches for the operating system environments in your private data centers through a single management console. For more information, see Overview of OS Management Hub (https://docs.cloud.oracle.com/iaas/osmh/doc/overview.htm).
//

package osmanagementhub

import (
	"fmt"
	"github.com/oracle/oci-go-sdk/v65/common"
	"strings"
)

// ModuleStreamProfileFilter Used to select module stream/profiles from VendorSoftwareSources to create/update CustomSoftwareSources.
type ModuleStreamProfileFilter struct {

	// Module name.
	ModuleName *string `mandatory:"true" json:"moduleName"`

	// The type of the filter, which can be of two types - INCLUDE or EXCLUDE.
	FilterType FilterTypeEnum `mandatory:"true" json:"filterType"`

	// Profile name.
	ProfileName *string `mandatory:"false" json:"profileName"`

	// Stream name.
	StreamName *string `mandatory:"false" json:"streamName"`
}

func (m ModuleStreamProfileFilter) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ModuleStreamProfileFilter) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingFilterTypeEnum(string(m.FilterType)); !ok && m.FilterType != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for FilterType: %s. Supported values are: %s.", m.FilterType, strings.Join(GetFilterTypeEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}
