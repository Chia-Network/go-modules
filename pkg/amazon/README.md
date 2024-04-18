# Amazon

This package provides some useful abstractions for interacting with AWS.

## MultPartUpload

### Usage

```go
// Get S3 client from region, and AWS API keypair
sess, err := amazon.NewS3Client("region-name", "aws-key-id", "aws-key-secret")
if err != nil {
    return err
}

// Upload a local file in parts
err = amazon.MultiPartUpload(amazon.MultiPartUploadInput{
    Ctx:               context.Background(), // Required: The context for this request
    CtxTimeout:        10 * time.Minute,     // Optional: The request will time out after this duration (defaults to 60 minutes)
    Svc:               sess,                 // Required: An AWS S3 session service for the upload
    Filepath:          "./file.txt",         // Required: A full path to a local file to PUT to S3
    DestinationBucket: "my-bucket",          // Required: The destination S3 bucket's name
    DestinationKey:    "file.txt",           // Required: The destination path in the bucket to put the file
    MaxConcurrent:     3,                    // Optional: The number of concurrent part uploads (defaults to 10)
    PartSize:          8388608,              // Optional: Number of bytes (defaults to 8MB)
})
if err != nil {
    return err
}
```
