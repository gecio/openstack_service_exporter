FROM golang:latest as builder

ENV GO111MODULE on
WORKDIR /go/src/github.com/innovocloud/openstack_service_exporter
COPY . /go/src/github.com/innovocloud/openstack_service_exporter
RUN make binary

FROM alpine:latest
LABEL maintainer="iNNOVO Cloud <openstack@innovo-cloud.de>"

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/innovocloud/openstack_service_exporter/openstack_service_exporter /bin/openstack_service_exporter
USER nobody

EXPOSE 9177
CMD ["/bin/openstack_service_exporter"]
