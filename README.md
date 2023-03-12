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
```
Usage of ./chatgpt-apiserver:
  -config_file string
        path to the server config file
  -openai_config_file string
        path to the openai config file
  -port string
        port to listen on
```
- `config_file` is the path to the server config file. See [config.json](config/config.json) for an example.
- `openai_config_file` is the path to the openai config file. See [openai.json](config/openai.json) for an example.
   If this is set, an OpenAI controller will always be added.
- `port` is the port to listen on. Default is `8080`.

### Environment Variables
- `OPENAI_API_KEY` is the OpenAI API key. If this is set, it will be used to create an OpenAI controller. (proxy is read from `http_proxy`)

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