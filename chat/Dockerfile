FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && apt-get install -y ca-certificates

ADD chat /usr/bin/chat

EXPOSE 8002

CMD ["chat"]