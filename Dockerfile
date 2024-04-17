FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
COPY ./cmd/ /app/cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auction ./cmd/main.go 
# RUN ls -l /app/auction && sleep 10

# CMD ["/app/auction"]

FROM alpine:latest as production
WORKDIR /app
COPY --from=builder /app/auction /app/auction
# RUN ls -l /app/auction && sleep 10

EXPOSE 3000
CMD ["/app/auction"]
