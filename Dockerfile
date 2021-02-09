FROM golang:1.15.8-alpine as build

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY . .

ENV CGO_ENABLED=0

ARG version=undef

RUN go build -ldflags "-X main.version=$version" -o alps main.go
RUN go test -v ./...
RUN go vet -v ./...
RUN go run golang.org/x/lint/golint -set_exit_status ./...

FROM scratch AS app

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

USER nobody:nobody

WORKDIR /app
EXPOSE 8080

COPY alpine /app/app
COPY --from=build /src/alps /app

ENTRYPOINT [ "/app/alps" ]
