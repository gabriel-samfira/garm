// Code generated by go-swagger; DO NOT EDIT.

package endpoints

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

	garm_params "github.com/cloudbase/garm/params"
)

// NewUpdateGithubEndpointParams creates a new UpdateGithubEndpointParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateGithubEndpointParams() *UpdateGithubEndpointParams {
	return &UpdateGithubEndpointParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateGithubEndpointParamsWithTimeout creates a new UpdateGithubEndpointParams object
// with the ability to set a timeout on a request.
func NewUpdateGithubEndpointParamsWithTimeout(timeout time.Duration) *UpdateGithubEndpointParams {
	return &UpdateGithubEndpointParams{
		timeout: timeout,
	}
}

// NewUpdateGithubEndpointParamsWithContext creates a new UpdateGithubEndpointParams object
// with the ability to set a context for a request.
func NewUpdateGithubEndpointParamsWithContext(ctx context.Context) *UpdateGithubEndpointParams {
	return &UpdateGithubEndpointParams{
		Context: ctx,
	}
}

// NewUpdateGithubEndpointParamsWithHTTPClient creates a new UpdateGithubEndpointParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateGithubEndpointParamsWithHTTPClient(client *http.Client) *UpdateGithubEndpointParams {
	return &UpdateGithubEndpointParams{
		HTTPClient: client,
	}
}

/*
UpdateGithubEndpointParams contains all the parameters to send to the API endpoint

	for the update github endpoint operation.

	Typically these are written to a http.Request.
*/
type UpdateGithubEndpointParams struct {

	/* Body.

	   Parameters used when updating a GitHub endpoint.
	*/
	Body garm_params.UpdateGithubEndpointParams

	/* Name.

	   The name of the GitHub endpoint.
	*/
	Name string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update github endpoint params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateGithubEndpointParams) WithDefaults() *UpdateGithubEndpointParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update github endpoint params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateGithubEndpointParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update github endpoint params
func (o *UpdateGithubEndpointParams) WithTimeout(timeout time.Duration) *UpdateGithubEndpointParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update github endpoint params
func (o *UpdateGithubEndpointParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update github endpoint params
func (o *UpdateGithubEndpointParams) WithContext(ctx context.Context) *UpdateGithubEndpointParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update github endpoint params
func (o *UpdateGithubEndpointParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update github endpoint params
func (o *UpdateGithubEndpointParams) WithHTTPClient(client *http.Client) *UpdateGithubEndpointParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update github endpoint params
func (o *UpdateGithubEndpointParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update github endpoint params
func (o *UpdateGithubEndpointParams) WithBody(body garm_params.UpdateGithubEndpointParams) *UpdateGithubEndpointParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update github endpoint params
func (o *UpdateGithubEndpointParams) SetBody(body garm_params.UpdateGithubEndpointParams) {
	o.Body = body
}

// WithName adds the name to the update github endpoint params
func (o *UpdateGithubEndpointParams) WithName(name string) *UpdateGithubEndpointParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the update github endpoint params
func (o *UpdateGithubEndpointParams) SetName(name string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateGithubEndpointParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}