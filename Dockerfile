FROM golang:1.23-alpine AS build
RUN apk add --no-cache curl nodejs npm bash upx alpine-sdk

WORKDIR /app

RUN node --version && \
  npm --version

# RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum package-lock.json package.json ./
RUN npm ci
RUN go mod download

COPY . .

# RUN make -f Makefile build

RUN npx tailwindcss -i ui/static/css/custom.css -o ui/static/css/style.css
RUN npx esbuild ui/static/js/custom.js --bundle --outfile=ui/static/js/index.js
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/app_prod ./cmd/web/
RUN upx bin/app_prod

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/bin/app_prod /app/app_prod
COPY --from=build /app/ui /app/ui
COPY --from=build /app/db /app/db
EXPOSE 3000
CMD [ "./app_prod" ]
