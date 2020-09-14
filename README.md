# Helix


## Install

If you don't have Go installed, install it with:
```s
$ {apt,yum,brew} install golang
$ echo 'export GOPATH=~/go' >> ~/.bashrc
$ source ~/.bashrc
$ mkdir $GOPATH
```


## Signing Keys

To generate new RS256 key pair for JWT signing:
```s
$ ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key
# Don't add passphrase

$ openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub
$ cat jwtRS256.key 
$ cat jwtRS256.key.pub
```

To copy a key (useful for the secrets configuration section below):
```s
$ cat jwtRS256.key | pbcopy
```


## Secrets

For secrets managmenet, Helix uses the open-source [SOPS](https://github.com/mozilla/sops) project by Mozilla. In short, SOPS provides a secure way to encrypt and decrypt application secrets such as database credentials, JWS keys, and more.

Your exact SOPS configuration will be custom depend on the infrastructure you deploy Helix to. In order to solve the "bootstrapping of trust" problem in deployment environments, SOPS works natively with the Key Management Systems (KMS) of all major cloud hosting providers (AWS, GCP, Azure). 

For more information on SOPS, please reference the [Youtube video](https://www.youtube.com/watch?v=V2PRhxphH2w) by Securing DevOps.


### Configuring SOPS with AWS
To setup SOPS with key managed by AWS KMS, please follow the steps below:

1. Install SOPS onto your local development machine:
```s
$ go get -u go.mozilla.org/sops/v3/cmd/sops
```

2. Go to your AWS KMS console and create a new symmetric key. Grant your user account permissions to use the key. Additionally grant permissions to the role assigned to your Helix deployment.

3. Copy the key's ARN and use it generate a new secrets file for your dev environment:
```s
$ sops --kms '{YOUR_KEY_ARN}' cfg/secrets.dev.yaml
```

4. Enter the following secrets:
- Database credentials
- HS256 signing secret
- RS256 pub/priv key pair
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
        secret: mysupersecret

    rs256:
        {YOUR_KEY_ID}: 
            public: "{YOUR_RS256_PUB_KEY}"
            private: "{YOUR_RS256_PRIVATE_KEY}"
```
