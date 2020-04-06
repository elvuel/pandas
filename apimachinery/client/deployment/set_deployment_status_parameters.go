// Code generated by go-swagger; DO NOT EDIT.

package deployment

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

	"github.com/cloustone/pandas/apimachinery/models"
)

// NewSetDeploymentStatusParams creates a new SetDeploymentStatusParams object
// with the default values initialized.
func NewSetDeploymentStatusParams() *SetDeploymentStatusParams {
	var ()
	return &SetDeploymentStatusParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewSetDeploymentStatusParamsWithTimeout creates a new SetDeploymentStatusParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewSetDeploymentStatusParamsWithTimeout(timeout time.Duration) *SetDeploymentStatusParams {
	var ()
	return &SetDeploymentStatusParams{

		timeout: timeout,
	}
}

// NewSetDeploymentStatusParamsWithContext creates a new SetDeploymentStatusParams object
// with the default values initialized, and the ability to set a context for a request
func NewSetDeploymentStatusParamsWithContext(ctx context.Context) *SetDeploymentStatusParams {
	var ()
	return &SetDeploymentStatusParams{

		Context: ctx,
	}
}

// NewSetDeploymentStatusParamsWithHTTPClient creates a new SetDeploymentStatusParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewSetDeploymentStatusParamsWithHTTPClient(client *http.Client) *SetDeploymentStatusParams {
	var ()
	return &SetDeploymentStatusParams{
		HTTPClient: client,
	}
}

/*SetDeploymentStatusParams contains all the parameters to send to the API endpoint
for the set deployment status operation typically these are written to a http.Request
*/
type SetDeploymentStatusParams struct {

	/*DeploymentControl
	  start or stop deployment

	*/
	DeploymentControl *models.DeploymentControl
	/*DeploymentID
	  deployment identifier

	*/
	DeploymentID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the set deployment status params
func (o *SetDeploymentStatusParams) WithTimeout(timeout time.Duration) *SetDeploymentStatusParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set deployment status params
func (o *SetDeploymentStatusParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set deployment status params
func (o *SetDeploymentStatusParams) WithContext(ctx context.Context) *SetDeploymentStatusParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set deployment status params
func (o *SetDeploymentStatusParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set deployment status params
func (o *SetDeploymentStatusParams) WithHTTPClient(client *http.Client) *SetDeploymentStatusParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set deployment status params
func (o *SetDeploymentStatusParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDeploymentControl adds the deploymentControl to the set deployment status params
func (o *SetDeploymentStatusParams) WithDeploymentControl(deploymentControl *models.DeploymentControl) *SetDeploymentStatusParams {
	o.SetDeploymentControl(deploymentControl)
	return o
}

// SetDeploymentControl adds the deploymentControl to the set deployment status params
func (o *SetDeploymentStatusParams) SetDeploymentControl(deploymentControl *models.DeploymentControl) {
	o.DeploymentControl = deploymentControl
}

// WithDeploymentID adds the deploymentID to the set deployment status params
func (o *SetDeploymentStatusParams) WithDeploymentID(deploymentID string) *SetDeploymentStatusParams {
	o.SetDeploymentID(deploymentID)
	return o
}

// SetDeploymentID adds the deploymentId to the set deployment status params
func (o *SetDeploymentStatusParams) SetDeploymentID(deploymentID string) {
	o.DeploymentID = deploymentID
}

// WriteToRequest writes these params to a swagger request
func (o *SetDeploymentStatusParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.DeploymentControl != nil {
		if err := r.SetBodyParam(o.DeploymentControl); err != nil {
			return err
		}
	}

	// path param deploymentId
	if err := r.SetPathParam("deploymentId", o.DeploymentID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
