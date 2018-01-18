package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	port := flag.String("port", "80", "Port")
	dir := flag.String("dir", ".", "Directory")
	host := flag.String("host", "home", "Hostname")
	flag.Parse()

	http.Handle("/", noCache(http.FileServer(http.Dir(*dir))))

	log.Printf("Serving %s on %s:%s\n", *dir, *host, *port)
	log.Fatal(http.ListenAndServe(*host+":"+*port, nil))
}

// wraps a handler to disable caching
func noCache(h http.Handler) http.Handler {
	var headersToSet = map[string]string{
		"Expires":         time.Unix(0, 0).Format(time.RFC1123),
		"Cache-Control":   "no-cache, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}

	var headersToDelete = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// delete
		for _, v := range headersToDelete {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// set
		for k, v := range headersToSet {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	})
}
