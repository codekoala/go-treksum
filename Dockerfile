FROM ubuntu:16.04
MAINTAINER Josh VanderLinden <codekoala@gmail.com>

COPY ./bin/treksum-* /usr/bin/

EXPOSE 1323
CMD ["treksum-api"]
