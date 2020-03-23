package shortener

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
)

const (
	rCollection = "Redirects"
	datasetName = "Redirects"
	tableName   = "Usage"
)

var (
	fsClient *firestore.Client
	bqClient *bigquery.Client
)

// Redirect is a single redirect instruction.
type Redirect struct {
	Key         string
	URL         string
	RequestInfo *RequestInfo
}

// RequestInfo is what we can glean about a user's request.
type RequestInfo struct {
	RemoteAddr string
	URL        string
	Referer    string
	UserAgent  string
	Timestamp  time.Time
}

// NewRequestInfo returns a RequestInfo struct by examining a http.Request.
func NewRequestInfo(r *http.Request) *RequestInfo {
	return &RequestInfo{
		RemoteAddr: r.RemoteAddr,
		URL:        r.URL.String(),
		Referer:    r.Referer(),
		UserAgent:  r.UserAgent(),
		Timestamp:  time.Now(),
	}
}

// URLHandler retrieves URLs from the database and serves a redirect.
func URLHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	// Retrieve key from Datastore
	ctx := r.Context()
	client, err := getFirestoreClient(ctx)
	if err != nil {
		log.Printf("Could not create Firestore client: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	snap, err := client.Collection(rCollection).Doc(key).Get(ctx)
	if err != nil {
		log.Printf("Could not retrieve Firestore document: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var red Redirect
	snap.DataTo(&red)

	// Redirect.  URL must include scheme (http, https)
	defer storeUsage(ctx, &red, r)
	http.Redirect(w, r, red.URL, 303) // https://en.wikipedia.org/wiki/HTTP_303
}

func storeUsage(ctx context.Context, red *Redirect, r *http.Request) {
	// Ensure BigQuery dataset/tables
	client, err := getBigqueryClient(ctx)
	if err != nil {
		log.Printf("Could not get BigQuery client, failing write")
		return
	}
	schema, err := bigquery.InferSchema(Redirect{})
	dataset := client.Dataset(datasetName)
	_ = dataset.Create(ctx, &bigquery.DatasetMetadata{Location: "US"})
	table := dataset.Table(tableName)
	_ = table.Create(ctx, &bigquery.TableMetadata{Name: "Redirect Usage", Schema: schema})

	// Overwrite RequestInfo
	red.RequestInfo = NewRequestInfo(r)

	// Store redirect
	ins := table.Inserter()
	if err := ins.Put(ctx, red); err != nil {
		log.Printf("Stream to BigQuery failed: %+v", err)
		return
	}
}

// RootGetHandler handles a get to /.
func RootGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post {'Key': '...', 'URL': '...'} to this URL to store a redirect\n")
}

// RootPostHandler handles saving URLs to the database.
func RootPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Deserialize data
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not read from http.Request body: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer r.Body.Close()
	var red Redirect
	err = json.Unmarshal(b, &red)
	if err != nil {
		log.Printf("Could not unmarshal JSON: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write URL to database
	client, err := getFirestoreClient(ctx)
	if err != nil {
		log.Printf("Could not create Firestore client: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = client.Collection(rCollection).Doc(red.Key).Set(ctx, red)
	if err != nil {
		log.Printf("Could not store Firestore object: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, "OK\n")
}
