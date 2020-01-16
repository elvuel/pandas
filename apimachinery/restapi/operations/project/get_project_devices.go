// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/cloustone/pandas/models"
)

// GetProjectDevicesHandlerFunc turns a function with the right signature into a get project devices handler
type GetProjectDevicesHandlerFunc func(GetProjectDevicesParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetProjectDevicesHandlerFunc) Handle(params GetProjectDevicesParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetProjectDevicesHandler interface for that can handle valid get project devices params
type GetProjectDevicesHandler interface {
	Handle(GetProjectDevicesParams, *models.Principal) middleware.Responder
}

// NewGetProjectDevices creates a new http.Handler for the get project devices operation
func NewGetProjectDevices(ctx *middleware.Context, handler GetProjectDevicesHandler) *GetProjectDevices {
	return &GetProjectDevices{Context: ctx, Handler: handler}
}

/*GetProjectDevices swagger:route GET /projects/{projectId}/devices Project getProjectDevices

get project's all devices

get project's devices

*/
type GetProjectDevices struct {
	Context *middleware.Context
	Handler GetProjectDevicesHandler
}

func (o *GetProjectDevices) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetProjectDevicesParams()

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