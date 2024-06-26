// Code generated by go-swagger; DO NOT EDIT.

package enterprises

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewDeleteEnterpriseParams creates a new DeleteEnterpriseParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteEnterpriseParams() *DeleteEnterpriseParams {
	return &DeleteEnterpriseParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteEnterpriseParamsWithTimeout creates a new DeleteEnterpriseParams object
// with the ability to set a timeout on a request.
func NewDeleteEnterpriseParamsWithTimeout(timeout time.Duration) *DeleteEnterpriseParams {
	return &DeleteEnterpriseParams{
		timeout: timeout,
	}
}

// NewDeleteEnterpriseParamsWithContext creates a new DeleteEnterpriseParams object
// with the ability to set a context for a request.
func NewDeleteEnterpriseParamsWithContext(ctx context.Context) *DeleteEnterpriseParams {
	return &DeleteEnterpriseParams{
		Context: ctx,
	}
}

// NewDeleteEnterpriseParamsWithHTTPClient creates a new DeleteEnterpriseParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteEnterpriseParamsWithHTTPClient(client *http.Client) *DeleteEnterpriseParams {
	return &DeleteEnterpriseParams{
		HTTPClient: client,
	}
}

/*
DeleteEnterpriseParams contains all the parameters to send to the API endpoint

	for the delete enterprise operation.

	Typically these are written to a http.Request.
*/
type DeleteEnterpriseParams struct {

	/* EnterpriseID.

	   ID of the enterprise to delete.
	*/
	EnterpriseID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete enterprise params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteEnterpriseParams) WithDefaults() *DeleteEnterpriseParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete enterprise params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteEnterpriseParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete enterprise params
func (o *DeleteEnterpriseParams) WithTimeout(timeout time.Duration) *DeleteEnterpriseParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete enterprise params
func (o *DeleteEnterpriseParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete enterprise params
func (o *DeleteEnterpriseParams) WithContext(ctx context.Context) *DeleteEnterpriseParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete enterprise params
func (o *DeleteEnterpriseParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete enterprise params
func (o *DeleteEnterpriseParams) WithHTTPClient(client *http.Client) *DeleteEnterpriseParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete enterprise params
func (o *DeleteEnterpriseParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithEnterpriseID adds the enterpriseID to the delete enterprise params
func (o *DeleteEnterpriseParams) WithEnterpriseID(enterpriseID string) *DeleteEnterpriseParams {
	o.SetEnterpriseID(enterpriseID)
	return o
}

// SetEnterpriseID adds the enterpriseId to the delete enterprise params
func (o *DeleteEnterpriseParams) SetEnterpriseID(enterpriseID string) {
	o.EnterpriseID = enterpriseID
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteEnterpriseParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param enterpriseID
	if err := r.SetPathParam("enterpriseID", o.EnterpriseID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
