// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/cloustone/pandas/vms"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	saveModelOp               = "save_model"
	saveModelsOp              = "save_model"
	updateModelOp             = "update_model"
	updateModelKeyOp          = "update_model_by_key"
	retrieveModelByIDOp       = "retrieve_model_by_id"
	retrieveModelByKeyOp      = "retrieve_model_by_key"
	retrieveAllModelsOp       = "retrieve_all_models"
	retrieveModelsByChannelOp = "retrieve_vms_by_chan"
	removeModelOp             = "remove_model"
	retrieveModelIDByKeyOp    = "retrieve_id_by_key"
)

var (
	_ vms.ModelRepository = (*modelRepositoryMiddleware)(nil)
	_ vms.ModelCache      = (*modelCacheMiddleware)(nil)
)

type modelRepositoryMiddleware struct {
	tracer opentracing.Tracer
	repo   vms.ModelRepository
}

// ModelRepositoryMiddleware tracks request and their latency, and adds spans
// to context.
func ModelRepositoryMiddleware(tracer opentracing.Tracer, repo vms.ModelRepository) vms.ModelRepository {
	return modelRepositoryMiddleware{
		tracer: tracer,
		repo:   repo,
	}
}

func (trm modelRepositoryMiddleware) Save(ctx context.Context, ths ...vms.Model) ([]vms.Model, error) {
	span := createSpan(ctx, trm.tracer, saveModelsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Save(ctx, ths...)
}

func (trm modelRepositoryMiddleware) Update(ctx context.Context, th vms.Model) error {
	span := createSpan(ctx, trm.tracer, updateModelOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Update(ctx, th)
}

func (trm modelRepositoryMiddleware) RetrieveByID(ctx context.Context, owner, id string) (vms.Model, error) {
	span := createSpan(ctx, trm.tracer, retrieveModelByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveByID(ctx, owner, id)
}

func (trm modelRepositoryMiddleware) Retrieve(ctx context.Context, id string) (vms.Model, error) {
	span := createSpan(ctx, trm.tracer, retrieveModelByIDOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Retrieve(ctx, id)
}

func (trm modelRepositoryMiddleware) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata vms.Metadata) (vms.ModelsPage, error) {
	span := createSpan(ctx, trm.tracer, retrieveAllModelsOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.RetrieveAll(ctx, owner, offset, limit, name, metadata)
}

func (trm modelRepositoryMiddleware) Remove(ctx context.Context, owner, id string) error {
	span := createSpan(ctx, trm.tracer, removeModelOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return trm.repo.Remove(ctx, owner, id)
}

type modelCacheMiddleware struct {
	tracer opentracing.Tracer
	cache  vms.ModelCache
}

// ModelCacheMiddleware tracks request and their latency, and adds spans
// to context.
func ModelCacheMiddleware(tracer opentracing.Tracer, cache vms.ModelCache) vms.ModelCache {
	return modelCacheMiddleware{
		tracer: tracer,
		cache:  cache,
	}
}

func (tcm modelCacheMiddleware) Save(ctx context.Context, modelKey string, modelID string) error {
	span := createSpan(ctx, tcm.tracer, saveModelOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Save(ctx, modelKey, modelID)
}

func (tcm modelCacheMiddleware) ID(ctx context.Context, modelKey string) (string, error) {
	span := createSpan(ctx, tcm.tracer, retrieveModelIDByKeyOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.ID(ctx, modelKey)
}

func (tcm modelCacheMiddleware) Remove(ctx context.Context, modelID string) error {
	span := createSpan(ctx, tcm.tracer, removeModelOp)
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)

	return tcm.cache.Remove(ctx, modelID)
}
