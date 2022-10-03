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
	s := Server{
		username: "rita",
		password: "123",
		limiter:  rate.NewLimiter(100, 30),
	}

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
		assertResponseBody(t, response.Body.String(), "<!DOCTYPE html><html><em>Hello, world</em>\n<p>Query parameters:<ul><li>foo: [&lt;strong&gt;bar&lt;/strong&gt;]</li></ul>")
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
		// Computed by running echo -n Aladdin:open sesame | base64
		//(base64.StdEncoding.EncodeToString([]byte("Aladdin:open sesame"))).
		request.Header.Add("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
		response := httptest.NewRecorder()
		s.handlerAuth(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "<!DOCTYPE html><html><em>Hello, Aladdin!</em>\n")
	})
	t.Run("GET /authenticated (FAIL)", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/authenticated", nil)
		request.Header.Add("Authorization", "Basic dXNlcm5hbWU6cGFzc3dvcmQ==")
		response := httptest.NewRecorder()
		s.handlerAuth(response, request)
		assertStatus(t, response.Code, http.StatusUnauthorized)
	})
	t.Run("GET /limited (OK)", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/limited", nil)
		var responsesCodes = make(map[int]int)
		for i := 0; i < 32; i++ {
			rr := httptest.NewRecorder()
			s.handlerLimit(rr, request)
			_, ok := responsesCodes[rr.Code]
			if ok {
				responsesCodes[rr.Code]++
			} else {
				responsesCodes[rr.Code] = 1
			}
		}
		fmt.Println(responsesCodes)

		got429, ok429 := responsesCodes[429]

		if !ok429 {
			t.Errorf("Response Code 429 not found")
		}

		got200, ok200 := responsesCodes[200]

		if !ok200 {
			t.Errorf("Response Code 200 not found")
		}
		want200, want429 := 30, 2
		if got200 != want200 {
			t.Errorf("did not get correct of responses codes. Got number of successful requests %v, want %v. Got number of rejected requests %v, want %v", got200, want200, got429, want429)
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
