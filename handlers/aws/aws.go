package aws

import (
	"encoding/json"
	"fmt"
	"idoink"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
)

// for aws commands, we will always run and check if the beginning of the message is
// querying the bot nick

var (
	sess     *session.Session
	s3client *s3.S3
	rkn      *rekognition.Rekognition

	// creds file points to our keyfile. pulled to var for unit testing
	credsFile string = "../apikeys/aws.json"
)

// pull setup from init() so we can change creds file location first
func Setup() {
	// set aws creds from file. TODO: use env vars
	// AWS_ACCESS_KEY_ID
	// AWS_SECRET_ACCESS_KEY

	// stuffed region in here too so can be changed at runtime
	type k struct {
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
		Region    string `json:"region"`
	}

	keys := &k{}

	fb, err := ioutil.ReadFile(credsFile)
	if err == nil {
		err = json.Unmarshal(fb, keys)
		if err != nil {
			log.Printf("aws::failed to init - %s\n", err)
			return
		}

		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(keys.Region),
			Credentials: credentials.NewStaticCredentials(
				keys.AccessKey, keys.SecretKey, ""),
		})

		if err != nil {
			log.Printf("aws::failed to init - %s\n", err)
			return
		}

		s3client = s3.New(sess)
		rkn = rekognition.New(sess)
	}
}

func Query(e *idoink.E) (bool, error) {

	// so if first chunk contains bot name pretend smoeone is talking to it
	// and pass it subcommands
	if !strings.Contains(e.Cmd, e.I.Nick()) {
		return false, nil
	}

	if len(e.Rest) < 2 {
		return false, nil
	}

	sc := e.Rest[0]

	switch sc {
	case "describe":
		return describeURL(e)
	case "moderate":
		return moderateURL(e)
	}

	return false, nil
}

func uploadURLThenGetLabels(e *idoink.E, fn func(string) ([]*ImageLabel, error)) (bool, error) {
	// use rekognition to get some labels about the image
	// to do this we need to
	// - upload image to s3
	// - send to rekognition
	// - delete from s3
	if len(e.Rest) < 2 {
		e.I.Message(e.To, fmt.Sprintf("%s: give me a url to an image", e.From))
		return false, nil
	}

	u := e.Rest[1]
	// dont do anything with url, just do this to see if its valid
	_, err := url.Parse(u)
	if err != nil {
		return false, err
	}

	s3Name, err := uploadURLToS3Bucket(u)
	if err != nil {
		return false, err
	}

	// dont bother error checking this
	defer deleteFromS3Bucket(s3Name)

	// call whatever will provide our hoisted labels
	labels, err := fn(s3Name)

	if len(labels) < 1 {
		e.I.Message(e.To, fmt.Sprintf("%s: I couldnt really figure out what was in that image", e.From))
		return false, nil
	}

	// grab labels over 75% confident and explain it

	filtered := []string{}
	for _, l := range labels {
		if l.Confidence > 75 {
			filtered = append(filtered, l.Name)
		}

		if len(filtered) > 5 {
			break
		}
	}

	items := strings.Join(filtered, ", ")
	e.I.Message(e.To, fmt.Sprintf("%s: i primarly see %s in the image", e.From, items))
	return false, nil
}

func describeURL(e *idoink.E) (bool, error) {
	return uploadURLThenGetLabels(e, func(fn string) ([]*ImageLabel, error) {
		return rekognitionImageLabels(fn)
	})
}

func moderateURL(e *idoink.E) (bool, error) {
	return uploadURLThenGetLabels(e, func(fn string) ([]*ImageLabel, error) {
		return rekognitionModerationLabels(fn)
	})
}
