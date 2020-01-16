// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/cloustone/pandas/models"
)

// GetProjectDevicesOKCode is the HTTP code returned for type GetProjectDevicesOK
const GetProjectDevicesOKCode int = 200

/*GetProjectDevicesOK successful operation

swagger:response getProjectDevicesOK
*/
type GetProjectDevicesOK struct {

	/*
	  In: Body
	*/
	Payload []models.Device `json:"body,omitempty"`
}

// NewGetProjectDevicesOK creates GetProjectDevicesOK with default headers values
func NewGetProjectDevicesOK() *GetProjectDevicesOK {

	return &GetProjectDevicesOK{}
}

// WithPayload adds the payload to the get project devices o k response
func (o *GetProjectDevicesOK) WithPayload(payload []models.Device) *GetProjectDevicesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project devices o k response
func (o *GetProjectDevicesOK) SetPayload(payload []models.Device) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectDevicesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]models.Device, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetProjectDevicesNotFoundCode is the HTTP code returned for type GetProjectDevicesNotFound
const GetProjectDevicesNotFoundCode int = 404

/*GetProjectDevicesNotFound project not found

swagger:response getProjectDevicesNotFound
*/
type GetProjectDevicesNotFound struct {
}

// NewGetProjectDevicesNotFound creates GetProjectDevicesNotFound with default headers values
func NewGetProjectDevicesNotFound() *GetProjectDevicesNotFound {

	return &GetProjectDevicesNotFound{}
}

// WriteResponse to the client
func (o *GetProjectDevicesNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}