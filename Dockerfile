FROM golang:latest

WORKDIR /Go-course-task/

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 8080

RUN make

CMD ["./apiserver"]