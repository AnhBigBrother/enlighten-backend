rm -rf *.pem
rm -rf *.srl

# 1. Generate CA's private key and self-signed certificte // https://man.openbsd.org/openssl#req
openssl req -nodes -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=vn/ST=Hanoi/L=Hanoi/O=wololo/OU=foo/CN=AnhBroo/emailAddress=anh.bigbrother@gmail.com"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text


# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -nodes -newkey rsa:4096 -keyout server-key.pem -out server-req.pem -subj "/C=vn/ST=Hanoi/L=Hanoi/O=Enlighten/OU=foo/CN=Enlighten/emailAddress=anh.bigbro@gmail.com"


# 3. Use CA's private key to sign web server CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -days 60 -extfile server-ext.cnf
openssl x509 -in server-cert.pem -noout -text

