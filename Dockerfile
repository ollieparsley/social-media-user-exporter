FROM alpine

LABEL maintainer="Ollie Parsley <ollie@ollieparsley.com>"

# Sort out CA certs
RUN apk --update add ca-certificates

# Timezones
RUN apk --update add tzdata

# Needed for go binary to run properly
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /

COPY target/build/root/usr/bin/social-media-user-exporter /usr/bin/social-media-user-exporter

ENTRYPOINT ["/usr/bin/social-media-user-exporter"]
