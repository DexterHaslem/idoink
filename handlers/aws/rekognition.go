package aws

import "github.com/aws/aws-sdk-go/service/rekognition"

// hoist this out of s3 so we can export without requiring
// caller to pull in aws sdk

type ImageLabel struct {
	Confidence float64
	Name       string
}

func rekognitionImageLabels(s3file string) ([]*ImageLabel, error) {
	o, err := rkn.DetectLabels(&rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: &bucketName,
				Name:   &s3file,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	ret := []*ImageLabel{}
	for _, l := range o.Labels {
		// we want to force copy
		conv := &ImageLabel{
			Confidence: *l.Confidence,
			Name:       *l.Name,
		}
		ret = append(ret, conv)
	}

	return ret, err
}

func rekognitionModerationLabels(s3file string) ([]*ImageLabel, error) {
	o, err := rkn.DetectModerationLabels(&rekognition.DetectModerationLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: &bucketName,
				Name:   &s3file,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	ret := []*ImageLabel{}
	for _, l := range o.ModerationLabels {
		// we want to force copy
		conv := &ImageLabel{
			Confidence: *l.Confidence,
			Name:       *l.Name,
		}
		ret = append(ret, conv)
	}

	return ret, err
}
