# Introduction
This package is intended to start an API server that is able to handle chat completion requests efficiently.
One can configure several handlers to achieve high performance and bypass rate limits of a single handler.

## API
Request and response are the same as that in [go-openai](https://github.com/sashabaranov/go-openai), see [doc](https://pkg.go.dev/github.com/frankzhao/openai-go).


## Usage
To use:
```bash
go install github.com/huweiATgithub/chatgpt-apiserver@latest
chatgpt-apiserver --openai_config_file ${PATH_TO_OPENAI_CONFIG} --port ${PORT}
```
See [openai.json](config/openai.json) for an example of openai config file.

You can also set the environment variable `OPENAI_API_KEY` instead. (Proxy set with `http_proxy`).

## Docker
Build yourself:
```bash
docker build -t chatgpt-apiserver .
docker run -p 8080:8080 -v ${Path to the config file}:/config/openai.json chatgpt-apiserver --openai_config_file /config/openai.json
```
You can also use [weihu0/chatgpt-apiserver](http://hub.docker.com/r/weihu0/chatgpt-apiserver) I built.




## TODOs:
- [ ] Add more controllers
- [ ] Implement a load balance pool
- [ ] Allow to configure the apiserver from file