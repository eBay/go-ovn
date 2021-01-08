FROM openvswitch/ovn:2.12_e60f2f2_debian_master

ARG GOVERPKG=go1.13.9.linux-amd64.tar.gz
RUN apt-get update
RUN apt-get install -y wget git gcc python3-pip python3-dev build-essential autoconf automake libtool

RUN wget -P /tmp https://dl.google.com/go/$GOVERPKG

RUN tar -C /usr/local -xzf /tmp/$GOVERPKG
RUN rm /tmp/$GOVERPKG

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH

RUN apt-get install -y supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

ENTRYPOINT ["/usr/bin/supervisord"]