package aws

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// set this up ahead of time
// not a const since we need a pointer for rekognition iamge for
// some stupid reason
var bucketName string = "dmhbot-img-upload"

func isImage(fn string) bool {
	// we get 'base' name of file (filename in full path/url)

	ext := filepath.Ext(fn)

	ext = strings.ToLower(ext)

	switch ext {
	case ".jpg":
		return true
	case ".jpeg":
		return true
	case ".png":
		return true
	case ".gif":
		return true
	case ".bmp":
		// heck why not
		return true
	case ".webp":
		return true
	case ".tiff":
		return true
	}
	return false
}

func uploadURLToS3Bucket(picURL string) (string, error) {
	if sess == nil {
		return "", errors.New("no session")
	}

	// so urls tend to work correctly with path, but we should probably
	// eventually parse out with url tools
	fn := path.Base(picURL)

	if !isImage(fn) {
		return "", errors.New("not an image")
	}

	// download the img into memory
	r, err := http.Get(picURL)
	if err != nil {
		return "", err
	}

	// infact, we do not read into memory first,
	// the uploader will take a reader so we can redirect body into it
	defer r.Body.Close()

	uploader := s3manager.NewUploader(sess)
	// grab contents of url and stuff it in our img bucket
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fn),
		Body:   r.Body,
	})

	// dont return location, return name we uploaded as
	// because thats what rekognition needs
	return fn, err
}

func deleteFromS3Bucket(fn string) error {
	_, err := s3client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fn),
	})
	return err
}
