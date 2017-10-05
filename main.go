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

	http.Handle("/", NoCache(http.FileServer(http.Dir(*dir))))

	log.Printf("Serving %s on %s:%s\n", *dir, *host, *port)
	log.Fatal(http.ListenAndServe(*host+":"+*port, nil))
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// delete etags
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// set no cache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	})
}
