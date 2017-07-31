FROM ubuntu:16.04
MAINTAINER Josh VanderLinden <codekoala@gmail.com>

COPY ./bin/treksum-api /usr/bin/treksum-api

EXPOSE 1323
CMD ["/usr/bin/treksum-api"]
