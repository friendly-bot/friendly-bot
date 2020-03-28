FROM gcr.io/distroless/base

COPY ./bin/friendly-bot /go/bin/app /

CMD ["/app"]

