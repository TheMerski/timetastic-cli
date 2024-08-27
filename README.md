# timetastic cli

> [!WARNING]
> This tool uses cookies from your local chrome browser to hijack your timetastic session, please make sure to validate how your cookies are being used yourself before running it!

> [!NOTE]
> this tool uses unsupported reverse engineered api's.<br>
> So be warned! It might break at any point in time.

## Using the CLI

Timetastic has no way to create recurring leaves in your calendar, this cli tool fixes that and provides you with an cli to interact with timetastic.

Pre-requisites:
- [Go](https://go.dev/doc/install)
- Timetastic logged in on chrome

Installing:

```bash
go install github.com/themerski/timetastic-cli/timetastic@latest
```

Using the CLI:

```bash
timetastic --help
```

## How it works

Timetastic [provides an API](https://timetastic.co.uk/api/#introduction) but this is only available for Admins, luckily with some reverse engineering we can call the API's that their website calls (and some actual api's).

In order to do this we need 2 things:

### 1. Authentication

Authentication on timetastic is handled by cookies, so we use the great [kooky](https://github.com/browserutils/kooky) library to read the timetastic cookies from Chrome.

> [!NOTE]
> As mentioned by the kooky library, reaching into a browsers cookies is a bad idea.
> And as a users you will never want to give an tool access to you're cookie storage. I'm literally using it to hijack your timetastic session, so please make sure you are comfortable with what this tool does before using it.

### 2. Cross-site request forgery (xsrf) token

This is injected into the timetastic webpage, so with a single request and some [xpath filtering](./internal/authentication/getPagedata.go#L36) we can get it from the page.
