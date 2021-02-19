package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"frontend-service/internal/image_service"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type kubernetesTestPage struct {
	PodName         string
	NodeName        string
	GreetingPodName string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Make a quick grpc request
	// conn, err := grpc.Dial("greeter-service:9000", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("Failed to connect: %v", err)
	// }
	// defer conn.Close()

	// client := greeter.NewGreetingServiceClient(conn)

	// response, err := client.Greeter(context.Background(), &greeter.Greeting{Message: "WASSAAAA"})
	// if err != nil {
	// 	log.Fatalf("Something went wrong sending an RPC: %v", err)
	// }

	data := kubernetesTestPage{PodName: os.Getenv("POD_NAME"), NodeName: os.Getenv("NODE_NAME"), GreetingPodName: "POOP"}

	template, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatalf("Failed to parse index.html template: %v", err)
	}

	template.Execute(w, data)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Unable to establish websocket connection: %v", err)
	}
	defer ws.Close()

	messageType, data, err := ws.ReadMessage()
	if messageType == websocket.BinaryMessage {
		fmt.Println("Received data!")
	} else {
		fmt.Println("Text data?")
	}

	// Send this data to the image service
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error dialing grpc image service: %v", err)
	}
	defer conn.Close()

	client := image_service.NewImageGrayscaleServiceClient(conn)

	stream, err := client.UploadImage(context.Background())
	if err != nil {
		log.Fatalf("Something went wrong sending an RPC: %v", err)
	}
	waitc := make(chan struct{})

	// Response image buffer
	buffer := new(bytes.Buffer)

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if buffer.Len() < 10552 {
				log.Printf("Buffer size: %v", buffer.Len())
			}

			// Add chunk to buffer
			buffer.Write(in.GetProcessedImage())
		}
	}()

	reader := bufio.NewReader(bytes.NewReader(data))
	part := make([]byte, 1024)
	var count int

	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		stream.Send(&image_service.ImageRequest{Image: part[:count]})
	}
	if err != io.EOF {
		log.Fatalf("Error sending grayscale image in chunks: %v", err)
	}
	stream.CloseSend()
	<-waitc

	ws.WriteMessage(websocket.BinaryMessage, buffer.Bytes())
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/socket", imageHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":3000", nil)
}
