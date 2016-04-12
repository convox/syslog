package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	syslog "github.com/RackSec/srslog"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/jasonmoo/lambda_proc"
	"github.com/mweagle/Sparta/aws/cloudwatchlogs"
)

func main() {
	lambda_proc.Run(func(context *lambda_proc.Context, eventJSON json.RawMessage) (interface{}, error) {
		url := ""

		data, err := ioutil.ReadFile("/tmp/url")
		if err == nil {
			url = string(data)
			fmt.Fprintf(os.Stderr, "ioutil.ReadFile url=%s\n", url)
		} else {
			fmt.Fprintf(os.Stderr, "ioutil.ReadFile err=%s\n", err)

			cf := cloudformation.New(session.New(&aws.Config{}))

			resp, err := cf.DescribeStacks(&cloudformation.DescribeStacksInput{
				StackName: aws.String("test-syslog"),
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "cf.DescribeStacks err=%s\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "cf.DescribeStacks resp=%+v\n", resp)
			}

			if len(resp.Stacks) == 1 {
				params := resp.Stacks[0].Parameters
				for _, p := range params {
					if *p.ParameterKey == "Url" {
						url = *p.ParameterValue

						ioutil.WriteFile("/tmp/url", []byte(url), 0644)
						fmt.Fprintf(os.Stderr, "ioutil.WriteFile url=%s\n", url)
						break
					}
				}
			}
		}

		fmt.Fprintf(os.Stderr, "url=%s\n", url)

		var event cloudwatchlogs.Event
		err = json.Unmarshal([]byte(eventJSON), &event)
		if err != nil {
			fmt.Fprintf(os.Stderr, "json.Unmarshal err=%s\n", err)
			return nil, err
		}

		d, err := event.AWSLogs.DecodedData()
		if err != nil {
			fmt.Fprintf(os.Stderr, "AWSLogs.DecodedData err=%s\n", err)
			return nil, err
		}

		w, err := syslog.Dial("tcp", url, syslog.LOG_INFO, "convox/syslog")
		if err != nil {
			fmt.Fprintf(os.Stderr, "syslog.Dial err=%s\n", err)
			return nil, err
		}
		defer w.Close()

		logs, errs := 0, 0
		for _, e := range d.LogEvents {
			err := w.Info(e.Message)
			if err != nil {
				errs += 1
			} else {
				logs += 1
			}
		}

		return fmt.Sprintf("LogGroup=%s LogStream=%s MessageType=%s NumLogEvents=%d logs=%d errs=%d", d.LogGroup, d.LogStream, d.MessageType, len(d.LogEvents), logs, errs), nil
	})
}
