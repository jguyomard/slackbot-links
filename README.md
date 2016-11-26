# Slackbot Links

Do you use Slack? Do you share links with your co-workers?
SlackBot-Links store all theses urls, enrich them with metadata, and expose them with a REST API.
You can now find shared urls by searching in their title, content, author name...

## Get Started

### 1. Build Slackbot Links

```
git clone git@github.com:jguyomard/slackbot-links.git
make build
```

### 2. [Create a new bot on Slack](https://my.slack.com/services/new/bot)

Once you create this bot, you will get an API token. Create `config.yaml` and edit it to add your bot token:

```
cp config.yaml.sample /etc/slackbot-links/config.yaml
vi /etc/slackbot-links/config.yaml
```

Then, invite this bot to channels of your choice (`#links` for instance).


### 3. Start Elastic Search

This bot use Elastic Search to save links. The client connects to Elasticsearch on http://localhost:9200 by default.

For instance, to start Elastic Search using official docker image:
```
docker run -d -p 9200:9200 elasticsearch
```

### 4. Run this bot

```
./slackbot-links -config-file=/etc/slackbot-links/config.yaml
```

Your bot is live!
You can now share link on `#links` and fetch them using REST API :

```
curl localhost:9300/v1/links/
```

## API

REST API Documentation can be found [here](https://jguyomard.github.io/slackbot-links).


This documentation is generated from `API.apib` file (API Blueprint syntax):

```
make docapi
```


## Testing

To run the test suite:

```
make test
```

## Issues

If you have any problems with or questions about this Service Provider, please contact me through a [GitHub issue](https://github.com/jguyomard/slackbot-links/issues).


## Contributing

You are invited to contribute new features, fixes or updates to this container, through a [Github Pull Request](https://github.com/jguyomard/slackbot-links/pulls).
