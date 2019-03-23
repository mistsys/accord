package status

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/mistsys/mist_go_utils/cloud"
	"github.com/mistsys/mist_go_utils/pprofutil"
)

var defaultHttpPort = 9110

// startup time of this process
var STARTTIME = time.Now().UTC()

// GITCOMMIT is the git commit id (the hex of the sha1 hash) of the source code used when building this code
// It should be set to the head commit of the master repo, not mist_go_utils. In the release code's repo
// the mist_go_utils should be vendored, and from that the mist_go_utils version can be derived from the main
// repo's commit.
var GITCOMMIT string

// VERSION is a version string for humans. We can try and set it to the circle-ci supplied version tag
var VERSION string

// BUILDTIME is the time at which the code was compiled
var BUILDTIME string

// a useful handler for the "/about" URL. You can register this with
//   http.HandleFunc("/about", cloud.HandleAbout)
func HandleAbout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
	"version": %q,
	"git-commit": %q,
	"build-time": %q,
	"go-runtime": %q,
	"start-time": %q,
	"uptime": %f
}`, VERSION, GITCOMMIT, BUILDTIME, runtime.Version(), STARTTIME.Format(time.RFC3339), time.Since(STARTTIME).Seconds())
}

func Serve() Mux {
	return ServePort(defaultHttpPort)
}

func ServePort(port int) Mux {
	addr := ":" + strconv.Itoa(port)
	mux := http.NewServeMux()
	mux.HandleFunc("/about", HandleAbout)

	pprofutil.InitMux(mux)

	go func(addr string, mux http.Handler) {
		err := http.ListenAndServe(addr, mux)
		msg := fmt.Sprintf("Error: can't serve metrics at %q: %s", addr, err)
		fmt.Println(msg)
		if !cloud.IsPrivateEnv() {
			panic(msg)
		}
	}(addr, mux)

	return mux
}

type Mux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	Handle(pattern string, handler http.Handler)
}
