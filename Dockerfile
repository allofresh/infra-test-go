FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

COPY deploy/_output/rest/prod /usr/local/bin/prod
RUN microdnf update -y && microdnf clean all

EXPOSE 8080

ENTRYPOINT ["prod"]
