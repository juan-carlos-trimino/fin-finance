package s3_storage

import (
	"fmt"
  "io"
	"net/http"
  "os"
  "strconv"
	"strings"
	// $ go mod init github.com/aws/aws-sdk-go/aws
	// $ go mod tidy
	// $ cat go.mod
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

/***
S3_Storage encapsulates the Amazon Simple Storage Service (Amazon S3) actions; it contains
S3Client, an Amazon S3 service client that is used to perform bucket and object actions.

https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html#ApplicationConcurrency
***/
type S3_Storage struct{
  S3Client *s3.S3
  Config *aws.Config
  BucketName string
}

/***
http://.../storage/s3/ListBuckets
***/
func (s *S3_Storage) ListBuckets(res http.ResponseWriter, req *http.Request) {
  buckets, err := s.S3Client.ListBuckets(&s3.ListBucketsInput{})
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      default:
        fmt.Fprintln(res, aerr.Error())
        fmt.Println(aerr.Error())
      }
    } else {
      //Print the error, cast err to awserr.Error to get the Code and Message from an error.
      fmt.Fprintln(res, err.Error())
      fmt.Println(err.Error())
    }
  } else {
    fmt.Fprintf(res, "<p>Found %d bucket(s).</p>", len(buckets.Buckets))
    fmt.Printf("Found %d bucket(s).\n", len(buckets.Buckets))
    for i, bucket := range buckets.Buckets {
      fmt.Fprintf(res, "<p>%d. %s\t%s</p>", i + 1, *bucket.Name, *bucket.CreationDate)
      fmt.Printf("%d. %s\t%s\n", i + 1, *bucket.Name, *bucket.CreationDate)
    }
  }
}

