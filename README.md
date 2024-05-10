## Example Usage
Generate the Root CA key
`openssl genrsa -des3 -out ./cmd/cli/rootCA.key 4096`

Create and self sign the Root CA Certificate
`openssl req -x509 -new -nodes -key ./cmd/cli/rootCA.key -sha256 -days 1024 -out ./cmd/cli/rootCA.crt`

Generate dummy CSR and Key
`cd exampleFiles && ./csr_generation.sh && cd ../`

Build the CLI
`go build -o ./ssl-cert-signer ./cmd/cli`

Check for usage info
`./ssl-cert-signer --help`

Run the CLI to sign your cert
`./ssl-cert-signer exampleFiles/mydomain.com.csr`

## Useful Links
* https://github.com/Andrew4d3/go-csv2json/blob/master/csv2json.go
* https://levelup.gitconnected.com/tutorial-how-to-create-a-cli-tool-in-golang-a0fd980264f
* https://gist.github.com/salrashid123/1fd267cf213c1a1fe9e6c35c78b47e83
* https://gist.github.com/fntlnz/cf14feb5a46b2eda428e000157447309
* https://www.sslshopper.com/certificate-decoder.html
