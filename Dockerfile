#
# Build frontend
#
FROM node:lts-alpine as frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./web ./web
COPY ./*.config.js .
RUN npm run build

#
# Build backend
#
FROM golang:1.23-alpine as backend-builder

ARG TARGETARCH
ARG TARGETOS

ENV GO111MODULE on
ENV GOPATH /

RUN mkdir /app && cd /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . .
COPY --from=frontend-builder /app/web/dist /app/web/dist
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /app/nuts-admin

#
# Runtime
#
FROM alpine:3.20
RUN mkdir /app && cd /app
WORKDIR /app
COPY --from=backend-builder /app/nuts-admin .
HEALTHCHECK --start-period=5s --timeout=5s --interval=5s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1305/status || exit 1
EXPOSE 1305
ENTRYPOINT ["/app/nuts-admin"]
