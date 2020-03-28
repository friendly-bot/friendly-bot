FROM gcr.io/distroless/base

COPY ./bin/friendly-bot /

CMD ["/app"]

