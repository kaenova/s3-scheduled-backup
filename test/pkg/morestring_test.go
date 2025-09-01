package pkg_test

import (
	"testing"

	"github.com/kaenova/s3-scheduled-backup/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGetLastKStringOne(t *testing.T) {
	data := []string{"a", "b", "c"}
	expect := []string{"c"}
	finalData := pkg.GetLastKString(1, data)
	assert.EqualValues(t, expect, finalData)
}

func TestGetLastKStringTwo(t *testing.T) {
	data := []string{"a", "b", "c"}
	expect := []string{"b", "c"}
	finalData := pkg.GetLastKString(2, data)
	assert.EqualValues(t, expect, finalData)
}

func TestGetLastKStringThree(t *testing.T) {
	data := []string{"a", "b", "c"}
	expect := []string{"a", "b", "c"}
	finalData := pkg.GetLastKString(3, data)
	assert.EqualValues(t, expect, finalData)
}

func TestGetLastKStringFour(t *testing.T) {
	data := []string{"a", "b", "c"}
	assert.Panics(t, func() { pkg.GetLastKString(4, data) })
}

func TestGetFirstKStringOne(t *testing.T) {
	data := []string{"a", "b", "c"}
	expect := []string{"a"}
	finalData := pkg.GetFirstKString(1, data)
	assert.EqualValues(t, expect, finalData)
}

func TestGetFirstKStringTwo(t *testing.T) {
	data := []string{"a", "b", "c"}
	expect := []string{"a", "b"}
	finalData := pkg.GetFirstKString(2, data)
	assert.EqualValues(t, expect, finalData)
}

func TestGetFirstKStringThree(t *testing.T) {
	data := []string{"a", "b", "c"}
	expect := []string{"a", "b", "c"}
	finalData := pkg.GetFirstKString(3, data)
	assert.EqualValues(t, expect, finalData)
}

func TestGetFirstKStringFour(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	expect := []string{"a", "b", "c"}
	finalData := pkg.GetFirstKString(3, data)
	assert.EqualValues(t, expect, finalData)
}

func TestGetFirstKStringPanic(t *testing.T) {
	data := []string{"a"}
	assert.Panics(t, func() { pkg.GetLastKString(2, data) })
}
