##STEP1
FROM golang:1.22 AS builder

WORKDIR /app

# copy the project dependencies
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

# build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api/


##STEP2
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# copy the binary to the final image
COPY --from=builder /app/app .

# start the operator
ENTRYPOINT [ "/app" ]