package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"encoding/binary"

	"github.com/saljam/mjpeg"
)

type jpeg []byte

func redirectJPEGs(ctx context.Context, mjpegPort, clientPort uint16) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream := mjpeg.NewStream()
	stream.FrameInterval = 30 * time.Millisecond

	log("stream created")
	log(stream)

	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", mjpegPort), stream)
	log(fmt.Sprintf("started listening to: 0.0.0.0:%d", mjpegPort))

	for ctx.Err() == nil {
		jpegChan, err := listenToJPEG(ctx, clientPort)
		if err != nil {
			return err
		}
		log("jpeg channel opened")

		for jpg := range jpegChan {
			log("jpeg received")
			stream.UpdateJPEG(jpg)
			log("stream updated")
		}

	}
	return nil
}

func listenToJPEG(ctx context.Context, port uint16) (<-chan jpeg, error) {
	ctx, cancel := context.WithCancel(ctx)

	lb, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		cancel()
		return nil, err
	}
	log("started listening to port: ", port)

	con, err := lb.Accept()
	jpegChan := make(chan jpeg)

	go func() {
		<-ctx.Done()
		logErr(lb.Close())
		close(jpegChan)
	}()

	go func() {
		defer cancel()

		if logErr(err) != nil {
			return
		}

		chunk := make([]byte, 16*1024)
		buffer := &bytes.Buffer{}
		img := &bytes.Buffer{}
		nextImageSize := uint32(0)

		for readBytes, err := con.Read(chunk); logErr(err) != nil; readBytes, err = con.Read(chunk) {
			log("read chunk: ", chunk[:readBytes])

			buffer.Write(chunk[:readBytes])

			if nextImageSize == 0 && buffer.Len() >= 4 {
				sizeBytes := []byte{0, 0, 0, 0}
				buffer.Read(sizeBytes)
				nextImageSize = binary.BigEndian.Uint32(sizeBytes)
			}

			for uint32(img.Len()) < nextImageSize && buffer.Len() > 0 {
				tmp := make([]byte, min(int(nextImageSize)-img.Len(), buffer.Len()))
				read, _ := buffer.Read(tmp)
				tmp = tmp[:read]
				_, _ = img.Write(tmp)
			}

			if uint32(img.Len()) == nextImageSize {
				log("image complete ", img.Len())
				jpegChan <- img.Bytes()
				nextImageSize = 0
				img.Reset()
			}
		}
	}()

	return jpegChan, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func logErr(err error) error {
	if err != nil && err != io.EOF && err != context.Canceled {
		fmt.Printf("Error: %+v", err)
	}
	return err
}

func log(s ...interface{}) {
	fmt.Println(s...)
}
