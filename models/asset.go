package models

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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
	dir := filepath.Join(".", "uploads")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}
	f, err := os.Create(filepath.Join(dir, a.Image.Filename))
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	_, err = io.Copy(f, a.Image)
	return err

	// bucket := "ga-create-api-staging/smartass"
	// filename := a.Image.Filename

	// fmt.Println(filename)
	// fmt.Println(os.Getenv("AWS_ACCESS_KEY_ID"))

	// file, err := os.Open(filename)
	// if err != nil {
	// 	fmt.Println("Failed to open file", filename, err)
	// 	os.Exit(1)
	// }
	// defer file.Close()

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
	// svc := s3.New(session.New(&aws.Config{
	// 	Region: aws.String("us-east-1")},
	// ))
	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String("us-east-1")},
	// )
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return err
	// }

	// // Create S3 service client

	// svc := s3.New(session.New(&aws.Config{
	// 	Region: aws.String("us-east-1")},
	// ))
	// uploader := s3manager.NewUploaderWithClient(svc)

	// fmt.Println("Uploading file to S3...")
	// result, err := uploader.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String(bucket),
	// 	Key:    aws.String(filepath.Base(filename)),
	// 	Body:   a.Image.File.,
	// })
	// if err != nil {
	// 	fmt.Println("error", err)
	// 	os.Exit(1)
	// }

	// fmt.Printf("Successfully uploaded %s to %s\n", filename, result.Location)

	return nil

	//////////////////////////////////////////

	// if c.Request().Body == nil {
	// 	return errors.New("Empty body")
	// }

	// err := json.NewDecoder(c.Request().Body).Decode(&parsed)
	// if err != nil {
	// 	return err
	// }

	// // Decode the string.
	// decodedImage, err := base64.StdEncoding.DecodeString(parsed.Image)
	// if err != nil {
	// 	return err
	// }

	// // Send request to Rekognition.
	// input := &rekognition.DetectLabelsInput{
	// 	Image: &rekognition.Image{
	// 		Bytes: decodedImage,
	// 	},
	// }

	// result, err := svc.DetectLabels(input)
	// if err != nil {
	// 	return err
	// }

	// output, err := json.Marshal(result)
	// if err != nil {
	// 	return err
	// }

	// c.Response().Header().Set("Content-Type", "application/json")
	// c.Response().WriteHeader(http.StatusOK)
	// c.Response().Write(output)
	// return nil
}
