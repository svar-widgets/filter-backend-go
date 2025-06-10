FROM debian:12-slim
WORKDIR /app
COPY ./query /app
COPY ./dump.sql /app

CMD ["/app/query"]