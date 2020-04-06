// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/cloustone/pandas/apimachinery/models"
)

// GetModelsHandlerFunc turns a function with the right signature into a get models handler
type GetModelsHandlerFunc func(GetModelsParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetModelsHandlerFunc) Handle(params GetModelsParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetModelsHandler interface for that can handle valid get models params
type GetModelsHandler interface {
	Handle(GetModelsParams, *models.Principal) middleware.Responder
}

// NewGetModels creates a new http.Handler for the get models operation
func NewGetModels(ctx *middleware.Context, handler GetModelsHandler) *GetModels {
	return &GetModels{Context: ctx, Handler: handler}
}

/*GetModels swagger:route GET /models Model getModels

get specified all models

*/
type GetModels struct {
	Context *middleware.Context
	Handler GetModelsHandler
}

func (o *GetModels) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetModelsParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
