// Code generated by go-swagger; DO NOT EDIT.

package deployment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/cloustone/pandas/apimachinery/models"
)

// CreateDeploymentOKCode is the HTTP code returned for type CreateDeploymentOK
const CreateDeploymentOKCode int = 200

/*CreateDeploymentOK Successfully operation

swagger:response createDeploymentOK
*/
type CreateDeploymentOK struct {

	/*
	  In: Body
	*/
	Payload *models.Deployment `json:"body,omitempty"`
}

// NewCreateDeploymentOK creates CreateDeploymentOK with default headers values
func NewCreateDeploymentOK() *CreateDeploymentOK {

	return &CreateDeploymentOK{}
}

// WithPayload adds the payload to the create deployment o k response
func (o *CreateDeploymentOK) WithPayload(payload *models.Deployment) *CreateDeploymentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create deployment o k response
func (o *CreateDeploymentOK) SetPayload(payload *models.Deployment) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateDeploymentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateDeploymentBadRequestCode is the HTTP code returned for type CreateDeploymentBadRequest
const CreateDeploymentBadRequestCode int = 400

/*CreateDeploymentBadRequest Bad request

swagger:response createDeploymentBadRequest
*/
type CreateDeploymentBadRequest struct {
}

// NewCreateDeploymentBadRequest creates CreateDeploymentBadRequest with default headers values
func NewCreateDeploymentBadRequest() *CreateDeploymentBadRequest {

	return &CreateDeploymentBadRequest{}
}

// WriteResponse to the client
func (o *CreateDeploymentBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// CreateDeploymentMethodNotAllowedCode is the HTTP code returned for type CreateDeploymentMethodNotAllowed
const CreateDeploymentMethodNotAllowedCode int = 405

/*CreateDeploymentMethodNotAllowed Invalid condition

swagger:response createDeploymentMethodNotAllowed
*/
type CreateDeploymentMethodNotAllowed struct {
}

// NewCreateDeploymentMethodNotAllowed creates CreateDeploymentMethodNotAllowed with default headers values
func NewCreateDeploymentMethodNotAllowed() *CreateDeploymentMethodNotAllowed {

	return &CreateDeploymentMethodNotAllowed{}
}

// WriteResponse to the client
func (o *CreateDeploymentMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(405)
}
