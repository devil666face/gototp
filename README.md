# Install

Download and set in your `PATH`

- [windows](https://github.com/devil666face/gototp/releases/latest/download/gototp_windows_amd64.tar.gz)
- [linux](https://github.com/devil666face/gototp/releases/latest/download/gototp_linux_amd64.tar.gz)
- [termux](https://github.com/devil666face/gototp/releases/latest/download/gototp_termux.tar.gz)
- [darwin](https://github.com/devil666face/gototp/releases/latest/download/gototp_darwin_amd64.tar.gz)

# Usage

Set you secret password for database

```
❯ gototp
┃ 🔐 enter your secret
┃ >

enter submit
```

# Add

```
❯ gototp
┃ name
┃ >

  update period
  > 30

  set digits
  > 6
    8

  algorithm
  > SHA1
    SHA256
    SHA512
    MD5

  secret
  >

ctrl+e complete • enter next
```

- set name
- set period (default 30)
- set digits (default 6)
- set algorithm (default sha1)
- set secret key - from service what you want to add

# Show/Code

- show - all secrets
- code - update code in realtime
