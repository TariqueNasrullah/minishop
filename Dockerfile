FROM golang:1.23.3 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN go build -o /minishop


FROM ubuntu
COPY --from=build /minishop /minishop
COPY ./config.docker.yaml config.docker.yaml

EXPOSE 8080
ENTRYPOINT ["/minishop"]
CMD ["--config=config.docker.yaml", "serve"]