FROM alpine
MAINTAINER Sergey Nuzhdin <ipaq.lw@gmail.com>

COPY bin/kube-replay .

ENTRYPOINT ["./kube-replay"]
