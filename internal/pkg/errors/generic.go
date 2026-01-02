package errors

import "errors"

// MarshalFail is returned when the request body could not be marshalled
var MarshalFail = errors.New("encountered marshal error")

// UnmarshalFail is returned when the response body could not be unmarshalled
var UnmarshalFail = errors.New("encountered unmarshal error")

// FileUnreachable represents an error indicating that a file could not be reached or accessed.
var FileUnreachable = errors.New("file is unreachable")

// FileUnreadable represents an error indicating that a file could not be read.
var FileUnreadable = errors.New("file is unreadable")

// FileUnwriteable represents an error indicating that a file could not be created or written in.
var FileUnwriteable = errors.New("file is unwriteable")

// DirUnreachable represents an error indicating that a directory could not be reached or accessed.
var DirUnreachable = errors.New("directory is unreachable")

// DirUnreadable represents an error indicating that a directory could not be read.
var DirUnreadable = errors.New("directory is unreadable")

// DirUnwriteable represents an error indicating that a directory could not be created or written in.
var DirUnwriteable = errors.New("directory is unwriteable")
