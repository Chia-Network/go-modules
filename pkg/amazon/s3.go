package amazon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// MultiPartUploadInput holds the inputs for a multipart upload
type MultiPartUploadInput struct {
	Svc               *s3.S3          // Required: An AWS S3 session service for the upload
	Ctx               context.Context // Required: The context for this request
	CtxTimeout        time.Duration   // Optional: The request will time out after this duration (defaults to 60 minutes)
	MaxConcurrent     int             // Optional: The number of concurrent part uploads (defaults to 10)
	PartSize          int64           // Optional: Number of bytes (defaults to 8MB)
	Filepath          string          // Required: A full path to a local file to PUT to S3
	DestinationBucket string          // Required: The destination S3 bucket's name
	DestinationKey    string          // Required: The destination path in the bucket to put the file
	Logger            *slog.Logger    // Optional: Handles logging if supplied
}

// MultiPartUploadResult holds the result for an individual part upload
type MultiPartUploadResult struct {
	Error error
	Part  *s3.CompletedPart
}

// MultiPartUpload uploads a local file in multiple parts to AWS S3
func MultiPartUpload(input MultiPartUploadInput) error {
	// Exit if no S3 service given
	if input.Svc == nil {
		return fmt.Errorf("s3 service nil -- is a required option")
	}
	// Set part size to default 8MB if no part size specified or less than 5MB
	if input.PartSize < 5242880 {
		input.PartSize = 8 * 1024 * 1024
	}
	// Make sure max concurrent is at least 1, default to 10 if unspecified or less than 1
	if input.MaxConcurrent < 1 {
		input.MaxConcurrent = 10
	}
	// Set timeout to 60 minutes if not specified or zero value
	if input.CtxTimeout == 0 {
		input.CtxTimeout = 60 * time.Minute
	}

	// Set up context with timeout
	ctx, cancelFn := context.WithTimeout(input.Ctx, input.CtxTimeout)
	defer cancelFn()

	// Open local file
	file, err := os.Open(input.Filepath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			if input.Logger != nil {
				input.Logger.Error("encountered error closing file", "path", input.Filepath)
			}
		}
	}()

	// Get file and total file size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}
	fileSize := fileInfo.Size()

	// Initialize a multipart upload and get an upload ID back
	multipartUpload, err := input.Svc.CreateMultipartUploadWithContext(ctx, &s3.CreateMultipartUploadInput{
		Bucket: &input.DestinationBucket,
		Key:    &input.DestinationKey,
	})
	if err != nil {
		return fmt.Errorf("error creating multipart upload: %w", err)
	}

	// Record the upload ID from the multipart upload
	var uploadID string
	if multipartUpload != nil {
		if multipartUpload.UploadId != nil {
			if *multipartUpload.UploadId == "" {
				return errors.New("no upload ID returned in start upload request -- something wrong with the client or credentials?")
			}

			uploadID = *multipartUpload.UploadId
		}
	}

	// Get the total number of parts we will upload
	numParts := getTotalNumberParts(fileSize, input.PartSize)
	if input.Logger != nil {
		input.Logger.Debug("will upload file in parts", "file", input.Filepath, "parts", numParts)
	}

	var (
		wg  sync.WaitGroup
		ch  = make(chan error, numParts)
		sem = make(chan struct{}, input.MaxConcurrent)
	)

	// Start the individual part uploads
	orderedParts := make([]*s3.CompletedPart, numParts)
	for i := int64(0); i < numParts; i++ {
		partNumber := i + 1
		offset := i * input.PartSize
		bytesToRead := min(input.PartSize, fileSize-offset)

		partBuffer := make([]byte, bytesToRead)
		_, err := file.ReadAt(partBuffer, offset)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(partNumber int64, partBuffer []byte) {
			sem <- struct{}{}
			defer func() {
				<-sem
			}()
			defer wg.Done()

			if input.Logger != nil {
				input.Logger.Debug("uploading file part", "file", input.Filepath, "part", partNumber, "size", len(partBuffer))
			}

			resp, err := input.Svc.UploadPart(&s3.UploadPartInput{
				Bucket:     aws.String(input.DestinationBucket),
				Key:        aws.String(input.DestinationKey),
				UploadId:   &uploadID,
				PartNumber: aws.Int64(partNumber),
				Body:       bytes.NewReader(partBuffer),
			})
			if err != nil {
				ch <- fmt.Errorf("error uploading part %d : %w", partNumber, err)
				return
			}

			// Store the completed part in the uploadParts slice
			orderedParts[partNumber-1] = &s3.CompletedPart{
				ETag:       resp.ETag,
				PartNumber: aws.Int64(partNumber),
			}

			if input.Logger != nil {
				input.Logger.Debug("finished uploading file part", "file", input.Filepath, "part", partNumber, "size", len(partBuffer))
			}
		}(partNumber, partBuffer)
	}

	wg.Wait()

	// Check for errors from goroutines
	select {
	case err := <-ch:
		return err
	default:
		// No errors
	}
	close(ch)

	// Make a final call to AWS to say the file upload is complete
	// The file won't show up in S3 unless this is called
	_, err = input.Svc.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   &input.DestinationBucket,
		Key:      &input.DestinationKey,
		UploadId: &uploadID,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: orderedParts,
		},
	})
	if err != nil {
		return fmt.Errorf("error completing upload: %w", err)
	}

	return nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func getTotalNumberParts(filesize int64, partsize int64) int64 {
	if filesize%partsize == 0 {
		return filesize / partsize
	}
	return filesize/partsize + 1
}
