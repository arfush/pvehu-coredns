FROM debian

COPY coredns /usr/bin/coredns
#COPY Corefile /etc/coredns/Corefile

ENTRYPOINT ["/usr/bin/coredns", "-conf", "/Corefile"]
