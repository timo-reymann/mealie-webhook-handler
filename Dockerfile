FROM busybox AS bin
COPY ./dist /dist
RUN if [[ "$(arch)" == "x86_64" ]]; then \
        architecture="amd64"; \
    else \
        architecture="arm"; \
    fi; \
    cp /dist/mealie-webhook-handler_${architecture} /bin/mealie-webhook-handler && \
    chmod +x /bin/mealie-webhook-handler && \
    chown 1000:1000 /bin/mealie-webhook-handler \

FROM chainguard/wolfi-base
LABEL org.opencontainers.image.title="mealie-webhook-handler"
LABEL org.opencontainers.image.description="Webhook handler for recipes hosted on a mealie.io instance"
LABEL org.opencontainers.image.ref.name="main"
LABEL org.opencontainers.image.licenses='GPL v3'
LABEL org.opencontainers.image.vendor="Timo Reymann <mail@timo-reymann.de>"
LABEL org.opencontainers.image.authors="Timo Reymann <mail@timo-reymann.de>"
LABEL org.opencontainers.image.url="https://github.com/timo-reymann/mealie-webhook-handler"
LABEL org.opencontainers.image.documentation="https://github.com/timo-reymann/mealie-webhook-handler"
LABEL org.opencontainers.image.source="https://github.com/timo-reymann/mealie-webhook-handler.git"
RUN adduser -D -u 1000 mealie-webhook-handler
USER 1000
COPY --from=bin /bin/mealie-webhook-handler /bin/mealie-webhook-handler