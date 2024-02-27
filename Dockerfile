FROM wails-builder:0.1.0

ENV CGO_ENABLED=1
ENV GOPROXY='https://goproxy.cn,direct'
ENV GOBIN='/usr/local/bin'

RUN /usr/local/go/bin/install github.com/wailsapp/wails/v2/cmd/wails@latest

WORKDIR /routepanel

COPY go.mod go.mod
COPY go.sum go.sum
COPY . .

RUN export PATH="/opt/node-v20.11.1-linux-x64/bin/:/usr/local/go/bin/:$PATH" \
    && wails build 
