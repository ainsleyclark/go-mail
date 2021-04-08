<p align="center">
  <img alt="GitHub Logo" src="docs/github_logo.png" height="140" />
  <h3 align="center">Go Mail</h3>
  <p align="center">A cross platform Mailer for GoLang.</p>
  <p align="center">
    <a href="https://github.com/ainsleyclark/go-mail/latest"><img alt="Release" src="https://img.shields.io/github/release/ainsleyclark/go-mail.svg?style=flat-square"></a>
    <a href="https://travis-ci.com/ainsleyclark/go-mail"><img alt="Travis" src="https://www.travis-ci.com/ainsleyclark/go-mail.svg?branch=main"></a>
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
    <a href='https://coveralls.io/github/ainsleyclark/go-mail?branch=main'><img src='https://coveralls.io/repos/github/ainsleyclark/go-mail/badge.svg?branch=main' alt='Coverage Status' /></a>
    <a href="https://goreportcard.com/report/github.com/ainsleyclark/go-mail"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/ainsleyclark/go-mail"></a>
    <a href="https://pkg.go.dev/github.com/ainsleyclark/go-mail"><img src="https://godoc.org/github.com/ainsleyclark/go-mail?status.svg" alt="GoDoc"></a>
  </p>
</p>

## Introduction

Go Mail combines popular mail API's into a singular easy to use interface. Great for 

```go
cfg := mail.Config{
    URL:         "https://api.eu.sparkpost.com",
    APIKey:      "my-key",
    FromAddress: "hello@gophers.com",
    FromName:    "Gopher",
}

driver, err := mail.NewClient(mail.SparkPost, cfg)
if err != nil {
    fmt.Println(err)
    return
}

tx := &mail.Transmission{
    Recipients:  []string{"hello@gophers.com"},
    Subject:     "My email",
    HTML:        "<h1>Hello from go mail!</h1>",
}

result, err := driver.Send(tx)
if err != nil {
    fmt.Println(err)
    return
}

fmt.Println(result)
```

## Installation

```bash
go get -u github.com/ainsleyclark/go-mail
```

## Supported API's

| API         | File Read         |   Examples     |
|-------------|-------------------|----------------|
| SparkPost   | **VERSION**       | [Here](test-files/VERSION) |
| MailGun     | **VERSION**       | [Here](test-files/VERSION) |
| SendGrid    | **VERSION**       | [Here](test-files/VERSION) |


## Todo

- Add CC & BCC

## Contributing

We welcome contributors, but please read the [contributing document](CONTRIBUTING.md) before making a pull request.

## Licence

Code Copyright 2021 go mail. Code released under the [MIT Licence](LICENCE).