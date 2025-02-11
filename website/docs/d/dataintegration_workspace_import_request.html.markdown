---
subcategory: "Data Integration"
layout: "oci"
page_title: "Oracle Cloud Infrastructure: oci_dataintegration_workspace_import_request"
sidebar_current: "docs-oci-datasource-dataintegration-workspace_import_request"
description: |-
  Provides details about a specific Workspace Import Request in Oracle Cloud Infrastructure Data Integration service
---

# Data Source: oci_dataintegration_workspace_import_request
This data source provides details about a specific Workspace Import Request resource in Oracle Cloud Infrastructure Data Integration service.

This endpoint can be used to get the summary/details of object being imported.


## Example Usage

```hcl
data "oci_dataintegration_workspace_import_request" "test_workspace_import_request" {
	#Required
	import_request_key = var.workspace_import_request_import_request_key
	workspace_id = oci_dataintegration_workspace.test_workspace.id
}
```

## Argument Reference

The following arguments are supported:

* `import_request_key` - (Required) The key of the object export object request
* `workspace_id` - (Required) The workspace ID.


## Attributes Reference

The following attributes are exported:

* `bucket` - The name of the Object Storage Bucket where the objects will be imported from
* `created_by` - Name of the user who initiated import request.
* `error_messages` - Contains key of the error
* `file_name` - Name of the zip file from which objects will be imported.
* `import_conflict_resolution` - Import Objects Conflict resolution.
	* `duplicate_prefix` - In case of DUPLICATE mode, prefix will be used to disambiguate the object.
	* `duplicate_suffix` - In case of DUPLICATE mode, suffix will be used to disambiguate the object.
	* `import_conflict_resolution_type` - Import Objects Conflict resolution Type (RETAIN/DUPLICATE/REPLACE).
* `imported_objects` - The array of imported object details.
	* `aggregator_key` - Aggregator key
	* `identifier` - Object identifier
	* `name` - Name of the object
	* `name_path` - Object name path
	* `new_key` - New key of the object
	* `object_type` - Object type
	* `object_version` - Object version
	* `old_key` - Old key of the object
	* `resolution_action` - Object resolution action
	* `time_updated_in_millis` - time at which this object was last updated.
* `key` - Import object request key
* `name` - Name of the import request.
* `object_key_for_import` - Key of the object inside which all the objects will be imported
* `object_storage_region` - Region of the object storage (if using object storage of different region)
* `object_storage_tenancy_id` - Optional parameter to point to object storage tenancy (if using Object Storage of different tenancy)
* `status` - Import Objects request status.
* `time_ended_in_millis` - Time at which the request was completely processed.
* `time_started_in_millis` - Time at which the request started getting processed.
* `total_imported_object_count` - Number of objects that are imported.

