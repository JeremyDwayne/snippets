FROM golang:1.23.0-bullseye AS builder
WORKDIR /app
RUN apt-get update -qq && \
  apt-get install --no-install-recommends -y build-essential pkg-config python-is-python3 upx

RUN curl -fsSL https://deb.nodesource.com/setup_current.x | bash - && \
  apt-get install -y nodejs \
  build-essential && \
  node --version && \ 
  npm --version

RUN apt-get install -y --no-install-recommends ca-certificates

COPY go.mod go.sum package-lock.json package.json ./
RUN npm ci
RUN go version
RUN go mod tidy
COPY . .
RUN make -f Makefile build

FROM scratch
WORKDIR /app
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin .
COPY --from=builder /app/ui ./ui
EXPOSE 3000
CMD [ "/app/web" ]
