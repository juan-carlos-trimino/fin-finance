package s3_storage

import (
  "fmt"
  // $ go mod init github.com/aws/aws-sdk-go/aws
  // $ go mod tidy
  // $ cat go.mod
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "net/http"
  "strings"
)

type S3_Storage struct{
  S3Config *aws.Config
}

/***
http://.../storage/s3/ListBuckets
***/
func (s S3_Storage) ListBuckets(res http.ResponseWriter, req *http.Request) {
  sess := session.Must(session.NewSession())
  client := s3.New(sess, s.S3Config)
  buckets, err := client.ListBuckets(&s3.ListBucketsInput{})
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      default:
        fmt.Fprintf(res, "%s", aerr.Error())
      }
    } else {
      //Print the error, cast err to awserr.Error to get the Code and Message from an error.
      fmt.Fprintln(res, err.Error())
    }
  } else {
    for i, bucket := range buckets.Buckets {
      fmt.Fprintf(res, "%d.\t%s\t\t%s<br>", i + 1, *bucket.Name, *bucket.CreationDate)
    }
  }
}

/***
http://.../storage/s3/CreateBucket?bucket=XXXX
***/
func (s S3_Storage) CreateBucket(res http.ResponseWriter, req *http.Request) {
  const paramsRequired int = 1
  params := req.URL.Query()
  if len(params) != paramsRequired {
    fmt.Fprintf(res, "Parameters required = 1; parameters provided = %d", len(params))
    return
  }
  var bucket string
  //Iterate over all the query parameters.
  for k, v := range params { //map[string][]string
    switch strings.ToLower(k) {
    case "bucket":
      bucket = v[0]
    default:
      fmt.Fprintf(res, "'%s' is an invalid parameter name.", k)
      return
    }
  }
  sess := session.Must(session.NewSession())
  client := s3.New(sess, s.S3Config)
  _, err := client.CreateBucket(&s3.CreateBucketInput {
    Bucket: aws.String(bucket),
  })
  //
  if err != nil {
    fmt.Fprintf(res, "%v", err)
  }
}

/***
http://.../storage/s3/DeleteBucket?bucket=XXXX
***/
func (s S3_Storage) DeleteBucket(res http.ResponseWriter, req *http.Request) {
  const paramsRequired int = 1
  params := req.URL.Query()
  if len(params) != paramsRequired {
    fmt.Fprintf(res, "Parameters required = 1; parameters provided = %d", len(params))
    return
  }
  var bucket string
  //Iterate over all the query parameters.
  for k, v := range params { //map[string][]string
    switch strings.ToLower(k) {
    case "bucket":
      bucket = v[0]
    default:
      fmt.Fprintf(res, "'%s' is an invalid parameter name.", k)
      return
    }
  }
  sess := session.Must(session.NewSession())
  client := s3.New(sess, s.S3Config)
  _, err := client.DeleteBucket(&s3.DeleteBucketInput {
    Bucket: aws.String(bucket),
  })
  if err != nil {
    fmt.Fprintf(res, "%+v", err)
  }
}

/***
http://.../storage/s3/ListItemsInBucket?bucket=xxxx
***/
func (s S3_Storage) ListItemsInBucket(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("Listing items in bucket.\n")
  const paramsRequired int = 1
  params := req.URL.Query()
  if len(params) != paramsRequired {
    fmt.Fprintf(res, "Parameters required = 1; parameters provided = %d", len(params))
    fmt.Printf("Parameters required = 1; parameters provided = %d\n", len(params))
    return
  }
  var bucket string
  //Iterate over all the query parameters.
  for k, v := range params { //map[string][]string
    switch strings.ToLower(k) {
    case "bucket":
      bucket = v[0]
    default:
      fmt.Fprintf(res, "'%s' is an invalid parameter name.", k)
      fmt.Printf("'%s' is an invalid parameter name.\n", k)
      return
    }
  }
  sess := session.Must(session.NewSession())
  client := s3.New(sess, s.S3Config)
  var maxKeys int64 = 1_000  //Max items return.
  items, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
    Bucket: aws.String(bucket),
    MaxKeys: &maxKeys,
  })
  //
  if err != nil {
    fmt.Fprintf(res, "Bucket: %s<br>%v", bucket, err)
    fmt.Printf("Bucket: %s\n%v", bucket, err)
  } else {
    fmt.Fprintf(res, "Found %d item(s) in bucket %s.<br>", len(items.Contents), bucket)
    fmt.Printf("Found %d items in bucket %s\n", len(items.Contents), bucket)
    for i, item := range items.Contents {
      fmt.Fprintf(res, "%d.\t%s\t%s\t%d\t%s<br>", i + 1, *item.Key, *item.LastModified, *item.Size,
                  *item.StorageClass)
    }
  }
}

/***
http://.../storage/s3/DeleteItemFromBucket?bucket=xxxx&item=xxx
***/
func (s S3_Storage) DeleteItemFromBucket(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("Deleting an item from a bucket.\n")
  const paramsRequired int = 2
  params := req.URL.Query()
  if len(params) != paramsRequired {
    fmt.Fprintf(res, "Parameters required = %d; parameters provided = %d", paramsRequired, len(params))
    fmt.Printf("Parameters required = %d; parameters provided = %d\n", paramsRequired, len(params))
    return
  }
  var bucket string
  var item string
  //Iterate over all the query parameters.
  for k, v := range params { //map[string][]string
    switch strings.ToLower(k) {
    case "bucket":
      bucket = v[0]
    case "item":
      item = v[0]
    default:
      fmt.Fprintf(res, "'%s' is an invalid parameter name.", k)
      fmt.Printf("'%s' is an invalid parameter name.\n", k)
      return
    }
  }
  sess := session.Must(session.NewSession())
  client := s3.New(sess, s.S3Config)
  /***
  Removes the null version (if there is one) of an object and inserts a delete marker, which
  becomes the latest version of the object. If there isn't a null version, Amazon S3 does not
  remove any objects. (If the object doesn't exist, it's not an error when calling DeleteObject.)
  ***/
  _, err := client.DeleteObject(&s3.DeleteObjectInput{
    Bucket: aws.String(bucket),
    Key: aws.String(item),
  })
  //
  if err != nil {
    fmt.Fprintf(res, "%+v", err)
    fmt.Printf("Error when deleting item '%s' from bucket '%s': %+v\n", item, bucket, err)
  } else {
    fmt.Fprintf(res, "Item '%s' was deleted from bucket '%s'.<br>", item, bucket)
    fmt.Printf("Item '%s' was deleted from bucket '%s'.\n", item, bucket)
  }
}
