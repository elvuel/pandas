// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/cloustone/pandas/models"
)

// UpdateProjectDeviceHandlerFunc turns a function with the right signature into a update project device handler
type UpdateProjectDeviceHandlerFunc func(UpdateProjectDeviceParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateProjectDeviceHandlerFunc) Handle(params UpdateProjectDeviceParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// UpdateProjectDeviceHandler interface for that can handle valid update project device params
type UpdateProjectDeviceHandler interface {
	Handle(UpdateProjectDeviceParams, *models.Principal) middleware.Responder
}

// NewUpdateProjectDevice creates a new http.Handler for the update project device operation
func NewUpdateProjectDevice(ctx *middleware.Context, handler UpdateProjectDeviceHandler) *UpdateProjectDevice {
	return &UpdateProjectDevice{Context: ctx, Handler: handler}
}

/*UpdateProjectDevice swagger:route PUT /project/{projectId}/devices/{deviceId} Project updateProjectDevice

update project's specified device

update project's specified device

*/
type UpdateProjectDevice struct {
	Context *middleware.Context
	Handler UpdateProjectDeviceHandler
}

func (o *UpdateProjectDevice) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateProjectDeviceParams()

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