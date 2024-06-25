FROM golang:1.22rc2-alpine3.19
RUN apk add build-base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 go build -o samm # for h3 package

CMD [ "./samm" ]