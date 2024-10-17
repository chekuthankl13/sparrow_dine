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
	opt := option.WithCredentialsFile("service.json")
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err)
	}

}
