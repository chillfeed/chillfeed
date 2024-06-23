# ChillFeed

[![Fetch Feeds](https://github.com/chillfeed/chillfeed/actions/workflows/fetch_feeds.yml/badge.svg)](https://github.com/chillfeed/chillfeed/actions/workflows/fetch_feeds.yml)

ChillFeed is a relaxed feed aggregator that brings your feeds together in one place, with no pressure to keep up. It's designed for those who want a list of cool stuff so they can read what grabs their interest rather than a checklist of items that must be acknowledged.

## Features

- Aggregates multiple RSS, Atom, and JSON feeds
- Updates on a schedule via GitHub Actions
- Served via GitHub Pages
- Displays feed items in a clean, easy-to-scan format
- No read/unread tracking - browse at your leisure
- Paginated interface for easy navigation
- Dark theme for comfortable viewing
- **Not a reader** - opens articles on the source site, as the author intended

## Setup

1. [Use this template](https://github.com/new?template_name=chillfeed&template_owner=chillfeed) to create your own ChillFeed repo.
2. Copy `config.yaml.example` to `config.yaml` and edit it to define your feeds:
```yaml
articlesPerPage: 20                           # how many posts to show on each page
fetchWeeks: 4                                 # how many weeks to go back
feeds:
  - url: https://runtimeterror.dev/feed.xml
    title: jbowdre's Technical Blog           # override this title
  - url: http://blog.jbowdre.lol/feed/
    title: jbowdre's Weblog
```
3. Edit `.github/workflows/fetch_feeds.yaml` to set your preferred schedule.
```yaml
on:
  push:
    branches:
      - main
  schedule:
    - cron: '0 */4 * * *'  # Run every 4 hours
  workflow_dispatch:  # Allow manual trigger
```
4. *(Optional)* If you want to serve ChillFeed on a custom domain instead of `<username>.github.io`, set the repository secret `CNAME` to the desired domain and be sure that's [configured appropriately](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site) with your DNS provider.
5. Commit and push 🤞.

