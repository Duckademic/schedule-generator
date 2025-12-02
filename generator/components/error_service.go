package components

import (
	"fmt"
	"strings"
)

// GeneratorComponentErrorTypes defines categories of generator errors.
type GeneratorComponentErrorTypes int

const (
	SetDayTypeErrorType GeneratorComponentErrorTypes = iota
	BoneWeekErrorType
	MissingLessonsAdderErrorType
)

// GeneratorComponentError represents a typed generator component error.
type GeneratorComponentError interface {
	error                                         // Basic interface for errors
	GetTypeOfError() GeneratorComponentErrorTypes // Each error generator must have a category
}

// ErrorService aggregates and manages errors produced by generator components.
type ErrorService interface {
	error                             // Implements the error interface; represents the final accumulated error.
	AddError(GeneratorComponentError) // Add error to collection. The service automatically handles ordering or deduplication.
	IsClear() bool                    // Returns true if no errors have been collected.
}

// NewErrorService creates new ErrorService instance
func NewErrorService() ErrorService {
	return &errorService{errorMap: map[GeneratorComponentErrorTypes][]error{}}
}

type errorService struct {
	errorMap map[GeneratorComponentErrorTypes][]error
}

func (ec *errorService) AddError(err GeneratorComponentError) {
	errorType := err.GetTypeOfError()
	if _, ok := ec.errorMap[errorType]; !ok {
		ec.errorMap[errorType] = make([]error, 0, 1)
	}
	ec.errorMap[errorType] = append(ec.errorMap[errorType], err)
}

func (ec *errorService) IsClear() bool {
	return len(ec.errorMap) == 0
}

func (ec *errorService) Error() string {
	if len(ec.errorMap) == 0 {
		return ""
	}

	var b strings.Builder
	for key, errs := range ec.errorMap {
		b.WriteString(fmt.Sprintf("%d:\n", key))
		for _, err := range errs {
			b.WriteString(fmt.Sprintf("- %s\n", err.Error()))
		}
		b.WriteString("\n")
	}

	return b.String()
}
