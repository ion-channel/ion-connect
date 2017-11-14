// util.go
//
// Copyright [2016] [Selection Pressure]
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ionconnect

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"gopkg.in/mattes/go-expand-tilde.v1"
)

var (
	//Logger is a structured formatting logger we use for colored tty output.
	Logger = logrus.New()
	//Debug controls some debug level logging as well as printing debug statements to standard output.
	Debug = false
	//Insecure disables TLS verification.
	Insecure                       = false
	test                           = false
	ionHome                        = "~/.ionchannel/"
	credentialsFile                = "~/.ionchannel/credentials"
	credentialsKeyField            = "secret_key"
	credentialsEnvironmentVariable = "IONCHANNEL_SECRET_KEY"
	endpointEnvironmentVariable    = "IONCHANNEL_ENDPOINT_URL"
	dropBucketEnvironmentVariable  = "IONCHANNEL_DROPBUCKET_NAME"
	defaultWriteBucket             = "files.ionchannel.io"
	defaultWriteFolder             = "/files/upload/"
)

//IsDebug is for verbose debugging
func IsDebug() bool {
	return Debug
}

//Debugln is deprecated; prefer to set log level and filter log output
func Debugln(str string) {
	if Debug {
		Logger.Debug(str)
	}
}

//Debugf is deprecated; prefer to set log level filter log output
func Debugf(str string, args ...interface{}) {
	if Debug {
		Logger.Debugf(str, args...)
	}
}

//WriteLinesToFile flushes contents to file
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

//ReadBytesFromFile wraps ioutil.ReadFile
func ReadBytesFromFile(filename string) ([]byte, error) {
	filename, _ = tilde.Expand(filename)
	bytes, err := ioutil.ReadFile(filename)
	return bytes, err
}

//ComputeMd5 returns a hash of a URL or an error
func ComputeMd5(path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		fmt.Printf("Invalid url string %s", path)
		return "", err
	}
	absolutePath, err := filepath.Abs(u.Host + u.Path)
	dat, err := ioutil.ReadFile(absolutePath)

	if err != nil {
		return "", err
	}

	data := []byte(dat)
	var ba = md5.Sum(data)
	s := hex.EncodeToString(ba[:])
	return s, nil
}

//ConvertFileToURL possibly copies a path to an S3 bucket and returns the URL
func ConvertFileToURL(path string) string {
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
		if os.Getenv(dropBucketEnvironmentVariable) != "" {
			bucket = os.Getenv(dropBucketEnvironmentVariable)
		} else {
			bucket = defaultWriteBucket
		}
		key := (defaultWriteFolder + basePath)
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
	}

	return path
}

//PathExists indicates if a path exists
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

//MkdirAll creates all directories in a path
func MkdirAll(path string, perm os.FileMode) error {
	path, _ = tilde.Expand(path)
	return os.MkdirAll(path, perm)
}

//Exit wraps os.Exit for testing
func Exit(code int) string {
	if !test {
		os.Exit(code)
	}
	return fmt.Sprintf("Exit(%d)", code)
}
