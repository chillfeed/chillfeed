# ChillFeed

![Fetch Feeds](https://github.com/jbowdre/chillfeed/actions/workflows/fetch_feeds.yml/badge.svg)

ChillFeed is a relaxed feed aggregator that brings your feeds together in one place, with no pressure to keep up. It's designed for those who want a list of cool stuff so they can read what grabs their interest rather than a checklist of items that must be acknowledged.

## Features

- Aggregates multiple RSS, Atom, and JSON feeds
- Updates periodically via GitHub Actions
- Displays feed items in a clean, easy-to-scan format
- No read/unread tracking - browse at your leisure
- Paginated interface for easy navigation
- Dark theme for comfortable viewing
- **Not a reader** - opens articles on the source site, as the author intended

## Configuration
The list of feeds to be retrieve are stored in `feeds.yaml`. Each feed requires the feed URL; you can also set a title to override whatever is retrieved from the source.

```yaml
feeds:
  - url: https://runtimeterror.dev/feed.xml
    title: My Blog
  - url: http://whatever.scalzi.com/feed/
  - url: https://pluralistic.net/feed/
  - url: http://xkcd.com/rss.xml
```
