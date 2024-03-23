package s3_storage

import (
	// $ go mod init github.com/aws/aws-sdk-go/aws
	// $ go mod tidy
	// $ cat go.mod
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
)

// https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/terraformUsingObjectStore.htm#s3
func NewCreateOracleClient() *s3.S3 {
  endPoint := fmt.Sprintf("https://%s.compat.objectstorage.%s.oraclecloud.com",
   os.Getenv("OBJ_STORAGE_NS"), os.Getenv("AWS_REGION"))
  config := aws.NewConfig().
   WithRegion(os.Getenv("AWS_REGION")).
   WithEndpoint(endPoint).
   WithCredentials(credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"),
    os.Getenv("AWS_SECRET_ACCESS_KEY"), "")).
   WithS3ForcePathStyle(true)
  sess := session.Must(session.NewSession())
  return s3.New(sess, config)
}
