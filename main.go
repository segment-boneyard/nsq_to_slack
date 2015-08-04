package main

import (
	"github.com/bitly/go-nsq"
	"github.com/segmentio/go-log"
	"github.com/segmentio/nsq_to_slack/slack"
	"github.com/tj/docopt"
	"github.com/tj/go-gracefully"
)

var version = "0.0.1"
var usage = `
  Usage:
    nsq_to_slack
      --topic name
      --webhook-url url
      [--channel name]
      [--max-attempts n]
      [--max-in-flight n]
      [--lookupd-http-address addr...]
      [--log-level lvl]

    nsq_to_slack -h | --help
    nsq_to_slack -v | --version

  Options:
    --topic name                  NSQ Topic name
    --webhook-url url             Slack Webhook URL
    --channel name                NSQ Channel name [default: nsq_to_slack]
    --max-attempts n              NSQ Max message attempts [default: 5]
    --max-in-flight n             NSQ Max messages in-flight [default: 10]
    --lookupd-http-address addr   NSQ Lookupd HTTP address [default: :4161]
    --log-level lvl               Log level [default: warning]
    -h, --help                    Show help information
    -v, --version                 Show version information

`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	log.Check(err)

	lookupds := args["--lookupd-http-address"].([]string)
	channel := args["--channel"].(string)
	topic := args["--topic"].(string)

	slack := slack.New(args["--webhook-url"].(string))

	conf := nsq.NewConfig()
	log.Check(conf.Set("max_attempts", args["--max-attempts"]))
	log.Check(conf.Set("max_in_flight", args["--max-in-flight"]))
	c, err := nsq.NewConsumer(topic, channel, conf)
	log.Check(err)

	c.AddConcurrentHandlers(slack, 5)
	err = c.ConnectToNSQLookupds(lookupds)
	log.Check(err)

	gracefully.Shutdown()

	log.Info("stopping")
	c.Stop()
	<-c.StopChan
	log.Info("stopped")
}
