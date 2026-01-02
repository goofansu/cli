---
name: cli
description: Unified command-line interface for managing bookmarks (linkding) and feeds (miniflux). Use for authentication, managing bookmarks, and managing feeds.
---

# cli

A unified command-line interface for managing bookmarks (via Linkding) and feeds (via Miniflux).

## Commands

```bash
cli login <service>       # Authenticate with miniflux or linkding
cli logout <service>      # Remove stored credentials
cli add bookmark <url>    # Add bookmark to Linkding
cli add feed <url>        # Add feed to Miniflux
cli list bookmarks        # List bookmarks
cli list entries          # List feed entries
```

Use `--help` on any command for options.

## Workflows

### Get unread Entries

```bash
cli list entries --status unread --jq ".[] | { id, url, title, published_at, status, feed_id: .feed.id, feed_title: .feed.title }"
```

### Get More Entries by Feed

When you have the feed ID and title from a previous query, use `--feed-id` to get more entries from that specific feed:

```bash
cli list entries --feed-id 42 --limit 20 --jq ".[] | { id, url, title, published_at }"
```

### Add a Feed

Before subscribing to a new feed, verify the URL points to a valid RSS/Atom feed to prevent rate limiting issues:

```bash
cli add feed <url>
```

### Add Bookmark with Notes

When adding bookmarks with notes that contain double quotes, use single quotes around the entire notes value:

```bash
cli add bookmark <url> --notes 'Title: "Some Title"' --tags "tag1 tag2"
```

Alternatively, escape the inner double quotes:

```bash
cli add bookmark <url> --notes "Title: \"Some Title\"" --tags "tag1 tag2"
```
