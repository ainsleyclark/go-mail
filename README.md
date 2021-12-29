<p align="center">
  <img alt="Gopher" src="logo.svg" height="250" />
  <h3 align="center">Go Mail</h3>
  <p align="center">A cross platform mail driver for GoLang.</p>
  <p align="center">
		<a href="https://github.com/ainsleyclark/go-mail/actions/workflows/test.yml"><img src="https://github.com/ainsleyclark/go-mail/actions/workflows/test.yml/badge.svg?branch=main"></a>
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
		<a href="https://codecov.io/gh/ainsleyclark/go-mail"><img src="https://codecov.io/gh/ainsleyclark/go-mail/branch/main/graph/badge.svg?token=1ZI9R34CHQ"/></a>
    <a href="https://goreportcard.com/report/github.com/ainsleyclark/go-mail"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/ainsleyclark/go-mail?update=true"></a>
    <a href="https://pkg.go.dev/github.com/ainsleyclark/go-mail"><img src="https://godoc.org/github.com/ainsleyclark/go-mail?status.svg" alt="GoDoc"></a>
  </p>
</p>

## Introduction

Go Mail aims to unify multiple popular mail API's (SparkPost, MailGun, SendGrid, Postal & SMTP) into a singular easy to use interface. Email sending is seriously simple and great for allowing the developer to
choose what platform they use.


```go
cfg := mail.Config{
    URL:         "https://api.eu.sparkpost.com",
    APIKey:      "my-key",
    FromAddress: "hello@gophers.com",
    FromName:    "Gopher",
}

mailer, err := drivers.NewSparkPost(cfg)
if err != nil {
	log.Fatalln(err)
}

tx := &mail.Transmission{
    Recipients:  []string{"hello@gophers.com"},
    Subject:     "My email",
    HTML:        "<h1>Hello from Go Mail!</h1>",
}

result, err := mailer.Send(tx)
if err != nil {
	log.Fatalln(err)
}

fmt.Printf("%+v\n", result)
```

## Installation

```bash
go get -u github.com/ainsleyclark/go-mail
```

## Supported API's

Currently, Sparkpost, MailGun and SendGrid is supported, if you want to see more, just submit a feature request or create a new Driver and
submit a pull request.

| API         | Dependency                                                                   | Examples                      |
|-------------|------------------------------------------------------------------------------|-------------------------------|
| SparkPost   | [github.com/SparkPost/gosparkpost](https://github.com/SparkPost/gosparkpost) | [Here](examples/sparkpost.go) |
| MailGun     | [github.com/mailgun/mailgun-go/v4](github.com/mailgun/mailgun-go/v4])        | [Here](examples/mailgun.go)   |
| SendGrid    | [github.com/sendgrid/sendgrid-go](github.com/sendgrid/sendgrid-go)           | [Here](examples/sendgrid.go)  |
| Postmark    |  None         																														   | [Here](examples/postmark.go)  |
| Postal      |  None         																														   | [Here](examples/postal.go)  |
| SMTP        |  None - only use in development.                                             | [Here](examples/smtp.go)      |

## Docs

Documentation can be found at the [Go Docs](https://pkg.go.dev/github.com/ainsleyclark/go-mail), but we have included a kick start guide below to get you started.

### Creating a new client:

You can create a new driver by calling the `drivers` package and passing in a configuration type which is  needed to create a new mailer, each platform requires its own data,
for example, MailGun requires a domain, but SparkPost doesn't.
This is based of the requirements for the API. For more details see the examples above.

```go
cfg := mail.Config{
    URL:         "https://api.eu.sparkpost.com",
    APIKey:      "my-key",
    FromAddress: "hello@gophers.com",
    FromName:    "Gopher",
}

mailer, err := drivers.NewSparkpost(cfg)
if err != nil {
	log.Fatalln(err)
}
```

### Sending Data:

A transmission is required to transmit to a mailer as shown below. Once send is called, a `mail.Result` and error will be returned
indicating if the transmission was successful.

```go
tx := &mail.Transmission{
    Recipients: []string{"hello@gophers.com"},
    Subject:    "My email",
    HTML:       "<h1>Hello from Go Mail!</h1>",
    PlainText:  "plain text",
}

result, err := mailer.Send(tx)
if err != nil {
	log.Fatalln(err)
}

fmt.Printf("%+v\n", result)
```

### Adding attachments:

Adding attachments to the transmission is as simple as passing a byte slice and filename,
Go Mail takes care of the rest for you.

```go
image, err := ioutil.ReadFile("gopher.jpg")
if err != nil {
	log.Fatalln(err)
}

tx := &mail.Transmission{
    Recipients: []string{"hello@gophers.com"},
    Subject:    "My email",
    HTML:       "<h1>Hello from Go Mail!</h1>",
    PlainText:  "plain text",
    Attachments: mail.Attachments{
        mail.Attachment{
            Filename: "gopher.jpg",
            Bytes:    image,
        },
    },
}
```

## Todo

- Remove external dependencies.

## Contributing

We welcome contributors, but please read the [contributing document](CONTRIBUTING.md) before making a pull request.

## Licence

Code Copyright 2021 Go Mail. Code released under the [MIT Licence](LICENCE).
