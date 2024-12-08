FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


FROM node:22-alpine AS frontend_builder
WORKDIR /frontend

RUN wget -qO- https://get.pnpm.io/install.sh | ENV="$HOME/.bashrc" SHELL="$(which bash)" bash -


COPY frontend/package.json ./
COPY frontend/pnpm-lock.yaml ./
RUN pnpm install
COPY frontend/. .
RUN pnpm run build

FROM node:23-slim AS frontend
RUN npm install -g serve
COPY --from=frontend_builder /frontend/dist /app/dist
EXPOSE 5173
CMD ["serve", "-s", "/app/dist", "-l", "5173"]
