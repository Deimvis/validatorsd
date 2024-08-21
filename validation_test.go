package validatorsd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidateSelf(t *testing.T) {
	testCases := []validateTestCase{
		{
			&A{V: 42},
			nil,
		},
		{
			&A{V: 1},
			errors.New("wrong value"),
		},
		{
			&A{K: "key", V: 42},
			errors.New("non-empty key"),
		},
		{
			&A{K: "key", V: 1},
			errors.New("non-empty key"),
		},
	}
	runValidateSelfRecursively(t, testCases)
}

func Test_ValidateSelfFromNil(t *testing.T) {
	require.Panics(t, func() {
		ValidateSelfRecursively(nil)
	})
}

func Test_ValidateSelfWithRecursion(t *testing.T) {
	testCases := []validateTestCase{
		{
			&B{Valid: true, A: A{V: 42}},
			nil,
		},
		{
			&B{Valid: false, A: A{V: 42}},
			errors.New("not valid"),
		},
		{
			&B{Valid: true, A: A{V: 1}},
			errors.New("wrong value"),
		},
		{
			&С{},
			nil,
		},
		{
			&С{B: &B{Valid: true, A: A{V: 1}}},
			errors.New("wrong value"),
		},
	}
	runValidateSelfRecursively(t, testCases)
}

func Test_ValidateSelfWithEmbeds(t *testing.T) {
	testCases := []validateTestCase{
		{
			&D{A: A{V: 42}},
			nil,
		},
		{
			&D{A: A{V: 1}},
			errors.New("wrong value"),
		},
	}
	runValidateSelfRecursively(t, testCases)
}

func Test_ValidateSelfFromValue(t *testing.T) {
	testCases := []validateTestCase{
		{
			ValidatableByValue{Valid: true},
			nil,
		},
		{
			ValidatableByValue{Valid: false},
			errors.New("not valid"),
		},
		{
			ValidatableByPointer{Valid: true},
			nil,
		},
		{
			ValidatableByPointer{Valid: false},
			errors.New("not valid"),
		},
	}
	runValidateSelfRecursively(t, testCases)
}

func Test_ValidateSelfFromPointer(t *testing.T) {
	testCases := []validateTestCase{
		{
			&ValidatableByValue{Valid: true},
			nil,
		},
		{
			&ValidatableByValue{Valid: false},
			errors.New("not valid"),
		},
		{
			&ValidatableByPointer{Valid: true},
			nil,
		},
		{
			&ValidatableByPointer{Valid: false},
			errors.New("not valid"),
		},
	}
	runValidateSelfRecursively(t, testCases)
}

func runValidateSelfRecursively(t *testing.T, testCases []validateTestCase) {
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual := ValidateSelfRecursively(tc.obj)
			require.Equal(t, (tc.expected != nil), (actual != nil))
			if tc.expected != nil {
				require.Equal(t, tc.expected.Error(), actual.Error())
			}
		})
	}
}

type validateTestCase struct {
	obj      interface{}
	expected error
}

type A struct {
	K string
	V int
}

func (a *A) ValidateSelf() error {
	if len(a.K) > 0 {
		return errors.New("non-empty key")
	}
	if a.V != 42 {
		return errors.New("wrong value")
	}
	return nil
}

type B struct {
	Valid bool
	A     A
}

func (b *B) ValidateSelf() error {
	if !b.Valid {
		return errors.New("not valid")
	}
	return nil
}

type С struct {
	B *B
}

type D struct {
	A
	other string
}

type ValidatableByValue struct {
	Valid bool
}

func (v ValidatableByValue) ValidateSelf() error {
	if !v.Valid {
		return errors.New("not valid")
	}
	return nil
}

type ValidatableByPointer struct {
	Valid bool
}

func (v *ValidatableByPointer) ValidateSelf() error {
	if !v.Valid {
		return errors.New("not valid")
	}
	return nil
}
