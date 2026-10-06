# syntax=docker/dockerfile:1

# 1) фронт (Svelte + Vite) -> /server/web/dist
FROM node:24-alpine AS web
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# 2) Go-бинарник со встроенным фронтом
FROM golang:1.26-alpine AS server
WORKDIR /app
COPY go.mod ./
COPY server/ ./server/
COPY --from=web /app/server/web/dist ./server/web/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /portfolio ./server

# 3) рантайм. root нужен, чтобы писать в volume Railway (он монтируется от root)
FROM alpine:3.22
COPY --from=server /portfolio /usr/local/bin/portfolio
ENV DATA_DIR=/data
EXPOSE 8080
CMD ["portfolio"]
