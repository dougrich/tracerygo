package tracerygo

import (
	"fmt"
)

type StreamWriteError struct {
	Underlying error
}

func (s StreamWriteError) Error() string {
	return fmt.Sprintf("Error writing to stream: %v", s.Underlying)
}

type LookupError struct {
	Name       string
	Underlying error
}

func (l LookupError) Error() string {
	return fmt.Sprintf("Error looking up name %s: %v", l.Name, l.Underlying)
}

type NameNotFoundError struct {
	Name string
}

func (n NameNotFoundError) Error() string {
	return fmt.Sprintf("'%s' was not found in the grammar or declared; ensure either provided to grammar or declared inline", n.Name)
}

type UnsupportedModifierError struct {
	Modifier  string
}

func (u UnsupportedModifierError) Error() string {
	return fmt.Sprintf("unsupported modifier '%s' found", u.Modifier)
}

type ExpectationError struct {
	Expected  string
	Found     string
}

func (e ExpectationError) Error() string {
	return fmt.Sprintf("expected to find %s but found %s", e.Expected, e.Found)
}

type UnmatchedSymbolError struct {
	Index         int
	ExpectedFirst string
	MissingSecond string
}

func (u UnmatchedSymbolError) Error() string {
	return fmt.Sprintf("expected to find a %s to pair with %s starting at %d; went unpaired", u.MissingSecond, u.ExpectedFirst, u.Index)
}

type FieldError struct {
	FieldName string
	Underlying error
}

func (f FieldError) Error() string {
	return fmt.Sprintf("in field '%s': %s", f.FieldName, f.Underlying)
}