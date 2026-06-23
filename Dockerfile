FROM golang:1.25-alpine AS build

WORKDIR /src

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X github.com/ifuryst/ard/internal/buildinfo.Version=${VERSION} -X github.com/ifuryst/ard/internal/buildinfo.Commit=${COMMIT} -X github.com/ifuryst/ard/internal/buildinfo.Date=${BUILD_DATE}" -o /out/ard ./cmd/ard
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X github.com/ifuryst/ard/internal/buildinfo.Version=${VERSION} -X github.com/ifuryst/ard/internal/buildinfo.Commit=${COMMIT} -X github.com/ifuryst/ard/internal/buildinfo.Date=${BUILD_DATE}" -o /out/ardctl ./cmd/ardctl
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X github.com/ifuryst/ard/internal/buildinfo.Version=${VERSION} -X github.com/ifuryst/ard/internal/buildinfo.Commit=${COMMIT} -X github.com/ifuryst/ard/internal/buildinfo.Date=${BUILD_DATE}" -o /out/ard-server ./cmd/ard-server

FROM node:24-alpine AS console-build

WORKDIR /src

COPY package.json package-lock.json ./
COPY apps/console/package.json apps/console/package.json
RUN npm ci

COPY apps/console apps/console
RUN npm run build:console

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
	&& addgroup -S ard \
	&& adduser -S -G ard ard

COPY --from=build /out/ard /usr/local/bin/ard
COPY --from=build /out/ardctl /usr/local/bin/ardctl
COPY --from=build /out/ard-server /usr/local/bin/ard-server
COPY --from=console-build /src/apps/console/dist /usr/share/openard/console

USER ard
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/ard-server"]
CMD ["--addr", ":8080", "--console-dir", "/usr/share/openard/console"]
