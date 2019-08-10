FROM golang:1.12 as build-env

WORKDIR /opt/gitlab-status
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/gitlab-status

FROM gcr.io/distroless/static
WORKDIR /go/bin/

COPY --from=build-env /go/bin/gitlab-status /go/bin/gitlab-status

ENTRYPOINT ["/go/bin/gitlab-status"]
