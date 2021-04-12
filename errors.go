package tracerygo

import (
	"errors"
	"fmt"
)

var (
	// This is returned when a modifier doesn't write as many bytes as it expected to write
	ErrUnexpectedNumberOfBytesWritten = errors.New("Unexpected number of bytes written")
)

// This error wraps an error occuring in the underlying stream when writing
type ErrorStreamWrite struct {
	Underlying error
}

// Serializes the error message
func (s ErrorStreamWrite) Error() string {
	return fmt.Sprintf("Error writing to stream: %v", s.Underlying)
}

// This error wraps an error occuring in the underlying lookup for a name
type ErrorLookup struct {
	Name       string
	Underlying error
}

// Serializes the error message
func (l ErrorLookup) Error() string {
	return fmt.Sprintf("Error looking up name %s: %v", l.Name, l.Underlying)
}

// This error occurs if a name should be substituted (e.g. '#mention#') but no value can be found, either by lookup or by definition
type ErrorNameNotFound struct {
	Name string
}

// Serializes the error message
func (n ErrorNameNotFound) Error() string {
	return fmt.Sprintf("'%s' was not found in the grammar or declared; ensure either provided to grammar or declared inline", n.Name)
}

// This error occurs if a modifier on a substitution (e.g. '#mention.ed#') doesn't exist
type ErrorUnsupportedModifier struct {
	Modifier string
}

// Serializes the error message
func (u ErrorUnsupportedModifier) Error() string {
	return fmt.Sprintf("unsupported modifier '%s' found", u.Modifier)
}

// This error occurs during parsing when an expected type assertion fails; expected a string but got something else, expected an array but got something else, etc.
type ErrorExpectationFailed struct {
	Expected string
	Found    string
}

// Serializes the error message
func (e ErrorExpectationFailed) Error() string {
	return fmt.Sprintf("expected to find %s but found %s", e.Expected, e.Found)
}

// This error occurs during parsing when a symbol is unmatched; for example '#a##' will generate a substitution of 'a' and then fail because there is no closing tag
type ErrorUnmatchedSymbol struct {
	Index         int
	ExpectedFirst string
	MissingSecond string
}

// Serializes the error message
func (u ErrorUnmatchedSymbol) Error() string {
	return fmt.Sprintf("expected to find a %s to pair with %s starting at %d; went unpaired", u.MissingSecond, u.ExpectedFirst, u.Index)
}

// This error wraps an error that occurs when parsing or evaluting a field, and just decorates the field name for readability
type ErrorInField struct {
	FieldName  string
	Underlying error
}

// Serializes the error message
func (f ErrorInField) Error() string {
	return fmt.Sprintf("in field '%s': %s", f.FieldName, f.Underlying)
}
