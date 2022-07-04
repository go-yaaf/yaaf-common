// Copyright 2020. Motty Cohen.
//
// The recover utility helps to work with `recover` in bit pleasant way.

package utils

// All performs recover for all panics.
//
// Sample usage:
//		defer RecoverAll(func(err interface{}) {
//			fmt.Printf("got error: %s", err)
//		})
//
func RecoverAll(cb func(v any)) {
	r := recover()
	cb(r)
}

// One performs call cb function with recovered value
// in case when recovered value equals to e otherwise panic won't be recovered and will be propagated.
//
// Sample usage:
//		defer recover.One(ErrorUsernameBlank, func(err interface{}) {
// 			fmt.Printf("got error: %s", err)
//		})
//
func RecoverOne(e error, cb func(v any)) {
	r := recover()

	errors := []error{e}

	if inErrors(r, errors) {
		cb(r)
		return
	}
	panic(r)
}

// Any performs call cb function with recovered value
// in case when recovered value exists in slice errors.
//
// Sample usage:
//		defer recover.Any([]error{ErrorUsernameBlank, ErrorUsernameAlreadyTaken}, func(err interface{}) {
//  		fmt.Printf("got error: %s", err)
//		})
//
func RecoverAny(errors []error, cb func(v any)) {
	r := recover()

	if len(errors) == 0 || inErrors(r, errors) {
		cb(r)
		return
	}
	panic(r)
}

// Check if the given error included in the error list
func inErrors(e any, errors []error) bool {
	for _, err := range errors {
		if e == err {
			return true
		}
	}
	return false
}