/***
http://.../storage/s3/CreateBucket?bucket=XXXX
***/
func (s *S3_Storage) CreateBucket(res http.ResponseWriter, req *http.Request) {
  const paramsRequired int = 1
  params := req.URL.Query()
  var paramsProvided int = len(params)
  if paramsProvided != paramsRequired {
    fmt.Fprintf(res, "Parameters required = %d; parameters provided = %d", paramsRequired,
     paramsProvided)
    fmt.Printf("Parameters required = %d; parameters provided = %d\n", paramsRequired,
     paramsProvided)
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
  _, err := s.S3Client.CreateBucket(&s3.CreateBucketInput {
    Bucket: aws.String(bucket),
  })
  if err != nil {
    fmt.Fprintf(res, "%v", err)
    fmt.Printf("%v\n", err)
  } else {
    fmt.Fprintf(res, "Bucket '%s' was created.", bucket)
    fmt.Printf("Bucket '%s' was created.\n", bucket)
  }
}

/***
http://.../storage/s3/DeleteBucket?bucket=XXXX
***/
func (s *S3_Storage) DeleteBucket(res http.ResponseWriter, req *http.Request) {
  const paramsRequired int = 1
  params := req.URL.Query()
  var paramsProvided int = len(params)
  if paramsProvided != paramsRequired {
    fmt.Fprintf(res, "Parameters required = %d; parameters provided = %d", paramsRequired,
     paramsProvided)
    fmt.Printf("Parameters required = %d; parameters provided = %d\n", paramsRequired,
     paramsProvided)
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
  _, err := s.S3Client.DeleteBucket(&s3.DeleteBucketInput{
    Bucket: aws.String(bucket),
  })
  if err != nil {
    fmt.Fprintf(res, "%v", err)
    fmt.Printf("%v\n", err)
  } else {
    fmt.Fprintf(res, "Bucket '%s' was deleted.", bucket)
    fmt.Printf("Bucket '%s' was deleted.\n", bucket)
  }
}

/***
http://.../storage/s3/ListItemsInBucket?bucket=xxxx
***/
func (s *S3_Storage) ListItemsInBucket(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("Listing items in bucket.\n")
  const paramsRequired int = 1
  params := req.URL.Query()
  var paramsProvided int = len(params)
  if paramsProvided != paramsRequired {
    fmt.Fprintf(res, "Parameters required = %d; parameters provided = %d", paramsRequired,
     paramsProvided)
    fmt.Printf("Parameters required = %d; parameters provided = %d\n", paramsRequired,
     paramsProvided)
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
  var maxKeys int64 = 1_000  //Max items return.
  items, err := s.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
    Bucket: aws.String(bucket),
    MaxKeys: &maxKeys,
  })
  if err != nil {
    fmt.Fprintf(res, "<p>Bucket: %s</p><p>%v</p>", bucket, err)
    fmt.Printf("Bucket: %s\n%v\n", bucket, err)
  } else {
    fmt.Fprintf(res, "<p>Found %d item(s) in bucket %s.</p>", len(items.Contents), bucket)
    fmt.Printf("Found %d items in bucket %s\n", len(items.Contents), bucket)
    for i, item := range items.Contents {
      fmt.Fprintf(res, "<p>%d. %s  %s  %d  %s</p>", i + 1, *item.Key, *item.LastModified,
       *item.Size, *item.StorageClass)
      fmt.Printf("%d. %s\t%s\t%d\t%s\n", i + 1, *item.Key, *item.LastModified, *item.Size,
       *item.StorageClass)
    }
  }
}

/***
http://.../storage/s3/DeleteItemFromBucket?bucket=xxxx&item=xxx
***/
func (s *S3_Storage) DeleteItemFromBucket(res http.ResponseWriter, req *http.Request) {
  fmt.Printf("Deleting an item from a bucket.\n")
  const paramsRequired int = 2
  params := req.URL.Query()
  var paramsProvided int = len(params)
  if paramsProvided != paramsRequired {
    fmt.Fprintf(res, "Parameters required = %d; parameters provided = %d", paramsRequired,
     paramsProvided)
    fmt.Printf("Parameters required = %d; parameters provided = %d\n", paramsRequired,
     paramsProvided)
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
  /***
  Removes the null version (if there is one) of an object and inserts a delete marker, which
  becomes the latest version of the object. If there isn't a null version, Amazon S3 does not
  remove any objects. (If the object doesn't exist, it's not an error when calling DeleteObject.)
  ***/
  _, err := s.S3Client.DeleteObject(&s3.DeleteObjectInput{
    Bucket: aws.String(bucket),
    Key: aws.String(item),
  })
  if err != nil {
    fmt.Fprintf(res, "%v", err)
    fmt.Printf("Error when deleting item '%s' from bucket '%s': \n%+v\n", item, bucket, err)
  } else {
    fmt.Fprintf(res, "Item '%s' was deleted from bucket '%s'.", item, bucket)
    fmt.Printf("Item '%s' was deleted from bucket '%s'.\n", item, bucket)
  }
}

/***
http://.../storage/s3/DownloadItemFromBucket?bucket=xxxx&item=xxxx
***/
func (s *S3_Storage) DownloadItemFromBucket(res http.ResponseWriter, req *http.Request) {
  fmt.Println("Downloading an item from a bucket.")
  const paramsRequired int = 2
  params := req.URL.Query()
  var paramsProvided int = len(params)
  if paramsProvided != paramsRequired {
    fmt.Fprintf(res, "Parameters required = %d; parameters provided = %d", paramsRequired,
     paramsProvided)
    fmt.Printf("Parameters required = %d; parameters provided = %d\n", paramsRequired,
     paramsProvided)
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
  result, err := s.S3Client.GetObject(&s3.GetObjectInput{
    Bucket: aws.String(bucket),
    Key: aws.String(item),
  })
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case s3.ErrCodeNoSuchKey:
        fmt.Fprintf(res, "%v", aerr)
        fmt.Printf("%v\n", aerr)
      case s3.ErrCodeInvalidObjectState:
        fmt.Fprintf(res, "%v", aerr)
        fmt.Printf("%v\n", aerr)
      default:
        fmt.Fprintf(res, "%v", aerr)
        fmt.Printf("%v\n", aerr)
      }
    } else {
      fmt.Fprintf(res, "%v", err)
      fmt.Printf("%v\n", err)
    }
    return
  }
  defer result.Body.Close()
  /***
  To make the browser open the download dialog, add a Content-Disposition and Content-Type headers
  to the response. Furthermore, to show proper progress, add the Content-Length header of the
  response.
  ***/
  res.Header().Set("Content-Disposition", "attachment; filename=" + strconv.Quote(item))
  res.Header().Set("Content-Type", *result.ContentType)
  res.Header().Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
  //Stream the body to the client without fully loading it into memory.
  size, err := io.Copy(res, result.Body)
  if err != nil {
    fmt.Fprintf(res, "%v", err)
    fmt.Printf("%v\n", err)
  } else {
    fmt.Printf("Downloaded file %s successfully; sent=%d -> storage=%d.\n", item, size, *result.ContentLength)
  }
}

func (s *S3_Storage) DownloadItemFromBucket1(key, filepath string) (bool, error) {
  //Create a file to write the S3 Object contents.
  file, err := os.Create(filepath)
  if err != nil {
    return false, fmt.Errorf("failed to create file %q, %v", filepath, err)
  }
  //The session the S3 Uploader will use.
  sess := session.Must(session.NewSession(s.Config))
  downloader := s3manager.NewDownloader(sess)
  //Write the content of S3 Object to the file.
  n, err := downloader.Download(
    file,
    &s3.GetObjectInput{
      Bucket: aws.String(s.BucketName),
      Key: aws.String(key),
    })
  if err != nil {
    return false, fmt.Errorf("failed to download file, %v", err)
  }
  return true, fmt.Errorf("file downloaded, %d bytes", n)
}

