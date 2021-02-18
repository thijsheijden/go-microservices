package main

import (
	"context"
	"frontend-service/internal/greeter"
	"html/template"
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
)

type kubernetesTestPage struct {
	PodName         string
	NodeName        string
	GreetingPodName string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Make a quick grpc request
	conn, err := grpc.Dial("greeter-service:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := greeter.NewGreetingServiceClient(conn)

	response, err := client.Greeter(context.Background(), &greeter.Greeting{Message: "WASSAAAA"})
	if err != nil {
		log.Fatalf("Something went wrong sending an RPC: %v", err)
	}

	data := kubernetesTestPage{PodName: os.Getenv("POD_NAME"), NodeName: os.Getenv("NODE_NAME"), GreetingPodName: response.GetResponse()}

	template, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatalf("Failed to parse index.html template: %v", err)
	}

	template.Execute(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":3000", nil)
}
