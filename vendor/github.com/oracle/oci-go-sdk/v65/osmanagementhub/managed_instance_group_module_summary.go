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

// ManagedInstanceGroupModuleSummary Summary information pertaining to a module on a managed instance group.
type ManagedInstanceGroupModuleSummary struct {

	// The name of the module that contains the stream.
	Name *string `mandatory:"true" json:"name"`

	// The name of the module that contains the stream.
	EnabledStream *string `mandatory:"false" json:"enabledStream"`

	// The list of installed profiles under the currently enabled module stream.
	InstalledProfiles []string `mandatory:"false" json:"installedProfiles"`

	// The OCID of the software source that provides this module stream.
	SoftwareSourceId *string `mandatory:"false" json:"softwareSourceId"`
}

func (m ManagedInstanceGroupModuleSummary) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ManagedInstanceGroupModuleSummary) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}
