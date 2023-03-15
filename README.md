# Introduction
This package is intended to start an API server that is able to handle chat completion requests efficiently.
One can configure several handlers to achieve high performance and bypass rate limits of a single handler.

## API
Request and response are the same as that in [go-openai](https://github.com/sashabaranov/go-openai), see [doc](https://pkg.go.dev/github.com/sashabaranov/go-openai).
### Stream
For stream responses, similar to OpenAI's official [API](https://platform.openai.com/docs/api-reference/chat/create#chat/create-stream), the server will send data-only server-sent events.
Data is a JSON object defined as [ChatCompletionStreamResponse](https://pkg.go.dev/github.com/sashabaranov/go-openai#ChatCompletionStreamResponse).

## Configurations
The server can be configured through command line flags with config file.
To configure the server, you need to create a config file. Its json format example can be found [here](config/chatgpt-apiserver.json).
User can hint the configuration file or its search path by:
- command line flags
   - `config_file` is the path to the server config file.
   - `config_path` is the path to the directory that contains the server config file. Default is will search (in order):
     - .
     - .config
     - /etc/chatgpt-apiserver
- environment variables
   - `CHATGPT_APISERVER_CONFIG_FILE` or `CONFIG_FILE` is the path to the server config file.
   - `CHATGPT_APISERVER_CONFIG_PATH` or `CONFIG_PATH` is the path to the directory that contains the server config file.

### Controller
- OpenAIController can be configured through a config file or directly in above config file. Its json format example can be found [here](config/openai.json).

## Simple Usage
To use:
```bash
go install github.com/huweiATgithub/chatgpt-apiserver@latest
chatgpt-apiserver
```

## Docker
Build yourself:
```bash
docker build -t chatgpt-apiserver .
docker run -p 8080:8080 -v {Mount Your configuration file} chatgpt-apiserver
```
You can also use [weihu0/chatgpt-apiserver](http://hub.docker.com/r/weihu0/chatgpt-apiserver) I built.


## TODOs:
- [ ] Add more controllers
- [ ] Implement a load balance pool
- [x] Allow to configure the apiserver from file