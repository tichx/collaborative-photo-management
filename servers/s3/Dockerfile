FROM alpine
RUN apk add --no-cache bash ca-certificates
COPY s3 /s3
EXPOSE 8080
ENTRYPOINT ["/s3"]