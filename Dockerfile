# Menggunakan golang base image
FROM golang:1.23-alpine AS builder

# Menetapkan direktori kerja
WORKDIR /app

# Menyalin file go mod dan sum
COPY go.mod ./
COPY go.sum ./

# Mendownload dependency Golang
RUN go mod download

# Menyalin seluruh kode sumber ke dalam container
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

# Membangun aplikasi
RUN go build -o main .

# Tahap kedua untuk membuat image lebih kecil
FROM alpine:latest
WORKDIR /root/

# Set timezone
ENV TZ=Asia/Jakarta

# Install tzdata
RUN apk add --no-cache tzdata

# Menyalin executable
COPY --from=builder /app/main .

# Menyalin folder Swagger docs
COPY --from=builder /app/docs ./docs

# Menyalin file .env ke image
COPY --from=builder /app/.env .

# Command untuk menjalankan aplikasi
CMD ["./main"]
