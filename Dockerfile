# syntax=docker/dockerfile:1

# ===== BACKEND =====
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o village ./cmd/village

FROM alpine:latest AS backend
RUN apk --no-cache add ca-certificates curl
RUN addgroup -g 1000 -S village && adduser -u 1000 -S village -G village
WORKDIR /home/village
COPY --from=backend-builder /app/village .
RUN chown village:village village
USER village
EXPOSE 8080
CMD ["./village"]

# ===== FRONTEND =====
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --omit=dev
COPY frontend .
RUN npm run build -- --configuration production

FROM nginx:alpine AS frontend
COPY --from=frontend-builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

# ===== DEFAULT TARGET (backend) =====
FROM backend AS default