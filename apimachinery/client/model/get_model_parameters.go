// Code generated by go-swagger; DO NOT EDIT.

package model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetModelParams creates a new GetModelParams object
// with the default values initialized.
func NewGetModelParams() *GetModelParams {
	var ()
	return &GetModelParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetModelParamsWithTimeout creates a new GetModelParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetModelParamsWithTimeout(timeout time.Duration) *GetModelParams {
	var ()
	return &GetModelParams{

		timeout: timeout,
	}
}

// NewGetModelParamsWithContext creates a new GetModelParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetModelParamsWithContext(ctx context.Context) *GetModelParams {
	var ()
	return &GetModelParams{

		Context: ctx,
	}
}

// NewGetModelParamsWithHTTPClient creates a new GetModelParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetModelParamsWithHTTPClient(client *http.Client) *GetModelParams {
	var ()
	return &GetModelParams{
		HTTPClient: client,
	}
}

/*GetModelParams contains all the parameters to send to the API endpoint
for the get model operation typically these are written to a http.Request
*/
type GetModelParams struct {

	/*ModelID
	  device model id

	*/
	ModelID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get model params
func (o *GetModelParams) WithTimeout(timeout time.Duration) *GetModelParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get model params
func (o *GetModelParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get model params
func (o *GetModelParams) WithContext(ctx context.Context) *GetModelParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get model params
func (o *GetModelParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get model params
func (o *GetModelParams) WithHTTPClient(client *http.Client) *GetModelParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get model params
func (o *GetModelParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithModelID adds the modelID to the get model params
func (o *GetModelParams) WithModelID(modelID string) *GetModelParams {
	o.SetModelID(modelID)
	return o
}

// SetModelID adds the modelId to the get model params
func (o *GetModelParams) SetModelID(modelID string) {
	o.ModelID = modelID
}

// WriteToRequest writes these params to a swagger request
func (o *GetModelParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param modelId
	if err := r.SetPathParam("modelId", o.ModelID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}