package components

import (
	"fmt"
	"strings"
)

// ScheduleParameter defines a single component of the ScheduleFault.
// Each implementation provides its own fault calculation logic.
type ScheduleParameter interface {
	Fault() float64       // Calculates fault value based on the argument's internal calculation formula.
	GetArguments() string // Returns formatted, human-readable representation of the arguments
}

// NewSimpleScheduleParameter a new simpleScheduleParameter instance.
// The resulting fault value is computed as value (v) multiplied by its
// fault weight factor (f).
func NewSimpleScheduleParameter(v, f float64) ScheduleParameter {
	return &simpleScheduleParameter{value: v, fault: f}
}

type simpleScheduleParameter struct {
	fault float64
	value float64
}

func (ssp *simpleScheduleParameter) Fault() float64 {
	return ssp.fault * ssp.value
}
func (ssp *simpleScheduleParameter) GetArguments() string {
	return fmt.Sprintf("value: %f, fault %f", ssp.value, ssp.fault)
}

// ScheduleFault represents fault for generated schedule.
type ScheduleFault interface {
	AddParameter(name string, parameter ScheduleParameter) // Adds new ScheduleParameter or replace it with a new one.
	Fault() float64                                        // Returns sum of the ScheduleParameter Faults.
	GetParameters() string                                 // Returns formatted, human-readable representation of the parameters.
}

// NewScheduleFault creates a new ScheduleFault instance.
func NewScheduleFault() ScheduleFault {
	return &scheduleFault{parameters: map[string]ScheduleParameter{}}
}

type scheduleFault struct {
	parameters map[string]ScheduleParameter
}

func (sf *scheduleFault) AddParameter(name string, parameter ScheduleParameter) {
	sf.parameters[name] = parameter
}
func (sf *scheduleFault) Fault() (res float64) {
	for _, fault := range sf.parameters {
		res += fault.Fault()
	}
	return
}
func (sf *scheduleFault) GetParameters() string {
	if len(sf.parameters) == 0 {
		return ""
	}

	var b strings.Builder
	for key, value := range sf.parameters {
		if value.Fault() != 0 {
			b.WriteString(fmt.Sprintf("%s: %s, fault: %f \n", key, value.GetArguments(), value.Fault()))
		}
	}

	return b.String()
}
