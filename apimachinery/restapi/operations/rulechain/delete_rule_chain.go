// Code generated by go-swagger; DO NOT EDIT.

package rulechain

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/cloustone/pandas/models"
)

// DeleteRuleChainHandlerFunc turns a function with the right signature into a delete rule chain handler
type DeleteRuleChainHandlerFunc func(DeleteRuleChainParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteRuleChainHandlerFunc) Handle(params DeleteRuleChainParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// DeleteRuleChainHandler interface for that can handle valid delete rule chain params
type DeleteRuleChainHandler interface {
	Handle(DeleteRuleChainParams, *models.Principal) middleware.Responder
}

// NewDeleteRuleChain creates a new http.Handler for the delete rule chain operation
func NewDeleteRuleChain(ctx *middleware.Context, handler DeleteRuleChainHandler) *DeleteRuleChain {
	return &DeleteRuleChain{Context: ctx, Handler: handler}
}

/*DeleteRuleChain swagger:route DELETE /rulechains/{ruleChainId} Rulechain deleteRuleChain

delete rule chain

delete rule chain with Id

*/
type DeleteRuleChain struct {
	Context *middleware.Context
	Handler DeleteRuleChainHandler
}

func (o *DeleteRuleChain) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteRuleChainParams()

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