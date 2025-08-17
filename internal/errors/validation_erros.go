package errors

import (
	"fmt"

)

type ValidationError struct{
	Field string
	Message string
}

func (e *ValidationError) Error()string{
	return (fmt.Sprintf("validation failed for field '%s':%s",e.Field,e.Message))
}
func NewValidationError(field,message string) *ValidationError{
	return &ValidationError{
		Field: field,
		Message: message,
	}
}

type DuplicateError struct{
	Resource string
	Value string
}
func(e *DuplicateError) Error() string{
	return (fmt.Sprintf("%s already exists: %s",e.Resource,e.Value))
}
func NewDuplicateError(resource,value string)*DuplicateError{
	return &DuplicateError{
		Resource: resource,
		Value: value,
	}
}

type NotFoundError struct{
	ID int
	Resource string
}

func (e *NotFoundError) Error()string{
	return (fmt.Sprintf("%s not found with id %d",e.Resource,e.ID))
}

func NewNotFoundError(id int,resource string) *NotFoundError{
	return &NotFoundError{
		ID: id,
		Resource: resource,
	}
}