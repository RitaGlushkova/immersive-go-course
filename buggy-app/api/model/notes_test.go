package model

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	text := "This is an example #tag1 #tag2"
	expected := []string{"tag1", "tag2"}

	tags := extractTags(text)

	if !reflect.DeepEqual(expected, tags) {
		t.Fatalf("expected %v, got %v", expected, tags)
	}
}

func TestTagsTrim(t *testing.T) {
	text := "This is an example #tag1    #tag2    "
	expected := []string{"tag1", "tag2"}

	tags := extractTags(text)

	if !reflect.DeepEqual(expected, tags) {
		t.Fatalf("expected %v, got %v", expected, tags)
	}
}

func TestGetNotesForOwner(t *testing.T) {
	currTime := time.Now()
	ctx := context.Background()
	mock, err := pgxmock.NewPool(pgxmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectQuery("^SELECT (.+) FROM public.note WHERE owner = (.+)$").
		WithArgs("Qz6H-svr").
		WillReturnRows(pgxmock.NewRows([]string{"id", "owner", "content", "created", "modified"}).AddRow("id1", "test", "content1 #tag1", currTime, currTime))
	notes, err := GetNotesForOwner(ctx, mock, "457")
	require.Nil(t, notes)
	require.Error(t, err)
}
