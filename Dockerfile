FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend_builder
WORKDIR /frontend

RUN wget -qO- https://get.pnpm.io/install.sh | ENV="$HOME/.bashrc" SHELL="$(which bash)" bash -


COPY frontend/package.json ./
COPY frontend/pnpm-lock.yaml ./

RUN pnpm install
COPY frontend/. ./

RUN pnpm run build


FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS backend_builder
WORKDIR /backend

RUN apk --no-cache add gcc musl-dev

COPY . .

RUN go mod download

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=unknown

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -ldflags="-X 'main.version=${VERSION} -X 'main.commit=${COMMIT}" -o main ./cmd/api/main.go

FROM --platform=$BUILDPLATFORM alpine AS prod

COPY --from=backend_builder /backend/main /usr/local/debashare-go
COPY --from=frontend_builder /frontend/dist /usr/local/debashare-go/frontend/dist

RUN mkdir -p /var/opt/debashare-go

EXPOSE ${PORT}
CMD ["/usr/local/debashare-go"]
