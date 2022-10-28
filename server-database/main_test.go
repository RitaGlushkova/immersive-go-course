package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var TEST_DB_URL = "postgres://localhost:5432/go-server-database-test"
var DataForTests = []Image{
	{
		Title:   "White Cat",
		URL:     "https://images.freeimages.com/images/previews/13e/my-cat-1363423.jpg",
		AltText: "White cat sitting and looking to the left",
		Pixels:  6779900,
	},
	{
		Title:   "Catch a Ball",
		URL:     "https://images.freeimages.com/images/large-previews/12a/dog-1361473.jpg",
		AltText: "A dog jumping up catching a red ball",
		Pixels:  3333000,
	},
}

func setupSuite(tb testing.TB, addValues bool) func(tb testing.TB) {
	conn, err := pgx.Connect(context.Background(), TEST_DB_URL)
	if err != nil {
		tb.Fatalf("Teardown Error: Unable to connect to DB: %s", err.Error())
		os.Exit(1)
	}

	_, err = conn.Exec(context.Background(), "DELETE from public.images")

	if err != nil {
		tb.Fatalf("Teardown Error: Unable to delete from images: %s", err.Error())
	}
	if addValues {
		_, err = conn.Exec(context.Background(), `INSERT INTO public.images (title, url, alt_text, pixels) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)`, DataForTests[0].Title, DataForTests[0].URL, DataForTests[0].AltText, DataForTests[0].Pixels, DataForTests[1].Title, DataForTests[1].URL, DataForTests[1].AltText, DataForTests[1].Pixels)
		if err != nil {
			tb.Fatalf("Teardown Error: Unable to insert data: %s", err.Error())
		}
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
	teardownSuite := setupSuite(t, true)
	defer teardownSuite(t)
	conn, err := pgx.Connect(context.Background(), TEST_DB_URL)
	require.NoError(t, err)
	s := &Server{conn: conn}
	t.Run("GET /", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		handlerIndex(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		require.Equal(t, response.Body.String(), "Hello World")
	})
	t.Run("GET /images.json (WITHOUT QUERY)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/images.json", nil)
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		var got []Image
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Error(err)
		}
		assertStatus(t, http.StatusOK, response.Code)
		require.ElementsMatch(t, DataForTests, got)
	})
	t.Run("GET /images.json (WITH QUERY)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/images.json", nil)
		q := request.URL.Query()
		q.Add("indent", "2")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		want := "[\n  {\n    \"title\": \"White Cat\",\n    \"url\": \"https://images.freeimages.com/images/previews/13e/my-cat-1363423.jpg\",\n    \"alt_text\": \"White cat sitting and looking to the left\",\n    \"pixels\": 6779900\n  },\n  {\n    \"title\": \"Catch a Ball\",\n    \"url\": \"https://images.freeimages.com/images/large-previews/12a/dog-1361473.jpg\",\n    \"alt_text\": \"A dog jumping up catching a red ball\",\n    \"pixels\": 3333000\n  }\n]"
		assertStatus(t, http.StatusOK, response.Code)
		require.Equal(t, want, response.Body.String())
	})

	t.Run("POST /images.json (WITHOUT QUERY)", func(t *testing.T) {
		imageToSave := []byte(`{"title":"Big Ben","url":"https://images.freeimages.com/images/large-previews/3d0/london-1452422.jpg","alt_text":"Big Ben and Tube sign","pixels":1080000}`)
		request, err := http.NewRequest(http.MethodPost, "/images.json", bytes.NewBuffer(imageToSave))
		if err != nil {
			t.Fatalf("Unable to save data: %s", err.Error())
		}
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		require.Equal(t, string(imageToSave), response.Body.String())
	})
	t.Run("POST /images.json image already exists", func(t *testing.T) {
		imageToSave := []byte(`{"title":"Big Ben","url":"https://images.freeimages.com/images/large-previews/3d0/london-1452422.jpg","alt_text":"Big Ben and Tube sign"}`)
		request, err := http.NewRequest(http.MethodPost, "/images.json", bytes.NewBuffer(imageToSave))
		if err != nil {
			t.Fatalf("Unable to save data: %s", err.Error())
		}
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		want := "image with URL: https://images.freeimages.com/images/large-previews/3d0/london-1452422.jpg already exists\n"
		require.Equal(t, want, response.Body.String())
	})
	t.Run("POST /images.json (WITH QUERY)", func(t *testing.T) {
		imageToSave := []byte(`{"title": "Cat","url": "https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80","alt_text": "A cool cat"}`)
		request, err := http.NewRequest(http.MethodPost, "/images.json", bytes.NewBuffer(imageToSave))
		if err != nil {
			t.Fatalf("Unable to save data: %s", err.Error())
		}
		q := request.URL.Query()
		q.Add("indent", "2")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		want := "{\n  \"title\": \"Cat\",\n  \"url\": \"https://images.unsplash.com/photo-1533738363-b7f9aef128ce?ixlib=rb-1.2.1\\u0026ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8\\u0026auto=format\\u0026fit=crop\\u0026w=1000\\u0026q=80\",\n  \"alt_text\": \"A cool cat\",\n  \"pixels\": 1333000\n}"

		require.Equal(t, want, response.Body.String())
	})
}

func TestEncode(t *testing.T) {
	errBuf := bytes.NewBuffer(nil)
	_, err := MarshalJSON(make(chan int), "+", errBuf)
	expected := "couldn't proceed with Marshal:"
	got := errBuf.String()
	require.Contains(t, got, expected)
	require.Error(t, err)
}

func TestSaveImageFailsOnBadJSON(t *testing.T) {
	//clean DB prepare for use
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	//connect to database
	conn, err := pgx.Connect(context.Background(), TEST_DB_URL)
	require.NoError(t, err)
	defer conn.Close(context.Background())
	//call saseImage function with bad JSON
	body := bytes.NewBufferString(`hello`)
	img, err := saveImage(conn, body)
	//Assert for an error
	require.Error(t, err)
	// Assert the error message are clear

	//if should return nil instead of img Json
	require.Nil(t, img)
	// assert that nothing is added to the DB
	images, errFetch := FetchImages(conn)
	require.NoError(t, errFetch)
	require.Empty(t, images)
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, want %d, want %d", want, got)
	}
}
