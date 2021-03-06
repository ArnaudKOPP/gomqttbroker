FROM alpine
COPY GoMqttBroker /
COPY ssl /ssl
COPY conf /conf

EXPOSE 1883
EXPOSE 1888
EXPOSE 8883
EXPOSE 1993

CMD ["/GoMqttBroker"]