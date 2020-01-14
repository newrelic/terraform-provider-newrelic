// +build integration

package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationListTags(t *testing.T) {
	t.Parallel()

	var (
		testGUID = "MjUwODI1OXxBUE18QVBQTElDQVRJT058MjA0MjYxMzY4"
	)

	client := newIntegrationTestClient(t)

	actual, err := client.ListTags(testGUID)

	require.NoError(t, err)
	require.Greater(t, len(actual), 0)
}

func TestIntegrationAddTags(t *testing.T) {
	t.Parallel()

	var (
		testGUID = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
	)

	client := newIntegrationTestClient(t)

	tags := []Tag{
		{
			Key:    "test",
			Values: []string{"value"},
		},
	}
	err := client.AddTags(testGUID, tags)

	require.NoError(t, err)
}

func TestIntegrationReplaceTags(t *testing.T) {
	t.Parallel()

	var (
		testGUID = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
	)

	client := newIntegrationTestClient(t)

	tags := []Tag{
		{
			Key:    "test",
			Values: []string{"value"},
		},
	}
	err := client.ReplaceTags(testGUID, tags)

	require.NoError(t, err)
}

func TestIntegrationDeleteTags(t *testing.T) {
	t.Parallel()

	var (
		testGUID = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
	)

	client := newIntegrationTestClient(t)

	tagKeys := []string{"test"}
	err := client.DeleteTags(testGUID, tagKeys)

	require.NoError(t, err)
}

func TestIntegrationDeleteTagValues(t *testing.T) {
	t.Parallel()

	var (
		testGUID = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
	)

	client := newIntegrationTestClient(t)

	tagValues := []TagValue{
		{
			Key:   "test",
			Value: "value",
		},
	}
	err := client.DeleteTagValues(testGUID, tagValues)

	require.NoError(t, err)
}
