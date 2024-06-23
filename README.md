# ChillFeed

![Fetch Feeds](https://github.com/jbowdre/chillfeed/actions/workflows/fetch_feeds.yml/badge.svg)

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
The list of feeds and some basic configs are stored in `config.yaml`. Each feed requires the feed URL; you can also set a title to override whatever is retrieved from the source.

```yaml
feeds:
  - url: https://runtimeterror.dev/feed.xml
    title: My Blog
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