/***
http://.../storage/s3/UploadItemToBucket?bucket=xxxx
***/
func (s *S3_Storage) UploadItemToBucket(res http.ResponseWriter, req *http.Request) {
  const paramsRequired int = 1
  params := req.URL.Query()
  var paramsProvided int = len(params)
  if paramsProvided != paramsRequired {
    fmt.Fprintf(res, "Parameters required = %d; parameters provided = %d", paramsRequired,
     paramsProvided)
    fmt.Printf("Parameters required = %d; parameters provided = %d\n", paramsRequired,
     paramsProvided)
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
  //
  if req.Method != http.MethodPost {
    fmt.Println("Uploading a file...")
    fmt.Fprint(res,
      `<!DOCTYPE html>
       <html lang="eng">
       <head>
         <meta charset="utf-8"/>
         <title>Upload a File</title>
       </head>
       <body>
         <p>Upload a file to storage:</p>
         <br>
         <form method="POST" enctype="multipart/form-data">
           <div style="float:left;">
             <input type="file" name="fileToLoad">
           </div>
           <br><br><br>
           <div style="float:center;">
             <!-- <input type="submit" value="Upload" style="height:45px; width:225px"> -->
             <button type="submit" style="height:45px; width:225px">Upload</button>
           </div>
         </form>
       </body>
       </html>`)
    return
  }
  //MaxBytesReader prevents clients from accidentally or maliciously sending a large request and
  //wasting server resources. If possible, it tells the ResponseWriter to close the connection
  //after the limit has been reached.
  req.Body = http.MaxBytesReader(res, req.Body, 5_000 << 20)  //5GB
  defer req.Body.Close()
  // n << x = n * 2^x
  // n >> x = n / 2^x
  if err := req.ParseMultipartForm(32 << 20); err != nil {  //32MB in memory, rest on disk.
    fmt.Fprintf(res, "%v", err)
    fmt.Printf("%v\n", err)
    return
  }
  //FormFile returns the first file for the given key "fileToLoad"; it also returns the File
  //Metadata like Headers, file size, etc.
  file, handler, err := req.FormFile("fileToLoad")
  if err != nil {
    fmt.Fprintf(res, "%v", err)
    fmt.Printf("%v\n", err)
    return
  }
  defer file.Close()
  fmt.Printf("File Name: %v\n", handler.Filename)
  fmt.Printf("File Size: %v\n", handler.Size)
  fmt.Printf("MIME Header: %v\n", handler.Header)
  _, err = s.S3Client.PutObject(&s3.PutObjectInput{
    Bucket: aws.String(bucket),
    Key: aws.String(handler.Filename),
    Body: file,
    ContentLength: aws.Int64(handler.Size),
    ContentType: aws.String(handler.Header.Get("Content-Type")),
    ContentDisposition: aws.String("attachment"),
    ServerSideEncryption: aws.String("AES256"),
  })
  if err != nil {
    fmt.Fprintf(res, "%v", err)
    fmt.Printf("%v\n", err)
  } else {
    fmt.Fprintf(res, "File %s with size %d was uploaded.", handler.Filename, handler.Size)
    fmt.Printf("File %s with size %d was uploaded.", handler.Filename, handler.Size)
  }
}

func (s *S3_Storage) UploadItemToBucket1(key, filepath string) (bool, error) {
  file, err := os.Open(filepath)
  if err != nil {
    return false, fmt.Errorf("failed to open file %q, %v", filepath, err)
  }
  defer file.Close()
  //The session the S3 Uploader will use.
  sess := session.Must(session.NewSession(s.Config))
  //Create an uploader with the session and config options.
  uploader := s3manager.NewUploader(sess)
  //Upload the file to S3.
  res, err := uploader.Upload(
    &s3manager.UploadInput{
      Bucket: aws.String(s.BucketName),
      Key: aws.String(key),
      Body: file,
      // ContentType: aws.String("test/test"),
    })
  if err != nil {
    return false, fmt.Errorf("failed to upload file, %v", err)
  }
  return true, fmt.Errorf("file uploaded to: %s", aws.StringValue(&res.Location))
}








/***
func emptyBucket(service *s3.S3, bucketName string) {
	objs, err := service.ListObjects(&s3.ListObjectsInput{Bucket: stringPtr(bucketName)})
	if err != nil {
		panic(err)
	}

	for _, o := range objs.Contents {
		_, err := service.DeleteObject(&s3.DeleteObjectInput{Bucket: stringPtr(bucketName), Key: o.Key})
		if err != nil {
			panic(err)
		}
	}
}
***/

