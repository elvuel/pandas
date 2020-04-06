// Code generated by go-swagger; DO NOT EDIT.

package rulechain

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/cloustone/pandas/apimachinery/models"
)

// UploadRuleChainReader is a Reader for the UploadRuleChain structure.
type UploadRuleChainReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UploadRuleChainReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUploadRuleChainOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewUploadRuleChainInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUploadRuleChainOK creates a UploadRuleChainOK with default headers values
func NewUploadRuleChainOK() *UploadRuleChainOK {
	return &UploadRuleChainOK{}
}

/*UploadRuleChainOK handles this case with default header values.

excute successfully
*/
type UploadRuleChainOK struct {
	Payload []*models.RuleChain
}

func (o *UploadRuleChainOK) Error() string {
	return fmt.Sprintf("[POST /rulechains/{ruleChainId}/upload][%d] uploadRuleChainOK  %+v", 200, o.Payload)
}

func (o *UploadRuleChainOK) GetPayload() []*models.RuleChain {
	return o.Payload
}

func (o *UploadRuleChainOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUploadRuleChainInternalServerError creates a UploadRuleChainInternalServerError with default headers values
func NewUploadRuleChainInternalServerError() *UploadRuleChainInternalServerError {
	return &UploadRuleChainInternalServerError{}
}

/*UploadRuleChainInternalServerError handles this case with default header values.

Server internal error
*/
type UploadRuleChainInternalServerError struct {
}

func (o *UploadRuleChainInternalServerError) Error() string {
	return fmt.Sprintf("[POST /rulechains/{ruleChainId}/upload][%d] uploadRuleChainInternalServerError ", 500)
}

func (o *UploadRuleChainInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
