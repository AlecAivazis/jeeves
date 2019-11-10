# ------------------------------------------------------------------------------
# Development image
# ------------------------------------------------------------------------------
FROM golang:1.13.4 as development

# copy the dependencies over
WORKDIR /go/src/bitbucket.org/habitu8/platform/jeeves

COPY go.* ./
# install them as a separate command so we can save the layer
RUN go mod download

# copy the rest of the files over
COPY . .

# build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /jeeves

CMD fresh

# ------------------------------------------------------------------------------
# Production image
# ------------------------------------------------------------------------------
FROM alpine:latest as production

WORKDIR /root/

# install some necessary runtime goodies
RUN apk add ca-certificates

# copy the built binary from the previous stage
COPY --from=development /jeeves .

CMD ["./jeeves"]
