package shortener

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
)

// getProject gets the project ID.
func getProject() (string, error) {
	// Test if we're running on GCE.
	if metadata.OnGCE() {
		// Use the GCE Metadata service.
		projectID, err := metadata.ProjectID()
		if err != nil {
			log.Printf("Failed to get project ID from instance metadata")
			return "", err
		}
		return projectID, nil
	}
	// Shell out to gcloud.
	cmd := exec.Command("gcloud", "config", "get-value", "project")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to shell out to gcloud: %+v", err)
		return "", err
	}
	projectID := strings.TrimSuffix(out.String(), "\n")
	return projectID, nil
}

func getFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	if fsClient != nil {
		return fsClient, nil
	}

	project, err := getProject()
	if err != nil {
		return nil, err
	}
	fsClient, err = firestore.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	return fsClient, nil
}

func getBigqueryClient(ctx context.Context) (*bigquery.Client, error) {
	if bqClient != nil {
		return bqClient, nil
	}

	project, err := getProject()
	if err != nil {
		return nil, err
	}
	bqClient, err = bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, err
	}
	return bqClient, nil
}
