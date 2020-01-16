// Code generated by go-swagger; DO NOT EDIT.

package project

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

	"github.com/cloustone/pandas/models"
)

// NewAddProjectDeviceParams creates a new AddProjectDeviceParams object
// with the default values initialized.
func NewAddProjectDeviceParams() *AddProjectDeviceParams {
	var ()
	return &AddProjectDeviceParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewAddProjectDeviceParamsWithTimeout creates a new AddProjectDeviceParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewAddProjectDeviceParamsWithTimeout(timeout time.Duration) *AddProjectDeviceParams {
	var ()
	return &AddProjectDeviceParams{

		timeout: timeout,
	}
}

// NewAddProjectDeviceParamsWithContext creates a new AddProjectDeviceParams object
// with the default values initialized, and the ability to set a context for a request
func NewAddProjectDeviceParamsWithContext(ctx context.Context) *AddProjectDeviceParams {
	var ()
	return &AddProjectDeviceParams{

		Context: ctx,
	}
}

// NewAddProjectDeviceParamsWithHTTPClient creates a new AddProjectDeviceParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewAddProjectDeviceParamsWithHTTPClient(client *http.Client) *AddProjectDeviceParams {
	var ()
	return &AddProjectDeviceParams{
		HTTPClient: client,
	}
}

/*AddProjectDeviceParams contains all the parameters to send to the API endpoint
for the add project device operation typically these are written to a http.Request
*/
type AddProjectDeviceParams struct {

	/*Device*/
	Device models.Device
	/*ProjectID
	  specified project

	*/
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the add project device params
func (o *AddProjectDeviceParams) WithTimeout(timeout time.Duration) *AddProjectDeviceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the add project device params
func (o *AddProjectDeviceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the add project device params
func (o *AddProjectDeviceParams) WithContext(ctx context.Context) *AddProjectDeviceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the add project device params
func (o *AddProjectDeviceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the add project device params
func (o *AddProjectDeviceParams) WithHTTPClient(client *http.Client) *AddProjectDeviceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the add project device params
func (o *AddProjectDeviceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDevice adds the device to the add project device params
func (o *AddProjectDeviceParams) WithDevice(device models.Device) *AddProjectDeviceParams {
	o.SetDevice(device)
	return o
}

// SetDevice adds the device to the add project device params
func (o *AddProjectDeviceParams) SetDevice(device models.Device) {
	o.Device = device
}

// WithProjectID adds the projectID to the add project device params
func (o *AddProjectDeviceParams) WithProjectID(projectID string) *AddProjectDeviceParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the add project device params
func (o *AddProjectDeviceParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *AddProjectDeviceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Device != nil {
		if err := r.SetBodyParam(o.Device); err != nil {
			return err
		}
	}

	// path param projectId
	if err := r.SetPathParam("projectId", o.ProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}