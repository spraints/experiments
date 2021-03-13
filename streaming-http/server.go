// Usage: go run server.go
//   curl http://127.0.0.1:3434 | less
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	step = 100000
)

func main() {
	addr := flag.String("addr", "127.0.0.1:3434", "address for server to listen on")
	file := flag.String("file", "/usr/share/dict/words", "file to serve")
	count := flag.Int("n", 10, "number of times to repeat the file")

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("Listening on %v", *addr)
	http.ListenAndServe(*addr, &s{file: *file, count: *count})
}

type s struct {
	file  string
	count int
}

func (s *s) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(s.file)
	if err != nil {
		log.Printf("[%s] %s: %v", r.RemoteAddr, s.file, err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer f.Close()

	flusher := w.(http.Flusher)

	var rt, wt timer

	var buf [64 * 1024]byte
	xfer := 0
	nextReportAt := step
	for i := 0; i < s.count; i++ {
		_, err := f.Seek(0, 0)
		if err != nil {
			log.Printf("[%s] error: %v", r.RemoteAddr, err)
			return
		}
		for {
			stop := rt.start()
			n, err := f.Read(buf[:])
			stop()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("[%s] read error: %v", r.RemoteAddr, err)
				return
			}
			stop = wt.start()
			m, err := w.Write(buf[:n])
			stop()
			if err != nil {
				log.Printf("[%s] write error: %v", r.RemoteAddr, err)
				return
			}
			if n != m {
				log.Printf("[%s] error: incomplete write %d/%d bytes after %d bytes", r.RemoteAddr, m, n, xfer)
				return
			}
			flusher.Flush()
			xfer += n
			if xfer > nextReportAt {
				log.Printf("[%s] transferred %d bytes, %s in read, %s in write", r.RemoteAddr, xfer, rt.String(), wt.String())
				nextReportAt += step * (i + 1)
			}
		}
	}
	defer log.Printf("[%s] done!", r.RemoteAddr)
}

type timer struct {
	elapsed time.Duration
}

func (t *timer) start() func() {
	start := time.Now()
	return func() {
		t.elapsed += time.Since(start)
	}
}

func (t *timer) String() string {
	return t.elapsed.String()
}
