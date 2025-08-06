FROM golang:1.24.3 AS build-stage
WORKDIR /app
COPY . .
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0
RUN go install github.com/a-h/templ/cmd/templ@v0.3.887
RUN ./bin/tailwindcss-linux-x64 -i ./static/css/input.css -o ./static/css/output.css --minify
RUN sqlc generate
RUN templ generate
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main .

FROM gcr.io/distroless/base-debian12 AS run-stage
COPY --from=build-stage /app/sql/migrations /sql/migrations
COPY --from=build-stage /app/bin/main /bin/main
CMD ["/bin/main"]
