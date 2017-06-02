package pusher

import (
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func lookForPush(t *testing.T, w http.ResponseWriter) {
	v := runtime.Version()
	if !PushCapable(w) {
		log.Println("Current go version", v, "doesn't implement HTTP/2 Server Push")
		if strings.HasPrefix(v, "go") {
			parts := strings.Split(v[2:], ".")
			ubermajor, err := strconv.Atoi(parts[0])
			if err != nil {
				t.Error(err)
			}
			if ubermajor > 1 {
				t.Error("Expected to find http.Pusher in " + v)
				return
			}
			major, err := strconv.Atoi(parts[1])
			if err != nil {
				t.Error(err)
			}
			if major >= 8 {
				t.Error("Expected to find http.Pusher in " + v)
			}
		}
		return
	}
	log.Println("Current go version", v, "implements HTTP/2 Server Push")
	if strings.HasPrefix(v, "go") {
		parts := strings.Split(v[2:], ".")
		ubermajor, err := strconv.Atoi(parts[0])
		if err != nil {
			t.Error(err)
		}
		if ubermajor > 1 {
			return
		}
		major, err := strconv.Atoi(parts[1])
		if err != nil {
			t.Error(err)
		}
		if major < 8 {
			t.Error("Didn't expect to find http.Pusher in " + v)
		}
	}
	pushed := Push(w, "/app.css")
	if !pushed {
		t.Error("Expected sucessful push")
	}
}

func TestPushCapable(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lookForPush(t, w)
	})
	go http.ListenAndServe(":8080", nil)
	// http.Get("http://localhost:8080/")

	// req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// _, err = http.DefaultClient.Do(req)
	// if err != nil {
	// 	t.Error(err)
	// }

	// TODO: make a local HTTP/2 server+client dialog??
}

func TestFakePusher(t *testing.T) {
	var p fakePusher
	lookForPush(t, p)
}

type fakePusher struct{}

func (fakePusher) Header() http.Header         { return nil }
func (fakePusher) Write(b []byte) (int, error) { return len(b), nil }
func (fakePusher) WriteHeader(int)             {}
func (fakePusher) Push(target string, opts *http.PushOptions) error {
	// TODO: a build tag to not depend on the actual existence of http.PushOptions
	log.Println("Pushing", target, ":)")
	return nil
}
