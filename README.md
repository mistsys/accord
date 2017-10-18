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
