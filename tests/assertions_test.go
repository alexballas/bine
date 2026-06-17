package tests

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// requireAssertions is a small, fail-fast assertion helper built on top of the
// standard library testing package. It replaces the previous dependency on
// github.com/stretchr/testify, providing only the subset of assertions used by
// these integration tests. Every failing assertion aborts the current test via
// t.Fatalf.
type requireAssertions struct {
	t *testing.T
}

func newRequire(t *testing.T) *requireAssertions { return &requireAssertions{t: t} }

// NoError fails the test if err is non-nil.
func (r *requireAssertions) NoError(err error) {
	r.t.Helper()
	if err != nil {
		r.t.Fatalf("expected no error, got: %v", err)
	}
}

// True fails the test if value is false.
func (r *requireAssertions) True(value bool) {
	r.t.Helper()
	if !value {
		r.t.Fatal("expected condition to be true, got false")
	}
}

// Equal fails the test if expected and actual are not deeply equal.
func (r *requireAssertions) Equal(expected, actual any) {
	r.t.Helper()
	if !objectsAreEqual(expected, actual) {
		r.t.Fatalf("not equal:\n\texpected: %#v\n\tactual:   %#v", expected, actual)
	}
}

// Len fails the test if object does not have the given length.
func (r *requireAssertions) Len(object any, length int) {
	r.t.Helper()
	v := reflect.ValueOf(object)
	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Chan, reflect.String:
		if v.Len() != length {
			r.t.Fatalf("expected length %d, got %d for %#v", length, v.Len(), object)
		}
	default:
		r.t.Fatalf("cannot get length of %T", object)
	}
}

// Contains fails the test if container does not contain item. It supports
// substrings of strings, elements of slices/arrays and keys of maps.
func (r *requireAssertions) Contains(container, item any) {
	r.t.Helper()
	if !containsElement(container, item) {
		r.t.Fatalf("%#v does not contain %#v", container, item)
	}
}

// NotContains is the inverse of Contains.
func (r *requireAssertions) NotContains(container, item any) {
	r.t.Helper()
	if containsElement(container, item) {
		r.t.Fatalf("%#v should not contain %#v", container, item)
	}
}

// ElementsMatch fails the test unless listA and listB contain the same
// elements, ignoring order.
func (r *requireAssertions) ElementsMatch(listA, listB any) {
	r.t.Helper()
	if !elementsMatch(listA, listB) {
		r.t.Fatalf("elements differ:\n\tlistA: %#v\n\tlistB: %#v", listA, listB)
	}
}

func objectsAreEqual(expected, actual any) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}
	if exp, ok := expected.([]byte); ok {
		act, ok := actual.([]byte)
		if !ok {
			return false
		}
		return bytes.Equal(exp, act)
	}
	return reflect.DeepEqual(expected, actual)
}

// containsElement reports whether container contains item, mirroring the
// behaviour of testify's Contains: substrings for strings, elements for
// slices/arrays and keys for maps.
func containsElement(container, item any) bool {
	v := reflect.ValueOf(container)
	switch v.Kind() {
	case reflect.String:
		return strings.Contains(v.String(), fmt.Sprintf("%v", item))
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if objectsAreEqual(v.Index(i).Interface(), item) {
				return true
			}
		}
		return false
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if objectsAreEqual(key.Interface(), item) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// elementsMatch reports whether listA and listB hold the same elements
// regardless of order.
func elementsMatch(listA, listB any) bool {
	av := reflect.ValueOf(listA)
	bv := reflect.ValueOf(listB)
	if (av.Kind() != reflect.Slice && av.Kind() != reflect.Array) ||
		(bv.Kind() != reflect.Slice && bv.Kind() != reflect.Array) {
		return false
	}
	if av.Len() != bv.Len() {
		return false
	}
	matched := make([]bool, bv.Len())
	for i := 0; i < av.Len(); i++ {
		found := false
		for j := 0; j < bv.Len(); j++ {
			if matched[j] {
				continue
			}
			if objectsAreEqual(av.Index(i).Interface(), bv.Index(j).Interface()) {
				matched[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
