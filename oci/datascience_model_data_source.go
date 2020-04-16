// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package oci

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	oci_datascience "github.com/oracle/oci-go-sdk/datascience"
)

func init() {
	RegisterDatasource("oci_datascience_model", DatascienceModelDataSource())
}

func DatascienceModelDataSource() *schema.Resource {
	fieldMap := make(map[string]*schema.Schema)
	fieldMap["model_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return GetSingularDataSourceItemSchema(DatascienceModelResource(), fieldMap, readSingularDatascienceModel)
}

func readSingularDatascienceModel(d *schema.ResourceData, m interface{}) error {
	sync := &DatascienceModelDataSourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).dataScienceClient

	return ReadResource(sync)
}

type DatascienceModelDataSourceCrud struct {
	D               *schema.ResourceData
	Client          *oci_datascience.DataScienceClient
	Res             *oci_datascience.GetModelResponse
	ArtifactHeadRes *HeadModelArtifact
}

func (s *DatascienceModelDataSourceCrud) VoidState() {
	s.D.SetId("")
}

func (s *DatascienceModelDataSourceCrud) Get() error {
	request := oci_datascience.GetModelRequest{}

	if modelId, ok := s.D.GetOkExists("model_id"); ok {
		tmp := modelId.(string)
		request.ModelId = &tmp
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(false, "datascience")

	response, err := s.Client.GetModel(context.Background(), request)
	if err != nil {
		return err
	}

	s.Res = &response
	headModelArtifactRequest := oci_datascience.HeadModelArtifactRequest{}

	if modelId, ok := s.D.GetOkExists("model_id"); ok {
		tmp := modelId.(string)
		headModelArtifactRequest.ModelId = &tmp
	}

	headModelArtifactRequest.RequestMetadata.RetryPolicy = getRetryPolicy(false, "datascience")

	headModelArtifactResponse, err := s.Client.HeadModelArtifact(context.Background(), headModelArtifactRequest)
	if err != nil {
		return err
	}

	s.ArtifactHeadRes = &HeadModelArtifact{
		ContentLength:      headModelArtifactResponse.ContentLength,
		ContentDisposition: headModelArtifactResponse.ContentDisposition,
		ContentMd5:         headModelArtifactResponse.ContentMd5,
		LastModified:       headModelArtifactResponse.LastModified,
	}
	return nil
}

func (s *DatascienceModelDataSourceCrud) SetData() error {
	if s.Res == nil {
		return nil
	}

	s.D.SetId(*s.Res.Id)

	if s.Res.CompartmentId != nil {
		s.D.Set("compartment_id", *s.Res.CompartmentId)
	}

	if s.Res.CreatedBy != nil {
		s.D.Set("created_by", *s.Res.CreatedBy)
	}

	if s.Res.DefinedTags != nil {
		s.D.Set("defined_tags", definedTagsToMap(s.Res.DefinedTags))
	}

	if s.Res.Description != nil {
		s.D.Set("description", *s.Res.Description)
	}

	if s.Res.DisplayName != nil {
		s.D.Set("display_name", *s.Res.DisplayName)
	}

	s.D.Set("freeform_tags", s.Res.FreeformTags)

	if s.Res.ProjectId != nil {
		s.D.Set("project_id", *s.Res.ProjectId)
	}

	s.D.Set("state", s.Res.LifecycleState)

	if s.Res.TimeCreated != nil {
		s.D.Set("time_created", s.Res.TimeCreated.String())
	}

	if s.ArtifactHeadRes.ContentDisposition != nil {
		s.D.Set("artifact_content_disposition", *s.ArtifactHeadRes.ContentDisposition)
	}

	if s.ArtifactHeadRes.ContentLength != nil {
		s.D.Set("artifact_content_length", *s.ArtifactHeadRes.ContentLength)
	}

	if s.ArtifactHeadRes.ContentMd5 != nil {
		s.D.Set("artifact_content_md5", *s.ArtifactHeadRes.ContentMd5)
	}

	if s.ArtifactHeadRes.LastModified != nil {
		s.D.Set("artifact_last_modified", s.ArtifactHeadRes.LastModified.String())
	}

	return nil
}