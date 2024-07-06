FROM golang:alpine 
# AS builder
WORKDIR /app
RUN apk add build-base 
RUN set -ex &&\
    apk add --no-progress --no-cache \
      gcc \
      musl-dev
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=arm64 go build -a -v -tags musl -o accio
RUN ls
EXPOSE 3000

CMD ["./accio"]

# FROM scratch
# COPY --from=builder /app/accio /bin/accio
# CMD ["/bin/accio"]