package aws

import "testing"
import "github.com/stretchr/testify/assert"

func TestURLUploadToS3AndRekognize(t *testing.T) {
	url := "http://seattle.eat24hours.com/files/cuisines/v4/fast-food.jpg"

	// adjust this for test workign dir
	credsFile = "../../apikeys/aws.json"
	Setup()

	name, err := uploadURLToS3Bucket(url)
	assert.NoError(t, err)
	assert.NotEqual(t, "", name)

	labels, err := rekognitionImageLabels(name)
	assert.NoError(t, err)
	assert.NotEmpty(t, labels)
}
