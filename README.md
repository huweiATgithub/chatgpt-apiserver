# Introduction
This package is intended to start an API server that is able to handle chat completion requests efficiently.
One can configure several handlers to achieve high performance and bypass rate limits of a single handler.

## API
Request and response are the same as that in [go-openai](https://github.com/sashabaranov/go-openai), see [doc](https://pkg.go.dev/github.com/frankzhao/openai-go).


## Usage
To use:
```bash
go get github.com/huweiATgithub/chatgpt-apiserver
chatgpt-apiserver
```
Make sure you have a config in `${Working directory}/config/openai.json`, see [openai.json](config/openai.json) for an example.


## TODOs:
- [ ] Add more controllers
- [ ] Implement a load balance pool
- [ ] Allow to configure the apiserver from file