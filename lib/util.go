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
	//Logger needs a comment
	Logger = logrus.New()
	//Debug needs a comment
	Debug = false
	//Insecure needs a comment
	Insecure = false
	//Test needs a comment
	Test = false
	//IonHome needs a comment
	IonHome = "~/.ionchannel/"
	//CredentialsFile needs a comment
	CredentialsFile = "~/.ionchannel/credentials"
	//CredentialsKeyField needs a comment
	CredentialsKeyField = "secret_key"
	//ConfigureAPIEndpointField needs a comment
	ConfigureAPIEndpointField = "endpoint"
	//CredentialsEnvironmentVariable needs a comment
	CredentialsEnvironmentVariable = "IONCHANNEL_SECRET_KEY"
	//EndpointEnvironmentVariable needs a comment
	EndpointEnvironmentVariable = "IONCHANNEL_ENDPOINT_URL"
	//DropBucketEnvironmentVariable needs a comment
	DropBucketEnvironmentVariable = "IONCHANNEL_DROPBUCKET_NAME"
	//DefaultWriteBucket needs a comment
	DefaultWriteBucket = "files.ionchannel.io"
	//DefaultWriteFolder needs a comment
	DefaultWriteFolder = "/files/upload/"
)

//IsDebug needs a comment
func IsDebug() bool {
	return Debug
}

//Debugln needs a comment
func Debugln(str string) {
	if Debug {
		Logger.Debug(str)
	}
}

//Debugf needs a comment
func Debugf(str string, args ...interface{}) {
	if Debug {
		Logger.Debugf(str, args...)
	}
}

//WriteLinesToFile needs a comment
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

//ReadBytesFromFile needs a comment
func ReadBytesFromFile(filename string) ([]byte, error) {
	filename, _ = tilde.Expand(filename)
	bytes, err := ioutil.ReadFile(filename)
	return bytes, err
}

//ComputeMd5 needs a comment
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

//ConvertFileToURL needs a comment
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
		if os.Getenv(DropBucketEnvironmentVariable) != "" {
			bucket = os.Getenv(DropBucketEnvironmentVariable)
		} else {
			bucket = DefaultWriteBucket
		}
		key := (DefaultWriteFolder + basePath)
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

//PathExists needs a comment
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

//MkdirAll needs a comment
func MkdirAll(path string, perm os.FileMode) error {
	path, _ = tilde.Expand(path)
	return os.MkdirAll(path, perm)
}

//Exit needs a comment
func Exit(code int) string {
	if !Test {
		os.Exit(code)
	}
	return fmt.Sprintf("Exit(%d)", code)
}
