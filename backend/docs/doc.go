// Package docs provides the Swagger/OpenAPI specification metadata for the Nexus Tasks API.
//
// # Nexus Tasks REST API
//
// Nexus Tasks is a collaborative project and task management platform.
// This document describes the full REST API surface of the backend service.
//
// ## Authentication
//
// Most endpoints require an authenticated session established via the `/api/v2/auth/login`
// or `/api/v2/auth/register` endpoints. Authentication state is maintained via an
// HTTP-only session cookie (`nexus_user`). Alternatively, machine clients may
// authenticate using an API Key passed in the `X-API-Key` request header.
//
// ## Base URL
//
// All versioned resources are prefixed with `/api/v2`.
//
// Schemes: http, https
// Host: localhost:8000
// BasePath: /api/v2
// Version: 2.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// SecurityDefinitions:
//
//	cookieAuth:
//	  type: apiKey
//	  in: header
//	  name: Cookie
//	  description: >
//	    HTTP-only session cookie set by the login/register endpoints.
//	    The cookie name is `nexus_user` and contains a signed JWT.
//
//	apiKeyAuth:
//	  type: apiKey
//	  in: header
//	  name: X-API-Key
//	  description: >
//	    A long-lived API key generated via the `/api/v2/auth/api-keys` endpoint.
//	    Intended for machine-to-machine access.
//
// Security:
// - cookieAuth: []
// - apiKeyAuth: []
//
// swagger:meta
package docs
