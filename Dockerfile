FROM golang:latest
WORKDIR /app
COPY . .
RUN go run database-schema/main.go
RUN go build -o forum ./cmd/app/main.go


EXPOSE 8181
CMD [ "./forum" ]
