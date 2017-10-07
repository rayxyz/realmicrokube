FROM ray-xyz.com:9090/ubuntu1604

MAINTAINER Raywang

ADD bin/api /usr/local/bin

RUN ["chmod", "+x", "/usr/local/bin/api"]

ENTRYPOINT /usr/local/bin/api