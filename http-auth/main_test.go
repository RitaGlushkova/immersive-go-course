package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEndpoints(t *testing.T) {
	t.Run("GET /", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		handlerIndex(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "<!DOCTYPE html><html><em>Hello, world</em>\n")
	})
	t.Run("GET / (WITH QUERY)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		q := request.URL.Query()
		q.Add("foo", "<strong>bar</strong>")
		request.URL.RawQuery = q.Encode()
		response := httptest.NewRecorder()
		handlerIndex(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "<!DOCTYPE html><html><em>Hello, world</em><p>Query parameters:<ul><li>foo: &lt;strong&gt;bar&lt;/strong&gt;</li></ul>\n")
	})
	t.Run("POST /", func(t *testing.T) {
		b := strings.NewReader("<em>Hi</em>")
		request, _ := http.NewRequest(http.MethodPost, "/", b)
		response := httptest.NewRecorder()
		handlerIndex(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)
		assertResponseBody(t, response.Body.String(), "<!DOCTYPE html><html>&lt;em&gt;Hi&lt;/em&gt;\n")
	})
	t.Run("GET /200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/200", nil)
		response := httptest.NewRecorder()
		handler200(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "200\n")
	})
	t.Run("GET /500", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/500", nil)
		response := httptest.NewRecorder()
		handler500(response, request)
		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertResponseBody(t, response.Body.String(), "Internal server error\n")
	})
	t.Run("GET /404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/404", nil)
		response := httptest.NewRecorder()
		server := http.NotFoundHandler()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
	t.Run("GET /authenticated (OK)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/authenticated", nil)
		request.Header.Add("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
		response := httptest.NewRecorder()
		handlerAuth(response, request)
		assertStatus(t, response.Code, http.StatusOK)
	})
	t.Run("GET /authenticated (FAIL)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/authenticated", nil)
		request.Header.Add("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ==")
		response := httptest.NewRecorder()
		handlerAuth(response, request)
		assertStatus(t, response.Code, http.StatusUnauthorized)
	})
	t.Run("GET /limited (OK)", func(t *testing.T) {
		limiter := rate.NewLimiter(100, 30)
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("[]"))
		})
		request, _ := http.NewRequest(http.MethodGet, "/limited", nil)
		var recorder []int
		handler := handlerLimit(limiter, testHandler)
		for i := 0; i < 32; i++ {
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, request)
			recorder = append(recorder, rr.Code)
		}
		fmt.Printf("%v", recorder)
		countOK := 0
		count429 := 0
		for _, v := range recorder {
			if v == 200 {
				countOK++
			} else {
				count429++
			}
		}
		got1, got2 := countOK, count429
		want1, want2 := 30, 2
		fmt.Println(got1, got2, want1, want2)
		if (got1 != want1) || (got2 != want2) {
			t.Errorf("response body is wrong, got number of successful requests %v and rejected %v want %v and %v", got1, got2, want1, want2)
		}
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
