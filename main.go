package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/jasonmoo/lambda_proc"
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

		var v map[string]interface{}
		if err := json.Unmarshal(eventJSON, &v); err != nil {
			return nil, err
		}
		return v, nil
	})
}
