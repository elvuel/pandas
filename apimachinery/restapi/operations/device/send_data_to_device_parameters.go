// Code generated by go-swagger; DO NOT EDIT.

package device

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/cloustone/pandas/models"
)

// NewSendDataToDeviceParams creates a new SendDataToDeviceParams object
// no default values defined in spec.
func NewSendDataToDeviceParams() SendDataToDeviceParams {

	return SendDataToDeviceParams{}
}

// SendDataToDeviceParams contains all the bound params for the send data to device operation
// typically these are obtained from a http.Request
//
// swagger:parameters sendDataToDevice
type SendDataToDeviceParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*data to device
	  Required: true
	  In: body
	*/
	DeviceData *models.DeviceData
	/*device identifer
	  Required: true
	  In: path
	*/
	DeviceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSendDataToDeviceParams() beforehand.
func (o *SendDataToDeviceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.DeviceData
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("deviceData", "body"))
			} else {
				res = append(res, errors.NewParseError("deviceData", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.DeviceData = &body
			}
		}
	} else {
		res = append(res, errors.Required("deviceData", "body"))
	}
	rDeviceID, rhkDeviceID, _ := route.Params.GetOK("deviceId")
	if err := o.bindDeviceID(rDeviceID, rhkDeviceID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindDeviceID binds and validates parameter DeviceID from path.
func (o *SendDataToDeviceParams) bindDeviceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.DeviceID = raw

	return nil
}
