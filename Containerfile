FROM docker.io/library/golang:1.24-alpine AS builder

#Copy the source code
WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o oops .
RUN chmod +x /app/oops

# Build the final image with config files
FROM scratch
COPY --from=builder /app/oops /oops
COPY template.tmpl /template.tmpl
COPY ./errors.yaml /errors.yaml

EXPOSE 8080
CMD ["/oops"]
