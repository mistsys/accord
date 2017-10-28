## Accord: End to End secure SSH Certificate management in Public Cloud

This was originally inspired by the desire to use [Netflix's BLESS SSH CA](https://github.com/Netflix/bless) service and Facebook's [Scalable and Secure Access with SSH](https://code.facebook.com/posts/365787980419535/scalable-and-secure-access-with-ssh/) implementations. While these are great efforts, we wanted something that handled both host ssh certificates, client certificates and the certificate management lifecycle. Additionally, Lyft's [blessclient](https://eng.lyft.com/blessing-your-ssh-at-lyft-a1b38f81629d) is great but we have our authentication standardized on Google Apps, and wanted to utilize that.

Accord is a distillation of all that I have learned into looking at implementating a secure SSH-CA infrastructure.

Accord:
- Runs on a separate AWS account, and can be ported relatively easily to any other cloud. The terraform task is included to spin up the instances in your own AWS account and start using it
- You can deploy self-signed certificate or use LetsEncrypt certificate for the server
- Pure Go: This makes deployment and validation of code straightforward in the sense that you are not linking with version of openssl that the system provides. This is arguable, and I still believe in using `ssh-keygen` on a fresh machine to create the ssh-ca certificates.
- There are two SSH certificates allowed on the server, with potentially different lifecycles
   - User CA
   - Host CA
- Accord server uses AWS Parameter Store to read the passphrase to decrypt the ssh ca keys, this can only be done by assuming the role that allows reading the keys
   - This means any kind of tampering for the keys are logged in CloudTrail's audit log


## How to develop

Get the latest grpc

`go get google.golang.org/grpc`

`go get -u github.com/golang/protobuf/protoc-gen-go`

Then run `go generate`

If you need to implement or have added a new method in the interface, you can generate it with impl against the protocol

`go get -u github.com/josharian/impl`


`impl 's *CertServer' github.com/mistsys/accord/protocol.CertServer` will print out an implementation based on the interface

## Generating CA cert

There are two ways to generate CA cert.


### OpenSSH
The standard way is to use OpenSSH's `ssh-keygen`, you do want to make sure the following pre-requisites are met.

OpenSSH is at least newer than 7.3 and linked with LibreSSL or BoringSSL. You can check this with the following command:

```
ssh -V
```

In my OSX 10.12.15 machine I get

```
>OpenSSH_7.4p1, LibreSSL 2.5.0
```

If it is linked against OpenSSL, make sure it is after version 1.0.1. Make sure the system you are running this on has a reliable source of randomness with `/dev/urandom` (just make sure that this exists, it's very hard to determine with real accuracy your random number generator is truly random. For the adventurous, the [NIST-800-22](https://csrc.nist.gov/publications/detail/sp/800-22/rev-1a/final) document has details)

`keygen -b 4096 -f ssh_ca`

The default is RSA which, as far as we know is relatively safe for 4096 bits as long as its rotated frequently enough. Protecting against NSA in a public cloud isn't a reasonable threat model. There is a useful [website for checking keylengths](https://www.keylength.com/en/compare/) against various standard bodies recommendations.

Pick a reasonably random passphrase for it, save it in Parameter Store (TBD) so that it can be retrieved programmatically in code and used to decrypt the key.

### Accord

I initially wrote this to understand what was going under the hood of a SSH CA cert and it can possibly be good enough for dynamically rotating a certificate.


This syntax will change but for the time being this is how you can get the same certs

```
go run accord.go -task=genrootcert test_ca
```

This will generate a root cert without a passphrase, so it's not secure yet. In a production deployment, you want to use a passphrase for protecting the key. You can do that using the `-password` argument. The generated certs will include the validity period that tools can use to validate the host cert from known_hosts, and if the cert is close to expiry, get a new one.

```
go run accord.go -task=genrootcert -password hello test_encrypted_ca

```



### Rotation Procedure

- Delete the old passphrase key and then the certificate files, they shouldn't be accessible past the time

### Rotation Cycle

- Certificates should have a life time of no more than a year, ideally they should be only 3 months. This is for practical reasons

## Creating and signing certs

### Creating cert

Two certs are generated - for signing for hosts and users separately

Root Cert with a passphrase

```
go run cmd/accord/accord.go -task=genrootcert -password "staple horse apple newton" root_ca_20170927
```

This will create two files `root_ca_20170927` and `root_ca_20170927.pub`

User Cert with a passphrase

```
go run cmd/accord/accord.go -task=genrootcert -password "staple horse apple thatcher" user_ca_20170927
```

This will create two files `user_ca_20170927` and `user_ca_20170927.pub`

### Signing the User Key

```
go run cmd/accord/accord.go -certkey=user_ca_20170927 -password="staple horse apple thatcher" -task=genusercert -pubkey=$HOME/.ssh/id_rsa.pub
```

This will write the user key to stdout (this is generally intended to be signed and emailed to the user or replied in API)

### Signing the Host Key

```
go run cmd/accord/accord.go -certkey=root_ca_20170927 -password="staple horse apple newton" -hostname ec2-<ip-addr>.amazonaws.com -task=genhostcert -pubkey=ssh_host_rsa_key.pub
```

This will write the host key to stdout too.

## Testing end to end

The client and server talk over HTTP/2 gRPC protocol, the generated certificates can be used in `-insecure` mode (that works without TLS and intended only for development) to test end to end.

### Running the server

Run this from `cmd/accord_server`

```
go run server.go -rootca ../root_ca_20170927 -rootcapassword="staple horse apple newton" -userca ../user_ca_20170927 -usercapassword "staple horse apple thatcher" -insecure
```

### Running the client to sign host SSH keys

Run this from `cmd/accord_client`
```
go run client.go -task=hostcert -insecure -deploymentId=test -psk=JpUtbRukLuIFyjeKpA4fIpjgs6MTV8eH -hostkeys=test_assets/test_pubkeys/ -host=host.example.com
```

You can check the generated cert with `ssh-keygen`

```
-> % ssh-keygen -f test_assets/test_pubkeys/ssh_host_rsa_key-cert.pub -L
test_assets/test_pubkeys/ssh_host_rsa_key-cert.pub:
        Type: ssh-rsa-cert-v01@openssh.com host certificate
        Public key: RSA-CERT SHA256:uOYQUE2YSU0AJVfIgEafHcrldX++liMRc5hDDcitD2Y
        Signing CA: RSA SHA256:HCJ9E83f7KdVF+yolsAJx1B+a8WWZvOUoX8ZQtBZQrU
        Key ID: "f93b6b67-19f6-f5fc-2793-01e101dfb073"
        Serial: 1
        Valid: from 2017-10-08T00:26:51 to 2017-11-06T23:26:51
        Principals:
                host.example.com
        Critical Options: (none)
        Extensions: (none)
```

### Users requesting their SSH certificates

This will print the cert files after getting them signed by the server.


## Similar Projects

- https://github.com/cloudtools/ssh-cert-authority

# References

- https://speakerdeck.com/rlewis/how-netflix-gives-all-its-engineers-ssh-access-to-instances-running-in-production
- https://github.com/Netflix/bless
- https://github.com/lyft/python-blessclient
- https://code.facebook.com/posts/365787980419535/scalable-and-secure-access-with-ssh/
- https://ef.gy/hardening-ssh
- https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.certkeys
