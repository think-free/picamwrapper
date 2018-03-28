package httpserver

import (
	"net/http"
	"strconv"

	"github.com/think-free/other/picamwrapper/config"
)

// HTTPPicam is the http server
type HTTPPicam struct {
	stat     bool
	conf     *config.Config
	internal *config.Internal
}

// New create a new HttpPicam object
func New(conf *config.Config, internal *config.Internal) *HTTPPicam {

	return &HTTPPicam{conf: conf, internal: internal}
}

// Run start the server
func (h *HTTPPicam) Run() {

	http.HandleFunc("/record/start", func(w http.ResponseWriter, r *http.Request) {

		h.internal.Chwritest <- true
		w.Write([]byte("Ok"))
	})
	http.HandleFunc("/record/stop", func(w http.ResponseWriter, r *http.Request) {

		h.internal.Chwritest <- false
		w.Write([]byte("Ok"))
	})

	http.HandleFunc("/record/status", func(w http.ResponseWriter, r *http.Request) {

		h.internal.Lock()
		content := strconv.FormatBool(h.internal.State)
		h.internal.Unlock()

		w.Write([]byte(content))
	})

	http.HandleFunc("/auto/start", func(w http.ResponseWriter, r *http.Request) {

		h.internal.AutoMode <- true
		w.Write([]byte("Ok"))
	})
	http.HandleFunc("/auto/stop", func(w http.ResponseWriter, r *http.Request) {

		h.internal.AutoMode <- false
		w.Write([]byte("Ok"))
	})

	http.HandleFunc("/auto/status", func(w http.ResponseWriter, r *http.Request) {

		h.internal.Lock()
		content := strconv.FormatBool(h.internal.Auto)
		h.internal.Unlock()

		w.Write([]byte(content))
	})

	http.ListenAndServe(h.conf.HTTPListen, nil)
}
