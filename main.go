package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

func fetchHeader(req *http.Request, key string) (string, bool) {
	if _, ok := req.Header[key]; ok {
		return req.Header.Get(key), true
	}
	return "", false
}

type BucketProxy struct {
	bucket *storage.BucketHandle
}

func (s BucketProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	objectPath := req.URL.Path[1:]
	object := s.bucket.Object(objectPath)

	ctx := req.Context()
	switch req.Method {
	case http.MethodHead:
		_, err := object.Attrs(ctx)
		if err != nil {
			if err == storage.ErrObjectNotExist {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "File not found")
			} else {
				http.Error(w, err.Error(), http.StatusBadGateway)
			}
			return
		}
	case http.MethodGet:
		rc, err := object.NewReader(ctx)
		if err != nil {
			if err == storage.ErrObjectNotExist {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "File not found")
			} else {
				http.Error(w, err.Error(), http.StatusBadGateway)
			}
			return
		}
		defer rc.Close()

		io.Copy(w, rc)
	case http.MethodPut:
		// Write the object to GCS
		wc := object.NewWriter(ctx)

		// Copy the supported headers over from the original request
		if val, ok := fetchHeader(req, "Content-Type"); ok {
			wc.ContentType = val
		}
		if val, ok := fetchHeader(req, "Content-Language"); ok {
			wc.ContentLanguage = val
		}
		if val, ok := fetchHeader(req, "Content-Encoding"); ok {
			wc.ContentEncoding = val
		}
		if val, ok := fetchHeader(req, "Content-Disposition"); ok {
			wc.ContentDisposition = val
		}
		if val, ok := fetchHeader(req, "Cache-Control"); ok {
			wc.CacheControl = val
		}

		if _, err := io.Copy(wc, req.Body); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if err := wc.Close(); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		fmt.Fprintf(w, "OK")
	default:
		msg := fmt.Sprintf("Method '%s' is not supported", req.Method)
		http.Error(w, msg, http.StatusMethodNotAllowed)
	}
}

// Start the HTTP server
func run(addr, bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)

	handler := BucketProxy{bucket}

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(handler)

	fmt.Println("Starting proxy server on address", addr, "for bucket", bucketName)
	return http.ListenAndServe(addr, n)
}

// Urfave cli action
func action(c *cli.Context) error {
	addr := c.String("addr")
	bucketName := c.String("bucket-name")
	if bucketName == "" {
		return fmt.Errorf("please specify a bucket name")
	}
	return run(addr, bucketName)
}

func main() {
	app := cli.NewApp()
	app.Name = "nix-store-gcs-proxy"
	app.Usage = "A HTTP nix store that proxies requests to Google Storage"
	app.Version = "0.0.1"
	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bucket-name",
			Usage: "name of the bucket to proxy the data to",
		},
		cli.StringFlag{
			Name:  "addr",
			Value: "localhost:3000",
			Usage: "listening address of the HTTP server",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
