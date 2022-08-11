// Copyright 2022. Motty Cohen
//
// rest package definitions
//
package rest

import (
	"net/http"
	"strings"
)

func DefaultRestHandlerWrapperFunc(entry RestEntry) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		entry.Handler(NewRequestWithToken(rw, r, ""))
	}
}

// region Sortable collection of rest entries --------------------------------------------------------------------------
type RestEntries []*RestEntry

func (a RestEntries) Len() int {
	return len(a)
}

func (a RestEntries) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a RestEntries) Less(i, j int) bool {

	firstDynamic := isDynamic(a[i].Path)
	secondDynamic := isDynamic(a[j].Path)

	if firstDynamic && !secondDynamic {
		return false
	}

	if !firstDynamic && secondDynamic {
		return true
	}

	if len(a[i].Path) != len(a[j].Path) {
		return len(a[i].Path) > len(a[j].Path)
	}

	if a[i].Method != a[j].Method {
		return a[i].Method != "GET"
	}

	if a[i].Path == a[j].Path {
		panic(any("Two endpoints can't be the same"))
	}
	return true
}

func isDynamic(url string) bool {
	return strings.Contains(url, "{") && strings.Contains(url, "}")
}

// endregion
