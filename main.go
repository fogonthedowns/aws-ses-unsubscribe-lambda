package main

import (
	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/fogonthedowns/aws-ses-unsubscribe-lambda/lib"
)

var env *lib.Lambda = &lib.Lambda{}

func init() {
	lib.Init(env)
}

func main() {
	runtime.Start(env.HandleRequest)
}
