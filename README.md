# scrappy

A Go application that monitors product prices across configurable websites, tracks price history, and sends an email alert whenever a price drops.

Designed to run once per day via a cron job.

---

## How it works

1. Scrapes product listings from configurable websites (defined in a JSON file)
2. Saves each product and its current price to SQLite
3. Compares today's price with the last recorded one
4. If any prices dropped → sends an email alert

---

## Project structure

```
scrappy/
├── cmd/
│   ├── scrappy/                # Entry point
│   │   └── main.go
│   └── migrate/                # DB migration CLI
│       └── migrations/         # SQL up/down files
├── internal/
│   ├── database/               # SQLite models
│   │   ├── history.go
│   │   ├── models.go
│   │   ├── openDB.go
│   │   ├── product.go
│   │   └── site.go
│   ├── mail/                   # Email handling
│   │   ├── config.go
│   │   ├── mail.go
│   │   └── template.html
│   ├── scraper/                # Scraper modules
│   │   ├── builder.go
│   │   ├── extractPrice.go
│   │   ├── saveProducts.go
│   │   ├── scraper.go
│   │   └── scrapSites.go
│   └── types/                  # Shared types
│       └── target.go
├── assets/
│   └── config-example.json    # Example config file
└── Makefile
```

---

## Requirements

- Go 1.26+
- GCC (required by `go-sqlite3` via CGO)
- Chrome or Chromium (only if `waitFor` is set in config)

---

## Setup

### 1. Run the migrations

```bash
make upDB
```

### 2. Configure email

On first run, scrappy will look for `scrappy.env` in your OS config directory and create an empty template if it doesn't exist:

```
~/.config/scrappy.env                        # Linux
~/Library/Application Support/scrappy.env   # macOS
```

Fill it in:

```env
FROM=you@example.com
TO=recipient@example.com
PASSWORD=your_smtp_password
HOST=smtp.example.com
PORT=587
```

Gmail is supported. Use `smtp.gmail.com` as the host, port `587`, and an [App Password](https://support.google.com/accounts/answer/185833) instead of your regular account password.

### 3. Configure the websites to scrape

Copy the example config and edit it to define which websites to track:

```bash
cp assets/config-example.json my-scrapees.json
```

The config defines each website to scrape:

- `name` — unique identifier
- `site` — base URL
- `page` — product listing page
- CSS selectors:
  - `containerPath` — product card container
  - `productNamePath` — product name
  - `productDescPath` — product description (optional)
  - `productPricePath` — product price
  - `productURLPath` — product link
- Optional:
  - `urlWithSite` — prepend base URL to links
  - `waitFor` — CSS selector to wait for before scraping (uses headless Chrome)

### 4. Build and run

```bash
make build
./bin/scrappy --db data.db --cf my-scrapees.json
```

---

## Usage

Once set up, run it manually:

```bash
./bin/scrappy --db data.db --cf my-scrapees.json
```

Or schedule it as a daily cron job:

```bash
# Runs every day at 08:00
0 8 * * * /path/to/scrappy/bin/scrappy --db /path/to/data.db --cf /path/to/my-scrapees.json
```

On the first run, scrappy will populate the database with current prices — no email will be sent since there's no previous price to compare against. From the second run onwards, it will start detecting drops.

---

## Email preview

![Email Preview](assets/preview.png)

---

## Makefile targets

| Target         | Description                             |
| -------------- | --------------------------------------- |
| `make build`   | Format → vet → compile to `bin/scrappy` |
| `make upDB`    | Apply migrations (create tables)        |
| `make downDB`  | Roll back migrations (drop tables)      |
| `make resetDB` | Full schema wipe and recreate           |

---

## Exit codes

| Code | Meaning                          |
| ---- | -------------------------------- |
| `0`  | Success — email sent             |
| `1`  | Error                            |
| `2`  | No price drops — nothing to send |

## TODO

- [x] Code clean up
- [x] Show old price on email
- [x] Instead of detecting price drops, detect price changes - this way price increases would be flagged
- [x] Make websites to be scrapped configurable
