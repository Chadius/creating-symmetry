package utility

import (
	"fmt"
	"gopkg.in/check.v1"
	"math"
	"reflect"
)

// NumericallyCloseEnough asserts the two given numbers are within tolerance.
type NumericallyCloseEnough struct {}

// Info returns information for the NumericallyCloseEnough
func (checker NumericallyCloseEnough) Info() *check.CheckerInfo {
	return &check.CheckerInfo{Name: "NumericallyCloseEnough", Params: []string{"obtained", "expected", "tolerance"}}
}

// Check to see if the obtained and expected are within tolerance.
func (checker NumericallyCloseEnough) Check(params []interface{}, names []string) (result bool, error string) {
	interfaceObtained := params[0]
	interfaceExpected := params[1]
	interfaceTolerance := params[2]

	var floatType = reflect.TypeOf(float64(0))

	v := reflect.ValueOf(interfaceObtained)
	v = reflect.Indirect(v)
	obtained := v.Convert(floatType).Float()
	v = reflect.ValueOf(interfaceExpected)
	v = reflect.Indirect(v)
	expected := v.Convert(floatType).Float()
	v = reflect.ValueOf(interfaceTolerance)
	v = reflect.Indirect(v)
	tolerance := v.Convert(floatType).Float()

	if math.Abs(obtained - expected) < tolerance {
		return true, ""
	}

	return false, fmt.Sprintf("%f and %f are outside of tolerance %f", obtained, expected, tolerance)
}
