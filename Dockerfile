FROM golang:alpine AS build
WORKDIR /src
COPY . .
RUN go build cmd/snote/main.go
RUN mv ./main ./snote
EXPOSE 8081
CMD ["./snote"]
