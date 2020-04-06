// Code generated by go-swagger; DO NOT EDIT.

package rulechain

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/cloustone/pandas/apimachinery/models"
)

// UploadRuleChainOKCode is the HTTP code returned for type UploadRuleChainOK
const UploadRuleChainOKCode int = 200

/*UploadRuleChainOK excute successfully

swagger:response uploadRuleChainOK
*/
type UploadRuleChainOK struct {

	/*
	  In: Body
	*/
	Payload []*models.RuleChain `json:"body,omitempty"`
}

// NewUploadRuleChainOK creates UploadRuleChainOK with default headers values
func NewUploadRuleChainOK() *UploadRuleChainOK {

	return &UploadRuleChainOK{}
}

// WithPayload adds the payload to the upload rule chain o k response
func (o *UploadRuleChainOK) WithPayload(payload []*models.RuleChain) *UploadRuleChainOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the upload rule chain o k response
func (o *UploadRuleChainOK) SetPayload(payload []*models.RuleChain) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UploadRuleChainOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.RuleChain, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// UploadRuleChainInternalServerErrorCode is the HTTP code returned for type UploadRuleChainInternalServerError
const UploadRuleChainInternalServerErrorCode int = 500

/*UploadRuleChainInternalServerError Server internal error

swagger:response uploadRuleChainInternalServerError
*/
type UploadRuleChainInternalServerError struct {
}

// NewUploadRuleChainInternalServerError creates UploadRuleChainInternalServerError with default headers values
func NewUploadRuleChainInternalServerError() *UploadRuleChainInternalServerError {

	return &UploadRuleChainInternalServerError{}
}

// WriteResponse to the client
func (o *UploadRuleChainInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(500)
}
