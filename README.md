# ChillFeed

[![Fetch Feeds](https://github.com/jbowdre/chillfeed/actions/workflows/fetch_feeds.yml/badge.svg)](https://github.com/jbowdre/chillfeed/actions/workflows/fetch_feeds.yml)

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

## Configuration
`config.yaml` contains some basic configuration options as well as the list of feeds to retrieve. Each feed *must* have a `url` field, and may have an optional `title` that can be used in case the title returned by the feed is not descriptive enough.

```yaml
articlesPerPage: 20                           # how many posts to show on each page
fetchWeeks: 4                                 # how many weeks to go back
repo: github.com/jbowdre/chillfeed            # the name of your repo, for the status badge
feeds:
  - url: https://runtimeterror.dev/feed.xml
    title: My Blog                            # overriding this title
  - url: http://whatever.scalzi.com/feed/
  - url: https://pluralistic.net/feed/
  - url: http://xkcd.com/rss.xml
```

The schedule is defined in the workflow at `.github/workflows/fetch_feeds.yml`; adjust accordingly.

```yaml
on:
  push:
    branches:
      - main
    paths-ignore:
      - 'web/articles/*'
  schedule:
    - cron: '0 */4 * * *'  # Run every 4 hours
  workflow_dispatch:  # Allow manual trigger
```

If you want to serve on a custom domain instead of `<username>.github.io`, set the repository secret `CNAME` to the desired domain and be sure that's configured appropriately with your DNS provider.