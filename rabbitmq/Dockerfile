FROM rabbitmq:3.8
COPY rabbitmq.conf /etc/rabbitmq/rabbitmq.conf
RUN rabbitmq-plugins enable --offline rabbitmq_mqtt rabbitmq_web_mqtt rabbitmq_auth_backend_http
