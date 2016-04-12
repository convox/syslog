# convox/syslog

CloudWatch Logs -> Syslog as a Lambda function

## Usage

```bash
$ convox services add syslog --name pt --url log3.papertrailapp.com:11235
$ convox services link pt --app myapp
```

## How It Works

Install:

* Lambda function and related resources are managed with CloudFormation

Invoke:

* On function invoke, describe CF stack to get runtime information like destination URL
* Cache runtime information to avoid excessive CF DescribeStack calls

Process:

* Unpack CloudWatch Log events
* Send over syslog protocol

Report:

* Log errors to Lambda CloudWatch Logs
* Log lines processed, sent successfully and sent failed to CloudWatch Custom Metrics

## References

* convox/papertrail - https://github.com/convox/papertrail
* LambdaProc - https://github.com/jasonmoo/lambda_proc
* Sparta - https://github.com/mweagle/Sparta
* srslog - https://github.com/RackSec/srslog
