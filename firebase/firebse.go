package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	FirebaseApp  *firebase.App
	FirebaseAuth *auth.Client
}

func InitFirebase(credentials string) (*FirebaseClient, error) {
	opt := option.WithCredentialsFile(credentials)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	// Initialize Firebase Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Firebase Auth client: %v", err)
	}
	log.Println("Firebase initialized successfully")

	return &FirebaseClient{
		FirebaseApp:  app,
		FirebaseAuth: authClient,
	}, nil

}
