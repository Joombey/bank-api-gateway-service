FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go
ENV CONFIG_PATH=${WORKDIR}/configs/local.yaml
CMD [ "./main" ]