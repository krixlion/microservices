FROM fluent/fluentd
USER root
COPY conf /fluentd/etc
RUN ["gem", "install", "fluent-plugin-elasticsearch", "--no-document"]
USER fluent