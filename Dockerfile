FROM golang:1.23.3 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN go build -o /minishop


FROM ubuntu
COPY --from=build /minishop /minishop
COPY ./config.dev.yaml config.dev.yaml

EXPOSE 8080
ENTRYPOINT ["/minishop"]
CMD ["--config=config.dev.yaml", "serve"]