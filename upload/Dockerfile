FROM alpine
RUN apk add --no-cache bash ca-certificates
COPY upload /upload
EXPOSE 8080
ENTRYPOINT ["/upload"]