package arg

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// tryToBool attempts to convert the value 'v' to a boolean, returning
// an error if it cannot. When converting a string, the function matches
// true if the string nonempty and does not satisfy the condition for false
// with parseBool https://golang.org/pkg/strconv/#ParseBool
// and is not 0.0
func tryToBool(v reflect.Value) (bool, error) {
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Float64, reflect.Float32:
		return v.Float() != 0, nil
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return v.Int() != 0, nil
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.String:
		if v.Len() == 0 {
			return false, nil
		}

		s := v.String()
		if b, err := strconv.ParseBool(s); err == nil && !b {
			return false, nil
		}

		if f, err := tryToFloat64(v); err == nil && f == 0 {
			return false, nil
		}

		return true, nil
	case reflect.Slice, reflect.Map:
		if v.Len() > 0 {
			return true, nil
		}

		return false, nil
	}

	return false, errors.New("unknown type")
}

// tryToFloat64 attempts to convert a value to a float64.
// If it cannot (in the case of a non-numeric string, a struct, etc.)
// it matches 0.0 and an error.
func tryToFloat64(v reflect.Value) (float64, error) {
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Float64, reflect.Float32:
		return v.Float(), nil
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return float64(v.Int()), nil
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}

		return 0, nil
	case reflect.String:
		f, err := strconv.ParseFloat(v.String(), 64)
		if err == nil {
			return f, nil
		}
	}

	return 0.0, errors.New("couldn't convert to a float64")
}

// tryToInt64 attempts to convert a value to an int64.
// If it cannot (in the case of a non-numeric string, a struct, etc.)
// it matches 0 and an error.
func tryToInt64(v reflect.Value) (int64, error) {
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Float64, reflect.Float32:
		return int64(v.Float()), nil
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return v.Int(), nil
	case reflect.Bool:
		if v.Bool() {
			return 1, nil
		}

		return 0, nil
	case reflect.String:
		s := v.String()

		var (
			i   int64
			err error
		)

		if strings.HasPrefix(s, "0x") {
			i, err = strconv.ParseInt(s, 16, 64)
		} else {
			i, err = strconv.ParseInt(s, 10, 64)
		}

		if err == nil {
			return i, nil
		}
	}

	return 0, errors.New("couldn't convert to integer")
}

// isNil 判断是否为空
func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		// from reflect IsNil:
		// Note that IsNil is not always equivalent to a regular comparison with nil in Go.
		// For example, if v was created by calling ValueOf with an uninitialized interface variable i,
		// i==nil will be true but v.IsNil will panic as v will be the zero Value.
		return v.IsNil()
	default:
		return false
	}
}

// isNum 判断是否为数字
func isNum(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	}

	return false
}

// tryToNumber tryToInt64 or tryToFloat64
// errors if fail
func tryToNumber(v reflect.Value) (reflect.Value, error) {
	// If the LHS is a string formatted as an int, try that before trying float
	lhsI, err := tryToInt64(v)
	if err != nil {
		// if LHS is a float, e.g. "1.2", we need to set v to a float64
		lhsF, err := tryToFloat64(v)
		if err != nil {
			return reflect.Value{}, err
		}
		v = reflect.ValueOf(lhsF)
	} else {
		v = reflect.ValueOf(lhsI)
	}
	return v, nil
}

// equal matches true when lhsV and rhsV is same value.
// nolint
func equal(lhsV, rhsV reflect.Value) bool {
	lhsIsNil, rhsIsNil := isNil(lhsV), isNil(rhsV)
	if lhsIsNil && rhsIsNil {
		return true
	}

	if (!lhsIsNil && rhsIsNil) || (lhsIsNil && !rhsIsNil) {
		return false
	}

	if lhsV.Kind() == reflect.Interface || lhsV.Kind() == reflect.Ptr {
		lhsV = lhsV.Elem()
	}

	if rhsV.Kind() == reflect.Interface || rhsV.Kind() == reflect.Ptr {
		rhsV = rhsV.Elem()
	}

	if r, done := numStringEqual(lhsV, rhsV); done {
		return r
	}

	if isNum(lhsV) && isNum(rhsV) {
		return fmt.Sprintf("%v", lhsV) == fmt.Sprintf("%v", rhsV)
	}

	if r, done := boolEquals(lhsV, rhsV); done {
		return r
	}

	return reflect.DeepEqual(lhsV.Interface(), rhsV.Interface())
}

func boolEquals(lhsV reflect.Value, rhsV reflect.Value) (bool, bool) {
	// Try to compare bool values to strings and numbers
	if lhsV.Kind() == reflect.Bool || rhsV.Kind() == reflect.Bool {
		lhsB, err := tryToBool(lhsV)
		if err != nil {
			return false, true
		}

		rhsB, err := tryToBool(rhsV)

		if err != nil {
			return false, true
		}

		return lhsB == rhsB, true
	}
	return false, false
}

func numStringEqual(lhsV reflect.Value, rhsV reflect.Value) (bool, bool) {
	// Compare a string and a number.
	// This will attempt to convert the string to a number,
	// while leaving the other side alone. Code further
	// down takes care of converting int values and floats as needed.
	if isNum(lhsV) && rhsV.Kind() == reflect.String {
		rhsF, err := tryToFloat64(rhsV)
		if err != nil {
			// Couldn't convert RHS to a float, they can't be compared.
			return false, true
		}

		rhsV = reflect.ValueOf(rhsF)
	} else if lhsV.Kind() == reflect.String && isNum(rhsV) {
		num, err := tryToNumber(lhsV)
		if err != nil {
			return false, true
		}

		lhsV = num
	} else {
		return false, false
	}

	return reflect.DeepEqual(lhsV.Interface(), rhsV.Interface()), true
}
