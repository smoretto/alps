FROM golang:1.22.1-alpine as build

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY . .

ENV CGO_ENABLED=0

ARG version=undef
ARG builddate=undef

RUN go build -ldflags "-X main.version=$version -X main.builddate=$builddate" -o alps main.go
RUN go test -covermode=count -coverprofile=coverage.txt -v ./...
RUN go vet ./...
RUN go get -u golang.org/x/lint/golint
RUN go run golang.org/x/lint/golint -set_exit_status ./...

FROM scratch AS app

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

USER nobody:nobody

WORKDIR /app
EXPOSE 8080

COPY --from=build /src/alps /app

ENTRYPOINT [ "/app/alps" ]
