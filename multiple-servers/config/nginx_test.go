package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//if you send to static server localhost:8089/ - tests HTML/CSS/JS
//send request to api server - receive successful response
//send 10 requests to api server - each server is getting it at least once
func startServer (port int, path string, res string, numberOfReqPerPort map[int]int) {
	//start a server which han have isolated handlers
	serveMux := http.NewServeMux()
	serveMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		numberOfReqPerPort[port] += 1 
		fmt.Printf("server got request on port: %d\n", port)
		w.Write([]byte(res))
	})
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: serveMux,
	}
	go server.ListenAndServe()
}

func TestMain(t *testing.T) {
	numberOfReqPerPort := make(map[int]int)
	//decide what we send to nginx and what we expect back
	//staticRequest, _ := http.NewRequest(http.MethodGet, "/", nil)
	wantStaticResponseBody := "Hello static"
	startServer(8082, "/", wantStaticResponseBody, numberOfReqPerPort)
	wantApiResponseBody := "Hello api"
	startServer(8081, "/images.json", wantApiResponseBody, numberOfReqPerPort)
	startServer(8083, "/images.json", wantApiResponseBody, numberOfReqPerPort)
	startServer(8084, "/images.json", wantApiResponseBody, numberOfReqPerPort)
	// start mock servers
	abs, err := filepath.Abs("./nginx.conf")

	// Printing if there is no error
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("nginx", "-c", abs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	time.Sleep(1 * time.Second)
	//start nginx
	//nginx -c "`pwd`/config/nginx.conf"
	//send stuff
	resStatic := requestToServerAndResponseWithBody(t, "http://localhost:8089/")
	resApi := requestToServerAndResponseWithBody(t, "http://localhost:8089/api/images.json")

	require.Equal(t, wantStaticResponseBody, resStatic)
	require.Equal(t, wantApiResponseBody, resApi)

	for i := 1; i < 10; i++ {
 		requestToServerAndResponseWithBody(t, "http://localhost:8089/api/images.json")
	}
	if numberOfReqPerPort[8081] == 0 {
		t.Fatal("Expected to have at least one request")
	}
	if numberOfReqPerPort[8083] == 0 {
		t.Fatal("Expected to have at least one request")
	}
	if numberOfReqPerPort[8084] == 0 {
		t.Fatal("Expected to have at least one request")
	}
}

func requestToServerAndResponseWithBody(t *testing.T, url string) string{
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	//assert that we received what is expected
	bytes, err := io.ReadAll(res.Body)
    if err != nil {
        t.Fatal(err)
    }
	return string(bytes)
}

//Dockerise it
//start actual servers