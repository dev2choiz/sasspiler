FROM golang:1.13.5 as sasspiler-builder

MAINTAINER Dev2Choiz

RUN mkdir /src
WORKDIR /src

RUN useradd -ms /bin/bash sasspiler
USER sasspiler
