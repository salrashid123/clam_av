package gcfscan

import (
	"log"
	"fmt"
	"context"
	//"golang.org/x/net/context"
	"cloud.google.com/go/storage"
	"os"

	"github.com/dutchcoders/go-clamd"
	"cloud.google.com/go/functions/metadata"
)

var (
)

const (
	logName      = "go-functions"
	resourceType = "cloud_function"	
)

type Event struct {
	Bucket                  string `json:"bucket"`
	Name                    string `json:"name"`
	ContentType             string `json:"contentType"`
	Crc32c                  string `json:"crc32c"`
	Etag                    string `json:"etag"`
	Generation              string `json:"generation"`
	ID                      string `json:"id"`
	Kind                    string `json:"kind"`
	Md5Hash                 string `json:"md5Hash"`
	MediaLink               string `json:"mediaLink"`
	Metageneration          string `json:"metageneration"`
	SelfLink                string `json:"selfLink"`
	Size                    string `json:"size"`
	StorageClass            string `json:"storageClass"`
	TimeCreated             string `json:"timeCreated"`
	TimeStorageClassUpdated string `json:"timeStorageClassUpdated"`
	Updated                 string `json:"updated"`
}


func Scanner(ctx context.Context, event Event) error {

	meta, err := metadata.FromContext(ctx)
	if err != nil {
			return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Event ID: %v\n", meta.EventID)
	log.Printf("Event type: %v\n", meta.EventType)	

    sink := os.Getenv("BUCKET_DST")
    if len(sink) == 0 {
        log.Printf("Environment variable `BUCKET_DST` not set.")
    }

	ilbIP := os.Getenv("ILB_IP")
    if len(ilbIP) == 0 {
        log.Fatalf("Environment variable `ILB_IP` must be set and not empty.")
    }

	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("[Scanner] storage.NewCLient (%s) ", err)
	}
	defer gcsClient.Close()

	// read source file from GCS
	srcBucket := gcsClient.Bucket(event.Bucket)
	gcsSrcObject := srcBucket.Object(event.Name)
	gcsSrcReader, err := gcsSrcObject.NewReader(ctx)
	if err != nil {
		log.Fatalf("[Scanner] gcsReader (%s) ", err)
	}
	defer gcsSrcReader.Close()	

	log.Printf("[Scanner] Received: (%s) %s", event.Bucket, event.Name)

	c := clamd.NewClamd("tcp://" + ilbIP + ":3310")

	log.Printf("Ping: %v\n", c.Ping())

	stats, err := c.Stats()
	log.Printf("%v %v\n", stats, err)

	response, err := c.ScanStream(gcsSrcReader, make(chan bool))
	if err != nil {
		log.Fatalf("[Scanner] unable to scan file (%s) ", err)
	}	
	for s := range response {
		//log.Printf("%v\n", s.Status)

		switch s.Status {
		case "OK":
			log.Printf("Scan OK: (%s) %s", event.Bucket, event.Name)
		case "FOUND":
			log.Printf(">>>>>>> Virus Found in File file:(%s) %s", event.Bucket, event.Name)
			// do something here
		default:
			log.Printf("Unable to scan file:(%s) %s", event.Bucket, event.Name)
		}
	}

	return nil
}
