from alpine:latest

EXPOSE 5000
COPY bin/reg /usr/bin/registry
RUN chmod +x /usr/bin/registry
COPY storage.yaml /etc/config.yaml
ENV REGISTRY_CONFIGURATION_PATH /etc/config.yaml

CMD ["/usr/bin/registry"]
