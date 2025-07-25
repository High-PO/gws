# gws

> What is gws?
- go + aws cli v2 = gws

You can easily use AWS CLI v2 as an MFA-enabled IAM user on a Mac.

```
cd src && go build -o gws main.go
sudo cp ./gws /usr/local/bin/
```

## Command

### case 1
Use Default User
```
gws <mfa_token>
```

Use Custom User
```
gws <profile> <mfa_token>
```

Help Command
```
gws help
```
```
gws <USER> help
```

Show Version
```
gws --version
```


