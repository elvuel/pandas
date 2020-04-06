// Code generated by go-swagger; DO NOT EDIT.

package device

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

// NewSendDataToDeviceParams creates a new SendDataToDeviceParams object
// with the default values initialized.
func NewSendDataToDeviceParams() *SendDataToDeviceParams {
	var ()
	return &SendDataToDeviceParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewSendDataToDeviceParamsWithTimeout creates a new SendDataToDeviceParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewSendDataToDeviceParamsWithTimeout(timeout time.Duration) *SendDataToDeviceParams {
	var ()
	return &SendDataToDeviceParams{

		timeout: timeout,
	}
}

// NewSendDataToDeviceParamsWithContext creates a new SendDataToDeviceParams object
// with the default values initialized, and the ability to set a context for a request
func NewSendDataToDeviceParamsWithContext(ctx context.Context) *SendDataToDeviceParams {
	var ()
	return &SendDataToDeviceParams{

		Context: ctx,
	}
}

// NewSendDataToDeviceParamsWithHTTPClient creates a new SendDataToDeviceParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewSendDataToDeviceParamsWithHTTPClient(client *http.Client) *SendDataToDeviceParams {
	var ()
	return &SendDataToDeviceParams{
		HTTPClient: client,
	}
}

/*SendDataToDeviceParams contains all the parameters to send to the API endpoint
for the send data to device operation typically these are written to a http.Request
*/
type SendDataToDeviceParams struct {

	/*DeviceData
	  data to device

	*/
	DeviceData *models.DeviceData
	/*DeviceID
	  device identifer

	*/
	DeviceID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the send data to device params
func (o *SendDataToDeviceParams) WithTimeout(timeout time.Duration) *SendDataToDeviceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the send data to device params
func (o *SendDataToDeviceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the send data to device params
func (o *SendDataToDeviceParams) WithContext(ctx context.Context) *SendDataToDeviceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the send data to device params
func (o *SendDataToDeviceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the send data to device params
func (o *SendDataToDeviceParams) WithHTTPClient(client *http.Client) *SendDataToDeviceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the send data to device params
func (o *SendDataToDeviceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDeviceData adds the deviceData to the send data to device params
func (o *SendDataToDeviceParams) WithDeviceData(deviceData *models.DeviceData) *SendDataToDeviceParams {
	o.SetDeviceData(deviceData)
	return o
}

// SetDeviceData adds the deviceData to the send data to device params
func (o *SendDataToDeviceParams) SetDeviceData(deviceData *models.DeviceData) {
	o.DeviceData = deviceData
}

// WithDeviceID adds the deviceID to the send data to device params
func (o *SendDataToDeviceParams) WithDeviceID(deviceID string) *SendDataToDeviceParams {
	o.SetDeviceID(deviceID)
	return o
}

// SetDeviceID adds the deviceId to the send data to device params
func (o *SendDataToDeviceParams) SetDeviceID(deviceID string) {
	o.DeviceID = deviceID
}

// WriteToRequest writes these params to a swagger request
func (o *SendDataToDeviceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.DeviceData != nil {
		if err := r.SetBodyParam(o.DeviceData); err != nil {
			return err
		}
	}

	// path param deviceId
	if err := r.SetPathParam("deviceId", o.DeviceID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
