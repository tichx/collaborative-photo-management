FROM alpine
COPY gateway /gateway
RUN apk add --no-cache bash ca-certificates
EXPOSE 443
ENTRYPOINT ["/gateway"]