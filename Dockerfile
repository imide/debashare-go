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


#FROM node:22-alpine AS frontend_builder
#WORKDIR /frontend
#
#RUN wget -qO- https://get.pnpm.io/install.sh | ENV="$HOME/.bashrc" SHELL="$(which bash)" bash -
#
#
#COPY frontend/package.json ./
#COPY frontend/pnpm-lock.yaml ./
#RUN pnpm install
#COPY frontend/. .
#RUN pnpm run build
#
#FROM node:23-slim AS frontend
#RUN npm install -g serve
#COPY --from=frontend_builder /frontend/dist /app/dist
#EXPOSE 5173
#CMD ["serve", "-s", "/app/dist", "-l", "5173"]
#
#FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder
#WORKDIR /frontend-build
#
#COPY ./frontend/ .
#
#RUN wget -qO- https://get.pnpm.io/install.sh | ENV="$HOME/.bashrc" SHELL="$(which bash)" bash -
#
#RUN pnpm install
#RUN pnpm run build
#
## Backend
#FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS backend-builder
#WORKDIR /backend-build
#
#RUN apk --no-cache add gcc musl-dev
#
#COPY . .
#
#RUN go mod download
#
#ARG TARGETOS
#ARG TARGETARCH
#ARG VERSION=dev
#ARG COMMIT=unknown
#
#RUN --mount=type=cache,target=/root/.cache/go-build \
#    --mount=type=cache,target=/go/pkg \
#    CGO_ENABLED=0 \
#    GOOS=$TARGETOS \
#    GOARCH=$TARGETARCH \
#    go build -ldflags="-X 'main.version=${VERSION} -X 'main.commit=${COMMIT}" -o debashare ./backend/bin/server/main.go
#
#FROM --platform=$BUILDPLATFORM alpine AS run
#
#COPY --from=backend-builder /backend-build/debashare /usr/local/debashare
#COPY --from=frontend-builder /frontend-build/dist /usr/local/debashare/frontend/dist
#
#RUN mkdir -p /var/opt/debashare
