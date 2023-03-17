FROM rockylinux:8
ARG TAG=latest
RUN yum install -y go git file wget sudo rpm-build
RUN git clone https://github.com/nuvolaris/nuv -b $TAG --depth 1 
RUN cd nuv ;\ 
    go build -o bin/nuv ;\
    cd bin ;\
    ./nuv linux-rpm-all
