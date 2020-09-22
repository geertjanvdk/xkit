openssl req -config openssl.conf \
  -newkey rsa:2048 \
  -nodes -x509 \
  -keyout test_server.key.pem \
  -out test_server.crt.pem
