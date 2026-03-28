FROM busybox AS bin
COPY ./dist /dist
RUN if [[ "$(arch)" == "x86_64" ]]; then \
        architecture="amd64"; \
    else \
        architecture="arm"; \
    fi; \
    cp /dist/mealie-webhook-handler_linux-${architecture} /bin/mealie-webhook-handler && \
    chmod +x /bin/mealie-webhook-handler && \
    chown 1000:1000 /bin/mealie-webhook-handler

FROM alpine
RUN adduser -D -u 1000 mealie-webhook-handler
USER 1000

LABEL org.opencontainers.image.title="mealie-webhook-handler" \
      org.opencontainers.image.description="Webhook handler for recipes hosted on a mealie.io instance" \
      org.opencontainers.image.ref.name="main" \
      org.opencontainers.image.licenses='GPL v3' \
      org.opencontainers.image.vendor="Timo Reymann <mail@timo-reymann.de>" \
      org.opencontainers.image.authors="Timo Reymann <mail@timo-reymann.de>" \
      org.opencontainers.image.url="https://github.com/timo-reymann/mealie-webhook-handler" \
      org.opencontainers.image.documentation="https://github.com/timo-reymann/mealie-webhook-handler" \
      org.opencontainers.image.source="https://github.com/timo-reymann/mealie-webhook-handler.git"

COPY --from=bin /bin/mealie-webhook-handler /bin/mealie-webhook-handler
