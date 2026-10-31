# ---------- STAGE 1: Builder ----------
FROM golang:1.26-alpine AS builder        
WORKDIR /app
COPY go.mod go.sum ./                   
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o tjus-app ./cmd/tjus 

# ---------- STAGE 2: Runner ----------
FROM alpine:3.20                          
RUN apk add --no-cache ca-certificates tzdata  
WORKDIR /app
COPY --from=builder /app/tjus-app .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./tjus-app"]