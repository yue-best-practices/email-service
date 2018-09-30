FROM scratch
COPY ./email-service /
COPY ./cacert.pem /etc/ssl/certs/
ENTRYPOINT ["./email-service"]