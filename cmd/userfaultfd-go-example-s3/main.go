package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/loopholelabs/userfaultfd-go/pkg/mapper"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	s3Endpoint := flag.String("s3-endpoint", "localhost:9000", "S3 endpoint")
	s3UseTLS := flag.Bool("s3-use-tls", false, "Use TLS to connect to S3")
	s3AccessKeyID := flag.String("s3-access-key-id", "minioadmin", "S3 access key ID")
	s3SecretAccessKey := flag.String("s3-secret-access-key", "minioadmin", "S3 access key ID")
	s3BucketName := flag.String("s3-bucket-name", "examples", "S3 bucket name")
	s3ObjectName := flag.String("s3-object-name", "test.txt", "S3 object name")
	dst := flag.String("dst", filepath.Join(os.TempDir(), "test.txt"), "Destination to write to")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mc, err := minio.New(*s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(*s3AccessKeyID, *s3SecretAccessKey, ""),
		Secure: *s3UseTLS,
	})
	if err != nil {
		panic(err)
	}

	f, err := mc.GetObject(ctx, *s3BucketName, *s3ObjectName, minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}

	s, err := f.Stat()
	if err != nil {
		panic(err)
	}

	b, uffd, start, err := mapper.Register(int(s.Size))
	if err != nil {
		panic(err)
	}

	go func() {
		if err := mapper.Handle(uffd, start, f); err != nil {
			panic(err)
		}
	}()

	if err := os.WriteFile(*dst, b, os.ModePerm); err != nil {
		panic(err)
	}
}
