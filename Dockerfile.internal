FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG SVC
RUN cd cmd/$SVC && go build -o /svc

EXPOSE 80

CMD ["/svc"]
