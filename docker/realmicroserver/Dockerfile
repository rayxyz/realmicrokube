FROM ray-xyz.com:9090/ubuntu1604

MAINTAINER Raywang

ADD bin/server /usr/local/bin

RUN ["chmod", "+x", "/usr/local/bin/server"]

ENTRYPOINT /usr/local/bin/server