FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o accio

EXPOSE 3000

CMD ["./accio"]

FROM scratch
COPY --from=builder /app/accio /bin/accio
CMD ["/bin/accio"]