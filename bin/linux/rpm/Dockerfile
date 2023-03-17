FROM rockylinux:8
RUN yum install -y go git file wget
RUN git clone https://github.com/nuvolaris/nuv
RUN cd nuv ; go build -o bin/nuv
RUN cd nuv/bin ; ./nuv all ARCH=arm64 ; ./nuv all AMD=64
