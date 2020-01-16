// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/cloustone/pandas/models"
)

// DeleteProjectDeviceHandlerFunc turns a function with the right signature into a delete project device handler
type DeleteProjectDeviceHandlerFunc func(DeleteProjectDeviceParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteProjectDeviceHandlerFunc) Handle(params DeleteProjectDeviceParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// DeleteProjectDeviceHandler interface for that can handle valid delete project device params
type DeleteProjectDeviceHandler interface {
	Handle(DeleteProjectDeviceParams, *models.Principal) middleware.Responder
}

// NewDeleteProjectDevice creates a new http.Handler for the delete project device operation
func NewDeleteProjectDevice(ctx *middleware.Context, handler DeleteProjectDeviceHandler) *DeleteProjectDevice {
	return &DeleteProjectDevice{Context: ctx, Handler: handler}
}

/*DeleteProjectDevice swagger:route DELETE /project/{projectId}/devices/{deviceId} Project deleteProjectDevice

delete project's specified device

delete project's specified device

*/
type DeleteProjectDevice struct {
	Context *middleware.Context
	Handler DeleteProjectDeviceHandler
}

func (o *DeleteProjectDevice) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteProjectDeviceParams()

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