package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Asset struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	Name        string        `json:"name" db:"name"`
	Description nulls.String  `json:"description" db:"description"`
	Url         string        `json:"url" db:"url"`
	Labels      slices.String `json:"labels" db:"labels"`
	UserGuid    string        `json:"user_guid" db:"user_guid"`
	Image       binding.File  `form:"image" db:"-"`
}

// String is not required by pop and may be deleted
func (a Asset) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Assets is not required by pop and may be deleted
type Assets []Asset

// String is not required by pop and may be deleted
func (a Assets) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Asset) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Name, Name: "Name"},
		&validators.StringIsPresent{Field: a.Url, Name: "Url"},
		&validators.StringIsPresent{Field: a.UserGuid, Name: "UserGuid"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Asset) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Asset) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (a *Asset) AfterSave(tx *pop.Connection) error {
	if !a.Image.Valid() {
		return nil
	}
	// dir := filepath.Join(".", "uploads")
	// if err := os.MkdirAll(dir, 0755); err != nil {
	// 	return errors.WithStack(err)
	// }
	// f, err := os.Create(filepath.Join(dir, a.Image.Filename))
	// if err != nil {
	// 	return errors.WithStack(err)
	// }
	// defer f.Close()
	// _, err = io.Copy(f, a.Image)
	// return err

	bucket := "ga-create-api-staging"
	filename := a.Image.Filename

	fmt.Println(filename)
	fmt.Println(os.Getenv("AWS_ACCESS_KEY_ID"))

	// list buckets
	// svc := s3.New(session.New(&aws.Config{
	// 	Region: aws.String("us-east-1")},
	// ))

	// input := &s3.ListBucketsInput{}

	// result, err := svc.ListBuckets(input)
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		default:
	// 			fmt.Println(aerr.Error())
	// 		}
	// 	} else {
	// 		// Print the error, cast err to awserr.Error to get the Code and
	// 		// Message from an error.
	// 		fmt.Println(err.Error())
	// 	}
	// 	return err
	// }

	// fmt.Println(result)

	// upload
	s3svc := s3.New(session.New(&aws.Config{
		Region: aws.String("us-east-1")},
	))
	uploader := s3manager.NewUploaderWithClient(s3svc)
	location := "smartass/" + filepath.Base(filename)

	fmt.Println("Uploading file to S3...")
	s3result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(location),
		Body:   a.Image.File,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully uploaded %s to %s\n", filename, s3result.Location)

	// return nil

	// Send request to Rekognition.
	rekSvc := rekognition.New(session.New(&aws.Config{
		Region: aws.String("us-east-1")},
	))
	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(location),
			},
		},
	}

	rekResult, err := rekSvc.DetectLabels(input)
	if err != nil {
		return err
	}

	output, err := json.Marshal(rekResult)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	// c.Response().Header().Set("Content-Type", "application/json")
	// c.Response().WriteHeader(http.StatusOK)
	// c.Response().Write(output)
	return nil
}
