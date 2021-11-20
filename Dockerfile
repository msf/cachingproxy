FROM golang:1.17-bullseye AS builder
ARG APPLICATION
ARG LOC
RUN apt update && apt install -y curl make
COPY . /cachingproxy/
WORKDIR /cachingproxy/${LOC}/${APPLICATION}

RUN make setup build

FROM ubuntu:20.04
RUN apt update \
	&& apt install -y ca-certificates \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /cachingproxy/tmp/build/${APPLICATION} /app/
ENTRYPOINT []
CMD ["/bin/bash"]
