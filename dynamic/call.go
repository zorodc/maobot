/* This package includes runtime-checked function calls. */
package dynamic

import "fmt"
import "reflect"
import "errors"

type RtFunc interface{}

// Returns an instance of DynFunc from a function of any signature,
// verified to be a function.
// It is unreasonable to pass in a non-function, and so this causes a panic.
func New(function interface{}) RtFunc {
	if reflect.TypeOf(function).Kind() != reflect.Func {
		panic("Parameter `function` not a function")
	}
	return RtFunc(function)
}


// Determine whether a prospective parameter `param` satisfies the type
// specification of the formal parameter `formalParam`.
func Satisfactory(param, formalParam reflect.Type) bool {
	// A parameter satisfies the requirements of a formal parameter if
	// 1. The types are the same, or
	// 2. The formal param is an interface and the parameter implemements it.
	return param == formalParam ||
		(formalParam.Kind() == reflect.Interface &&
		param.Implements(formalParam))
}

// Executes a dynamic function call.
// Takes and returns a slice of interface{}s representing the function,
// parameters, and return values. This is so ordinary type assertions can be
// used in client code, rather than explicit use of the reflect package.
// In the implementation, this involves copious conversion between
// reflect.Value objects and interface{} objects, but this is opaque to clients.
func Call(function RtFunc,
	parameters ...interface{}) ([]interface{}, error) {
	funcType := reflect.TypeOf(function)

	var nArgs int = funcType.NumIn()
	if funcType.IsVariadic() {
		// The final parameter for a variadic function is visible as a slice
		// of the type taken in a variadic manner, so the params checked normally
		// for variadic functions are (0, ft.NumIn() - 1).
		nArgs -= 1 
	}

	// Detect arity mismatches.
	// calls to variadic functions must at least supply f.NumIn()-1 parameters,
	// calls to ordinary functions must supply exactly f.NumIn() parameters.
	if len(parameters) < nArgs || (
		!funcType.IsVariadic() && len(parameters) > nArgs) {
		return nil, errors.New("Arity mismatch") 
	}

	// Ensure for non-variadic parameters type equivalence.
	for i := 0; i < nArgs; i++ {
		// Check the param type implements the interface required.
		paramType := reflect.TypeOf(parameters[i])
		if !Satisfactory(paramType, funcType.In(i)) {
			return nil,
			errors.New(fmt.Sprintf("Type mismatch on parameter %d: %s != %s",
				i, paramType, funcType.In(i)));
		}
	}

	// If the function is variadic, then this will be run.
	// Ensure type equivalence for variadic parameters.
	for i := nArgs; i < len(parameters); i++ {
		// Check the param type implements the interface required.
		paramType := reflect.TypeOf(parameters[i])
		if !Satisfactory(paramType, funcType.In(nArgs).Elem()) {
			return nil,
			errors.New(fmt.Sprintf("Type mismatch on variadic param %d", i))
		}
	}

	// At this point, the crash-avoiding checks are over.

	/***** Generate parameter list *****/
	// Value.Call() requires []Value, not []interface{}
	var callValues []reflect.Value
	for _, e := range parameters {
		callValues = append(callValues, reflect.ValueOf(e))
	}

	/***** Execute a dynamic call *****/
	returnValues := reflect.ValueOf(function).Call(callValues)
	var returnLst []interface{} // Return a []interface{} instead of a []Value
	for _, e := range returnValues {
		returnLst = append(returnLst, e.Interface())
	}

	return returnLst, nil
}

