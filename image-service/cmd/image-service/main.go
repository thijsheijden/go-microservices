package main

import (
	"bufio"
	"bytes"
	"image"
	"image-service/internal/image_service"
	"image/png"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	image_service.UnimplementedImageGrayscaleServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen to port 9000: %v", err)
	}

	server := server{}

	grpcServer := grpc.NewServer()
	image_service.RegisterImageGrayscaleServiceServer(grpcServer, &server)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC failed to serve on port 9000: %v", err)
	}
}

func (s *server) UploadImage(stream image_service.ImageGrayscaleService_UploadImageServer) error {
	// New byte buffer
	var buffer bytes.Buffer

	// Keep reading new data until we reach EOF
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// Create the grayscale image and return that image to the client
			img := convertImageToGrayScale(buffer.Bytes())
			reader := bufio.NewReader(bytes.NewReader(img))
			part := make([]byte, 1024)
			var count int

			for {
				if count, err = reader.Read(part); err != nil {
					return nil
				}
				stream.Send(&image_service.ProgressResponse{Data: &image_service.ProgressResponse_ProcessedImage{ProcessedImage: part[:count]}})
			}
		}

		if err != nil {
			return err
		}

		// Add new data to the buffer
		buffer.Write(in.Image)
	}
}

func convertImageToGrayScale(buffer []byte) []byte {
	img, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		log.Fatalf("Error while decoding the image: %v", err)
	}

	// Convert image to grayscale
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	var grayBuffer bytes.Buffer
	err = png.Encode(&grayBuffer, grayImg)
	if err != nil {
		log.Fatalf("Failed to encode grayscale image: %v", err)
	}

	return grayBuffer.Bytes()
}
