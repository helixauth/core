# Helix


## Install

If you don't have Go installed, install it with:
```sh
$ {apt,yum,brew} install golang
$ echo 'export GOPATH=~/go' >> ~/.bashrc
$ source ~/.bashrc
$ mkdir $GOPATH
```


## JWS Keys

To generate new RS256 key pair for JWT signing:
```sh
$ ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key
# Don't add passphrase

$ openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub
$ cat jwtRS256.key 
$ cat jwtRS256.key.pub
```

When setting your secrets management (in the section below), it will be helpful to base64 encode and copy these keys:
```sh
$ cat jwtRS256.key | base64 | pbcopy
```


## Secrets

For secrets managmenet, Helix uses the open-source [SOPS](https://github.com/mozilla/sops) project by Mozilla. In short, SOPS provides a secure way to encrypt and decrypt application secrets such as database credentials, JWS keys, and more.

Your exact SOPS configuration will be custom depend on the infrastructure you deploy Helix to. In order to solve the "bootstrapping of trust" problem in deployment environments, SOPS works natively with the Key Management Systems (KMS) of all major cloud hosting providers (AWS, GCP, Azure). 

For more information on SOPS, setting it up for your infrastructure, and safely commiting encrypted secrets to Git, please reference the [Youtube video](https://www.youtube.com/watch?v=V2PRhxphH2w) by Securing DevOps.

For now, simply install SOPS on your local development machine:
```
$ go get -u go.mozilla.org/sops/cmd/sops
```


### Configuring SOPS with GPG
1. If you don't already have GnuPG installed locally, you can install it using the [binary installer](https://gnupg.org/download/index.html) or using your local pacakge manager:
```
$ {apt,yum,brew} install gnupg2
```

2. If you don't already have a GPG key, generate one with the command below. You will be asked for your full name and email address. When the information is okay, press "o" and "enter". You will be asked to set a optional password. For example:
<pre lang="sh">
$ gpg --generate-key
gpg (GnuPG) 2.2.15; Copyright (C) 2019 Free Software Foundation, Inc.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
Note: Use "gpg --full-generate-key" for a full featured key generation dialog.
GnuPG needs to construct a user ID to identify your key.
Real name: J Doe
Email address: jdoe@email.com
You selected this USER-ID:
    "J Doe <jdoe@email.com>"
Change (N)ame, (E)mail, or (O)kay/(Q)uit? O
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.
gpg: key 411F71D23B22E116 marked as ultimately trusted
gpg: revocation certificate stored as '/Users/jdoe/.gnupg/openpgp-revocs.d/0AB19F525F991CC847F744CA411F71D23B22E116.rev'
public and secret key created and signed.
pub   rsa2048 2019-05-17 [SC] [expires: 2021-05-16]
      <b>0AB19F525F991CC847F744CA411F71D23B22E116</b>
uid                      J Doe <jdoe@email.com>
sub   rsa2048 2019-05-17 [E] [expires: 2021-05-16]
</pre>

3. Now you have a PGP fingerprint (`0AB19F525F991CC847F744CA411F71D23B22E116` in the example above), you can register your public key to a public key server such as https://keyserver.ubuntu.com/. Simply export your public key, paste it into the key submission box, and press 'Submit'.
```sh
$ gpg --armor --export jdoe@email.com | pbcopy
```

4. Configure SOPS to use PGP and point at your public key server. `SOPS_GPG_EXEC` points to the GPG binary. `SOPS_PGP_FP` is a list is a list of comma separated fingerprints. Add the fingerprints of all team members that are allowed to access your Helix secrets. `SOPS_GPG_KEYSERVER` points to the key server where you and your team have registered your public keys.
```sh
export SOPS_GPG_EXEC="gpg"
export SOPS_PGP_FP="0AB19F525F991CC847F744CA411F71D23B22E116"
export SOPS_GPG_KEYSERVER="keyserver.ubuntu.com"
```

5. Create a new secrets file with the following command:
```sh
$ sops cfg/secrets.enc.dev.yaml
```

6. Configure the secrets for your development environment:
```yaml
postgres:
    username: helix
    password: password
    host: helixdb
    port: 5432
    db_name: helixdb
    ssl_mode: disable

jws:
    hs256:
        secret: mysupersecret # Feel free to rename this
    rs256:
        {KEY_ID}: 
            public: {BASE64_ENC_RS256_PUBLIC_KEY} # Please see the JWS Keys section above
            private: {BASE64_ENC_RS256_PRIVATE_KEY} # Please see the JWS Keys section above
```

7. Repeat steps 5 and 6 for your production environment secrets.

