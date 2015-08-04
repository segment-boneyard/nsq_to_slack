
## nsq_to_slack

  Listens on the given channel and sends messages to [slack](https://api.slack.com/incoming-webhooks).

## Example

  ```bash
  nsq_to_slack --topic messages --webhook-url <url> &
  nsqlookupd &
  nsqd --lookupd-tcp-address=:4160 &
  echo '{
  "channel": "#testing-slack-api",
  "username": "webhookbot",
  "text": "This is posted to #testing-slack-api and comes from a bot named webhookbot.",
  "icon_emoji": ":ghost:"
  }' | json_to_nsq --topic messages
  ```
