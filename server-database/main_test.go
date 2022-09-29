package main

import (
	//"bytes"
	"bytes"
	"context"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

var TEST_DB_URL = "postgres://localhost:5432/go-server-database-test"
var DataForTests = []Image{
	{
		Title:   "White Cat",
		URL:     "https://images.freeimages.com/images/previews/13e/my-cat-1363423.jpg",
		AltText: "White cat sitting and looking to the left",
	},
	{
		Title:   "Catch a Ball",
		URL:     "https://images.freeimages.com/images/large-previews/12a/dog-1361473.jpg",
		AltText: "A dog jumping up catching a red ball",
	},
}

func setupSuite(tb testing.TB) func(tb testing.TB) {
	conn, err := pgx.Connect(context.Background(), TEST_DB_URL)
	if err != nil {
		tb.Fatalf("Teardown Error: Unable to connect to DB: %s", err.Error())
		os.Exit(1)
	}

	_, err = conn.Exec(context.Background(), `INSERT INTO public.images (title, url, alt_text) VALUES ('White Cat','https://images.freeimages.com/images/previews/13e/my-cat-1363423.jpg','White cat sitting and looking to the left'),('Catch a Ball','https://images.freeimages.com/images/large-previews/12a/dog-1361473.jpg','A dog jumping up catching a red ball')`)
	if err != nil {
		tb.Fatalf("Teardown Error: Unable to insert data: %s", err.Error())
	}

	return func(tb testing.TB) {
		// teardown the database after testing
		_, err := conn.Exec(context.Background(), "DELETE from public.images")

		if err != nil {
			tb.Fatalf("Teardown Error: Unable to delete from images: %s", err.Error())
		}
	}
}

func TestMain(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	conn, err := pgx.Connect(context.Background(), TEST_DB_URL)
	if err != nil {
		t.Fatalf("Unable to insert data: %s", err.Error())
	}
	//require.NoError(t, err)
	s := &Server{conn: conn}
	// t.Run("GET /", func(t *testing.T) {
	// 	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	// 	response := httptest.NewRecorder()
	// 	handlerIndex(response, request)
	// 	assertStatus(t, response.Code, http.StatusOK)
	// 	require.Equal(t, response.Body.String(), "Hello World")
	// })
	t.Run("GET /images.json (WITHOUT QUERY)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/images.json", nil)
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		want := `[{"title":"White Cat","url":"https://images.freeimages.com/images/previews/13e/my-cat-1363423.jpg","alt_text":"White cat sitting and looking to the left"},{"title":"Catch a Ball","url":"https://images.freeimages.com/images/large-previews/12a/dog-1361473.jpg","alt_text":"A dog jumping up catching a red ball"}]`
		assertStatus(t, response.Code, http.StatusOK)
		require.Equal(t, want, response.Body.String())
	})
	t.Run("GET /images.json (WITH QUERY)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/images.json", nil)
		q := request.URL.Query()
		q.Add("indent", "2")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		if err != nil {
			t.Error(err)
		}
		want := `[
  {
    "title": "White Cat",
    "url": "https://images.freeimages.com/images/previews/13e/my-cat-1363423.jpg",
    "alt_text": "White cat sitting and looking to the left"
  },
  {
    "title": "Catch a Ball",
    "url": "https://images.freeimages.com/images/large-previews/12a/dog-1361473.jpg",
    "alt_text": "A dog jumping up catching a red ball"
  }
]`
		assertStatus(t, response.Code, http.StatusOK)
		require.Equal(t, want, response.Body.String())
	})

	t.Run("POST /images.json (WITH QUERY)", func(t *testing.T) {
		input := []byte(`{"title": "Big Ben", "url": "https://images.freeimages.com/images/large-previews/3d0/london-1452422.jpg","alt_text": "Big Ben and Tube sign"}`)
		request, _ := http.NewRequest(http.MethodPost, "/images.json", bytes.NewBuffer(input))
		q := request.URL.Query()
		q.Add("indent", "2")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		want := `{
  "title": "Big Ben",
  "url": "https://images.freeimages.com/images/large-previews/3d0/london-1452422.jpg",
  "alt_text": "Big Ben and Tube sign"
}`
		require.Equal(t, want, response.Body.String())

	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
