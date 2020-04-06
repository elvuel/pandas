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

// NewUpdateDeploymentParams creates a new UpdateDeploymentParams object
// with the default values initialized.
func NewUpdateDeploymentParams() *UpdateDeploymentParams {
	var ()
	return &UpdateDeploymentParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateDeploymentParamsWithTimeout creates a new UpdateDeploymentParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdateDeploymentParamsWithTimeout(timeout time.Duration) *UpdateDeploymentParams {
	var ()
	return &UpdateDeploymentParams{

		timeout: timeout,
	}
}

// NewUpdateDeploymentParamsWithContext creates a new UpdateDeploymentParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdateDeploymentParamsWithContext(ctx context.Context) *UpdateDeploymentParams {
	var ()
	return &UpdateDeploymentParams{

		Context: ctx,
	}
}

// NewUpdateDeploymentParamsWithHTTPClient creates a new UpdateDeploymentParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdateDeploymentParamsWithHTTPClient(client *http.Client) *UpdateDeploymentParams {
	var ()
	return &UpdateDeploymentParams{
		HTTPClient: client,
	}
}

/*UpdateDeploymentParams contains all the parameters to send to the API endpoint
for the update deployment operation typically these are written to a http.Request
*/
type UpdateDeploymentParams struct {

	/*Deployment
	  updated deployment

	*/
	Deployment *models.Deployment
	/*DeploymentID
	  deployment identifier

	*/
	DeploymentID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the update deployment params
func (o *UpdateDeploymentParams) WithTimeout(timeout time.Duration) *UpdateDeploymentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update deployment params
func (o *UpdateDeploymentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update deployment params
func (o *UpdateDeploymentParams) WithContext(ctx context.Context) *UpdateDeploymentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update deployment params
func (o *UpdateDeploymentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update deployment params
func (o *UpdateDeploymentParams) WithHTTPClient(client *http.Client) *UpdateDeploymentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update deployment params
func (o *UpdateDeploymentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDeployment adds the deployment to the update deployment params
func (o *UpdateDeploymentParams) WithDeployment(deployment *models.Deployment) *UpdateDeploymentParams {
	o.SetDeployment(deployment)
	return o
}

// SetDeployment adds the deployment to the update deployment params
func (o *UpdateDeploymentParams) SetDeployment(deployment *models.Deployment) {
	o.Deployment = deployment
}

// WithDeploymentID adds the deploymentID to the update deployment params
func (o *UpdateDeploymentParams) WithDeploymentID(deploymentID string) *UpdateDeploymentParams {
	o.SetDeploymentID(deploymentID)
	return o
}

// SetDeploymentID adds the deploymentId to the update deployment params
func (o *UpdateDeploymentParams) SetDeploymentID(deploymentID string) {
	o.DeploymentID = deploymentID
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateDeploymentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Deployment != nil {
		if err := r.SetBodyParam(o.Deployment); err != nil {
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
