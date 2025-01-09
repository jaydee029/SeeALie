FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

ADD user /usr/bin/user

EXPOSE 8001

CMD ["user"]