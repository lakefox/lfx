# LFX Forum Software

A simple Hacker News-like forum powered by SQLite.

---

## Setup

Create a `.env` file with the following content:

```env
DATABASE="./forum.db"
JWT_SECRET="your_secret"
SITE_TITLE="Your Forum Name"
POSTS_PER_PAGE="20"
THEME="Sunset Glow"
BAN_TIMEOUT="20"
```

### Customize Favicon

To use a custom favicon, replace the `favicon.webp` file in the `static` folder and rebuild the application.

### Deployment

All server files are embedded in the application, allowing you to copy the binary directly to your server.

---

## Running LFX

### On Your System

```bash
git clone https://github.com/lakefox/lfx.git
cd ./lfx
go build -o lfx
nano .env
./lfx
```

### On Linux Server

```bash
# Build the application for Linux
GOOS=linux GOARCH=amd64 go build -o lfx

# Transfer the binary to your server
scp ./lfx root@your.ip.add.ress:/

# Connect to your server
ssh root@your.ip.add.ress

# Set up the environment
nano .env
cd /
./lfx
```

---

## Themes

LFX includes a variety of themes to customize your forum’s appearance:

- Serenity Blue
- Sunset Glow
- Forest Whisper
- Monochrome Minimalist
- Vibrant Energy
- Autumn Harvest
- Ocean Breeze
- Midnight Elegance
- Lavender Bliss
- Amber Sunrise

Choose a theme by setting the `THEME` variable in your `.env` file.

---

## Spam Filtering

By default, posts, comments, and usernames containing words from [Google's Profanity Words List](https://github.com/coffee-and-fun/google-profanity-words/blob/main/data/en.txt) are rejected.

You can configure the ban duration using the `BAN_TIMEOUT` variable in the `.env` file (value in minutes).

---
