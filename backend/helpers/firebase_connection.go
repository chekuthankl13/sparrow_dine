package helpers

import (
	"context"
	"log"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	var err error
	opt := option.WithCredentialsFile("sparrowdine-7dca2-firebase-adminsdk-qu31m-c913048e8c.json")
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err)
	}

}
