// Copyright 2022. Motty Cohen
//
// Utility functions to extract data from HTTP REST request
//
package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/mottyc/yaaf-common/utils"
	"github.com/mottyc/yaaf-common/utils/collections"
)

// ResolveRemoteIpFromHttpRequest extracts remote ip from HTTP header X-Forwarded-For
func ResolveRemoteIpFromHttpRequest(r *http.Request) (ip string) {

	ip = r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return
}

// GetIntParamValue extract parameter value from query string as int
func GetIntParamValue(r *http.Request, paramName string, defaultValue int) (res int) {

	var (
		err error
	)
	res = defaultValue

	if tmp := r.URL.Query()[paramName]; len(tmp) > 0 {
		if res, err = strconv.Atoi(tmp[0]); err != nil {
			res = defaultValue
		}
	}

	return res
}

// GetFloatParamValue extract parameter value from query string as float64
func GetFloatParamValue(r *http.Request, paramName string, defaultValue float64) (res float64) {

	var (
		err error
	)
	res = defaultValue

	if tmp := r.URL.Query()[paramName]; len(tmp) > 0 {
		if res, err = strconv.ParseFloat(tmp[0], 32); err != nil {
			res = defaultValue
		}
	}

	return res
}

// GetIntParamArray extract parameter array values from query string
// to support multiple values query string e.g. https//some/domain?id=1&id=2&id=3
func GetIntParamArray(r *http.Request, paramName string) (res []int) {
	params := r.URL.Query()[paramName]
	for _, p := range params {
		arr := strings.Split(p, ",")

		for _, str := range arr {
			if num, err := strconv.Atoi(str); err == nil {
				res = append(res, num)
			}
		}
	}
	return
}

// GetUint64ParamValue extract parameter value from query string as unsigned-long
func GetUint64ParamValue(r *http.Request, paramName string, defaultValue uint64) (res uint64) {

	var err error

	res = defaultValue

	if tmp := r.URL.Query()[paramName]; len(tmp) > 0 {
		if res, err = strconv.ParseUint(tmp[0], 10, 64); err != nil {
			res = defaultValue
		}
	}
	return res
}

// GetInt64ParamValue extract parameter value from query string as long
func GetInt64ParamValue(r *http.Request, paramName string, defaultValue int64) (res int64) {

	var err error

	res = defaultValue

	if tmp := r.URL.Query()[paramName]; len(tmp) > 0 {
		if res, err = strconv.ParseInt(tmp[0], 10, 64); err != nil {
			res = defaultValue
		}
	}
	return res
}

// GetUInt64ParamArray extract parameter array values from query string
// This supports multiple values query string e.g. https//some/domain?id=1&id=2&id=3
func GetUInt64ParamArray(r *http.Request, paramName string) (res []uint64) {
	params := r.URL.Query()[paramName]
	for _, p := range params {
		arr := strings.Split(p, ",")

		for _, str := range arr {
			if num, err := strconv.ParseUint(str, 10, 64); err == nil {
				res = append(res, num)
			}
		}
	}
	return
}

// GetStringParamValue extract parameter value from query string as string
func GetStringParamValue(r *http.Request, paramName string, defaultValue string) (res string) {

	res = defaultValue

	if tmp := r.URL.Query()[paramName]; len(tmp) > 0 {
		res = tmp[0]
	}
	return
}

// GetStringParamArray extracts parameter array values from query string
// to support multiple values query string e.g. https//some/domain?id=1&id=2&id=3
func GetStringParamArray(r *http.Request, paramName string) (res []string) {
	params := r.URL.Query()[paramName]
	for _, p := range params {
		if len(p) > 0 {
			arr := strings.Split(p, ",")
			res = append(res, arr...)
		}
	}
	return
}

// GetBoolParamValue extracts parameter value from query string as bool
func GetBoolParamValue(r *http.Request, paramName string, defaultValue bool) (res bool) {

	res = defaultValue

	if tmp := r.URL.Query()[paramName]; len(tmp) > 0 {
		if t, err := strconv.ParseBool(tmp[0]); err == nil {
			res = t
		}
	}
	return
}

// GetIntParamValueFromPath extracts parameter value from path params as int
func GetIntParamValueFromPath(r *http.Request, paramName string, defaultValue int) int {

	var (
		res int
		err error
	)
	if tmp := mux.Vars(r)[paramName]; len(tmp) > 0 {
		if res, err = strconv.Atoi(tmp); err != nil {
			res = defaultValue
		}
	}

	return res
}

// GetStringParamValueFromPath extracts parameter value from path params as string
func GetStringParamValueFromPath(r *http.Request, paramName string, defaultValue string) (res string) {

	res = defaultValue

	if tmp := mux.Vars(r)[paramName]; len(tmp) > 0 {
		res = tmp
	}
	return
}

// GetEnumParamValue extracts parameter value from query string as enum
// Enum value can be passed as it's int or string value
func GetEnumParamValue(r *http.Request, paramName string, enum interface{}, defaultValue int) (res int) {

	var err error

	res = defaultValue

	if tmp := r.URL.Query()[paramName]; len(tmp) > 0 {
		if len(tmp) == 0 {
			return res
		}
		if res, err = strconv.Atoi(tmp[0]); err != nil {
			// try enum string
			if n, e := getEnumValueFromName(enum, tmp[0]); e != nil {
				res = defaultValue
			} else {
				res = n
			}
		}
	}
	return res
}

// GetEnumParamArray extracts parameter array values from query string as enums
// to support multiple values query string e.g. https//some/domain?id=1&id=2&id=3
func GetEnumParamArray(r *http.Request, paramName string, enum interface{}) (res []int) {
	params := r.URL.Query()[paramName]
	for _, p := range params {
		arr := strings.Split(p, ",")

		for _, str := range arr {
			if len(str) == 0 {
				continue
			}
			if num, err := strconv.Atoi(str); err != nil {
				// try enum string
				if n, e := getEnumValueFromName(enum, str); e == nil {
					res = append(res, n)
				}
			} else {
				res = append(res, num)
			}
		}
	}
	return
}

// getEnumValueFromName gets enum int value from enum name
func getEnumValueFromName(enum any, name string) (result int, err error) {

	// Handle any Panic error
	defer utils.RecoverAll(func(e interface{}) {
		err = fmt.Errorf("PANIC")
		return
	})

	ref := reflect.ValueOf(enum)
	typ := ref.Type()

	for i := 0; i < ref.NumField(); i++ {

		s := typ.Field(i).Name
		res := int(ref.Field(i).Int())

		if s == name {
			return res, nil
		}
	}

	return 0, fmt.Errorf("not found")
}

// getEnumValuesFromNames gets enum int value from enum name
func getEnumValuesFromNames(enum any, values []string) []int {

	// Handle any Panic error
	defer utils.RecoverAll(func(err any) {
		return
	})

	result := make([]int, 0)

	ref := reflect.ValueOf(enum)
	typ := ref.Type()

	for i := 0; i < ref.NumField(); i++ {

		s := typ.Field(i).Name
		res := int(ref.Field(i).Int())

		if collections.Include(values, s) {
			result = append(result, res)
		}
	}

	return result
}
