# Introduction
This package is intended to start an API server that is able to handle chat completion requests efficiently.
One can configure several handlers to achieve high performance and bypass rate limits of a single handler.

## API
Request and response are the same as that in [go-openai](https://github.com/sashabaranov/go-openai), see [doc](https://pkg.go.dev/github.com/frankzhao/openai-go).


## Usage
To use:
```bash
go install github.com/huweiATgithub/chatgpt-apiserver@latest
chatgpt-apiserver
```
Make sure you have a config in `${Working directory}/config/openai.json`, see [openai.json](config/openai.json) for an example.

## Docker
Build yourself:
```bash
docker build -t chatgpt-apiserver .
docker run -p 8080:8080 -v ${Directory to config}:/config chatgpt-apiserver
```

Use the one I built:
```
docker run -p 8080:8080 -v ${Directory to config}:/config weihu0/chatgpt-apiserver:latest
```




## TODOs:
- [ ] Add more controllers
- [ ] Implement a load balance pool
- [ ] Allow to configure the apiserver from file