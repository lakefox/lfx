# LFX Forum Software

Simple Hacker News like forum using Sqlite.

## Setup

Make a .env file with the following contents:

```sh
DATABASE="./forum.db"
JWT_SECRET="your_secret"
SITE_TITLE="Your Forum Name"
POSTS_PER_PAGE="20"
THEME="Sunset Glow"
```

To set a custom favicon replace the `favicon.webp` file in the static folder and rebuild the application.

All of the server files are embeded into the application so you can directly copy the binary to the server you plan on running on.

## Running LFX

```go
cd ./lfx
go build main.go
./main
```

## Themes

Serenity Blue
Sunset Glow
Forest Whisper
Monochrome Minimalist
Vibrant Energy
Autumn Harvest
Ocean Breeze
Midnight Elegance
Lavender Bliss
Amber Sunrise
