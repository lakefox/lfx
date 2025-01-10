# LFX Forum Software

Simple Hacker News like forum using Sqlite.

## Setup

Make a .env file with the following contents:

```sh
DATABASE="./forum.db"
JWT_SECRET="your_secret"
SITE_TITLE="Your Forum Name"
```

Set your favicon by adding a `favicon.webp` file to the static folder.

Change the look of the forum by modifing the CSS variables in `/static/styles.css`.

## Running LFX

```go
cd ./lfx
go build main.go
./main
```
