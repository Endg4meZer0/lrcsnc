package errors

import "errors"

// NotFound is returned when the requested song's lyrics are not found
var NotFound = errors.New("the requested song's lyrics were not found")

// ServerTimeout is returned when the server gets timed out before it can
// send anything
var ServerTimeout = errors.New("connection timed out")

// ServerError is returned when the server returns any other error code
// than 404 and timeout
var ServerError = errors.New("a server error occurred")

// BodyReadFail is returned when the body of the response could not be read
var BodyReadFail = errors.New("failed to read the response body")
