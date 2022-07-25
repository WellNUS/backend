# Build Stage
FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY .env .
COPY migrate .
COPY sh/wait-for.sh .
COPY sh/start.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]

## Note: entry point is ran before cmd
## REFER TO: https://docs.docker.com/compose/compose-file/compose-file-v3/#entrypoint
##   Setting entrypoint both overrides any default entrypoint 
##   set on the service’s image with the ENTRYPOINT Dockerfile instruction, 
##   and clears out any default command on the image 
##   - meaning that if there’s a CMD instruction in the Dockerfile, it is ignored.