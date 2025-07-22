# gws

> What is gws?
- go + aws cli v2 = gws

You can easily use AWS CLI v2 as an MFA-enabled IAM user on a Mac.

```
go build -o gws main.go
sudo cp ./gws /usr/local/bin/
```

```
gws <mfa_token>
```

```
gws <profile> <mfa_token>
```
