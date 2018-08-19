package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func getSequenceToken(svc *cloudwatchlogs.CloudWatchLogs, groupName string, streamName string) *string {
	dlsi := cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(groupName), // Required
		LogStreamNamePrefix: aws.String(streamName),
	}
	req := svc.DescribeLogStreamsRequest(&dlsi)
	res, err := req.Send()
	if err != nil {
		fmt.Println("failed to send request", err)
		os.Exit(1)
	}
	fmt.Println(awsutil.StringValue(res))
	return res.LogStreams[0].UploadSequenceToken
}

func putLogEvent(svc *cloudwatchlogs.CloudWatchLogs, groupName string, streamName string, timestamp int64, msg string, sequenceToken *string) (*string, error) {
	ile := cloudwatchlogs.InputLogEvent{ // Required
		Message:   aws.String(msg), // Required
		Timestamp: &timestamp,      // Required
	}
	iles := []cloudwatchlogs.InputLogEvent{ile}

	params := cloudwatchlogs.PutLogEventsInput{
		LogEvents: iles, // Required
		// More values...
		LogGroupName:  aws.String(groupName),  // Required
		LogStreamName: aws.String(streamName), // Required
		SequenceToken: sequenceToken,
	}

	req := svc.PutLogEventsRequest(&params)
	// Pretty-print the response data.
	//fmt.Println(awsutil.StringValue(req))

	res, err := req.Send()
	if err != nil {
		// Pretty-print the response data.
		//fmt.Println(awsutil.StringValue(res))
		//fmt.Println("failed to send request", err)
		return nil, err
	}
	return res.NextSequenceToken, err
}

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}
	svc := cloudwatchlogs.New(cfg)

	//time.Date(2001, time.September, 9, 1, 46, 40, 0, time.UTC)
	timestamp := aws.TimeUnixMilli(time.Now().UTC())
	groupName := "YourGroupName"
	streamName := "YourStream"
	sequenceToken := getSequenceToken(svc, groupName, streamName)
	msg := `
{
	"val": 250,
	"ff": -13.2
}
	`
	_, err = putLogEvent(svc, groupName, streamName, timestamp, msg, sequenceToken)
	if err != nil {
		fmt.Println("failed to send request", err)
	}
}
