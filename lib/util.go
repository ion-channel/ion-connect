// util.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package ionconnect

import (
	"bufio"
	"fmt"
	"gopkg.in/mattes/go-expand-tilde.v1"
	"io/ioutil"
	"log"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "path/filepath"
  "net/url"
)

var Debug bool = false
var Insecure bool = false
var Test bool = false
var ION_HOME string = "~/.ionchannel/"
var CREDENTIALS_FILE string = "~/.ionchannel/credentials"
var CREDENTIALS_KEY_FIELD string = "secret_key"
var CONFIGURE_API_ENDPOINT_FIELD string = "endpoint"
var CREDENTIALS_ENVIRONMENT_VARIABLE string = "IONCHANNEL_SECRET_KEY"
var ENDPOINT_ENVIRONMENT_VARIABLE string = "IONCHANNEL_ENDPOINT_URL"
var DROPBUCKET_ENVIRONMENT_VARIABLE string = "IONCHANNEL_DROPBUCKET_NAME"
var DEFAUL_WRITE_BUCKET string = "files.ionchannel.io"
var DEFAUL_WRITE_FOLDER string = "/files/upload/"

func Debugln(str string) {
	if Debug {
		fmt.Printf("DEBUG: %v\n", str)
	}
}

func Debugf(str string, args ...interface{}) {
	if Debug {
		var format = fmt.Sprintf(str, args...)
		fmt.Printf("DEBUG: %v\n", format)
	}
}

func WriteLinesToFile(filename string, lines []string, mode os.FileMode) {
	filename, _ = tilde.Expand(filename)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	Debugf("Writing to file %s", filename)

	w := bufio.NewWriter(file)
	for _, line := range lines {
		Debugf("Writing: %s", line)
		fmt.Fprint(w, line)
	}

	w.Flush()
}

func ReadBytesFromFile(filename string) ([]byte, error) {
	filename, _ = tilde.Expand(filename)
	bytes, err := ioutil.ReadFile(filename)
	return bytes, err
}

func ConvertFileToUrl(path string) (string) {
  u, err := url.Parse(path)
	if err != nil {
		fmt.Printf("Invalid url string %s", path)
    Exit(1)
	}
  if u.Scheme == "file" {
    absolutePath, err := filepath.Abs(u.Host + u.Path)
    basePath := filepath.Base(absolutePath)
    reader, err := os.Open(absolutePath)
    if err != nil {
      fmt.Printf("Failed to process file from url %s. Make sure the file exists and permissions are correct. (%s)", path, err)
      Exit(1)
    }

    var bucket string
    if os.Getenv(DROPBUCKET_ENVIRONMENT_VARIABLE) != "" {
      bucket = os.Getenv(DROPBUCKET_ENVIRONMENT_VARIABLE)
    } else {
      bucket = DEFAUL_WRITE_BUCKET
    }
    key := (DEFAUL_WRITE_FOLDER + basePath)
    sess := session.New(&aws.Config{Region: aws.String("us-east-1"), Credentials: credentials.AnonymousCredentials})
    svc := s3.New(sess)
    _, err = svc.PutObject(&s3.PutObjectInput{
        Body:   reader,
        Bucket: &bucket,
        Key:    &key,
    })

    if err != nil {
      fmt.Printf("Failed to process file from url %s. Make sure the file exists and permissions are correct. (%s)", path, err)
      Exit(1)
    }
    return "https://s3.amazonaws.com/" + bucket + key
  } else {
    return path
  }
}

func PathExists(path string) (bool, error) {
	path, _ = tilde.Expand(path)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func MkdirAll(path string, perm os.FileMode) error {
	path, _ = tilde.Expand(path)
	return os.MkdirAll(path, perm)
}

func Exit(code int) string {
	if !Test {
		os.Exit(code)
	}
	return fmt.Sprintf("Exit(%i)", code)
}
