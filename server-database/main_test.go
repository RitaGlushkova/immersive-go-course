package main

import (
	"bytes"
	"context"
	"encoding/json"
	
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
	_, err = conn.Exec(context.Background(), `INSERT INTO public.images (title, url, alt_text) VALUES ($1, $2, $3), ($4, $5, $6)`, DataForTests[0].Title, DataForTests[0].URL, DataForTests[0].AltText, DataForTests[1].Title, DataForTests[1].URL, DataForTests[1].AltText)
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
		assertStatus(t, response.Code, http.StatusOK)
		require.ElementsMatch(t, got, DataForTests)
	})
	t.Run("GET /images.json (WITH QUERY)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/images.json", nil)
		q := request.URL.Query()
		q.Add("indent", "2")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
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

	t.Run("POST /images.json (WITHOUT QUERY)", func(t *testing.T) {
		imageToSave := []byte(`{"title":"Big Ben","url":"https://images.freeimages.com/images/large-previews/3d0/london-1452422.jpg","alt_text":"Big Ben and Tube sign"}`)
		request, err := http.NewRequest(http.MethodPost, "/images.json", bytes.NewBuffer(imageToSave))
		if err != nil {
			t.Fatalf("Unable to save data: %s", err.Error())
		}
		response := httptest.NewRecorder()
		s.handlerImages(response, request)
		require.Equal(t, response.Body.String(), string(imageToSave))
	})
	t.Run("POST /images.json (WITH QUERY)", func(t *testing.T) {
		imageToSave := []byte(`{"title":"Big Ben","url":"https://images.freeimages.com/images/large-previews/3d0/london-1452422.jpg","alt_text":"Big Ben and Tube sign"}`)
		request, err := http.NewRequest(http.MethodPost, "/images.json", bytes.NewBuffer(imageToSave))
		if err != nil {
			t.Fatalf("Unable to save data: %s", err.Error())
		}
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
		require.Equal(t, response.Body.String(), want)
	})
}

func TestEncode(t *testing.T) {
	errBuf := bytes.NewBuffer(nil)
	_, err := EncodedMarshalJSON(make(chan int), "+", errBuf)
	expected := "Couldn't proceed with Marshal:"
	got := errBuf.String()
	//fmt.Println(got)
	require.Contains(t, got, expected)
	require.Error(t, err)
}

func TestFetchFunc(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	conn, err := pgx.Connect(context.Background(), TEST_DB_URL)
	require.NoError(t, err)
	s := &Server{conn: conn}
	_, err = FetchImages(s.conn, `SELECT title, url, alt_text FROM public.imag`)
	expected := `ERROR: relation "public.imag" does not exist (SQLSTATE 42P01)`
	require.EqualError(t, err, expected)
}

func TestSaveImageFailsOnBadJSON( t *testing.T) {
	//clean DB prepare for use
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	//connect to database
	conn, err := pgx.Connect(context.Background(), TEST_DB_URL)
	require.NoError(t, err)
	defer conn.Close(context.Background())
	//call saseImage function with bad JSON
	body := bytes.NewBufferString(`hello`)
	img, err := saveImage(conn, body )
	//Assert for an error
	require.Error(t, err)
	// Assert the error message are clear

	//if should return nil instead of img Json
	require.Nil(t, img)
	// assert that nothing is added to the DB
	images, errFetch := FetchImages(conn, `SELECT * FROM public.images`)
	require.NoError(t, errFetch)
	require.Empty(t, images)
}

// func requireEmtryDB () {
// 	return 
// }



// func assertError (t testing.TB, err error, expected string) {
// 	t.Helper()
// 	if !strings.Contains(err.Error(), expected) {
// 		t.Fatalf("expected \"%s\" got \"%s\"", expected, err.Error())
// 	}
// }

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
