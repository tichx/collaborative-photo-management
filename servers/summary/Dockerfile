FROM alpine
COPY summary /summary
RUN apk add --no-cache bash ca-certificates
EXPOSE 80
ENTRYPOINT ["/summary"]