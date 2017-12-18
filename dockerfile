FROM ubuntu:latest
MAINTAINER "github.com/lflxp"

ADD dockerapi /bin/dockerapi
ENTRYPOINT ["/bin/dockerapi"]
