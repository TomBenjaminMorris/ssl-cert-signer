openssl genrsa -out mydomain.com.key 2048

openssl req -new -sha256 -key mydomain.com.key -subj "/C=US/ST=CA/O=MyOrg, Inc./CN=mydomain.com" -out mydomain.com.csr

openssl req -in mydomain.com.csr -noout -text