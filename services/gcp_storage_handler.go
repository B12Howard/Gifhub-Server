package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

// FileUpload takes in a file and saves it to a specified GCP Cloud Storage bucket
func FileUpload(bucket string, object *os.File, fileName string) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	viper.SetConfigName("config")
	viper.AddConfigPath(exPath + "/config")
	viperErr := viper.ReadInConfig()

	if viperErr != nil {
		fmt.Printf("Error reading config file, %s", viperErr)
	}

	var gcpConfig GCPCloudStorageConfig
	viperErr = viper.Unmarshal(&gcpConfig)
	b, _ := json.Marshal(gcpConfig.GCPCLOUDSTORAGE)
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(b))

	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(fileName).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

	if _, err = io.Copy(wc, object); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(wc, "%v uploaded to %v.\n", fileName, bucket)

	return nil
}

// generateV4GetObjectSignedURL generates object signed URL with GET method.
func GenerateV4GetObjectSignedURL(bucketName, object string) (string, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println("generate", exPath)

	viper.SetConfigName("config")
	viper.AddConfigPath(exPath + "/config")
	viperErr := viper.ReadInConfig()

	if viperErr != nil {
		fmt.Printf("Error reading config file, %s", viperErr)
	}

	var gcpConfig GCPCloudStorageConfig
	viperErr = viper.Unmarshal(&gcpConfig)
	b, _ := json.Marshal(gcpConfig.GCPCLOUDSTORAGE)
	fmt.Println("v", b)

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(b))

	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Signing a URL requires credentials authorized to sign a URL. You can pass
	// these in through SignedURLOptions with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing.
	// In this example, none of these options are used, which means the SignedURL
	// function attempts to use the same authentication that was used to instantiate
	// the Storage client. This authentication must include a private key or have
	// iam.serviceAccounts.signBlob permissions.
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	u, err := client.Bucket(bucketName).SignedURL(object, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %v", bucketName, err)
	}

	return u, nil
}
