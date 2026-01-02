---
name: cli
description: Unified command-line interface for managing bookmarks (linkding) and feeds (miniflux). Use for authentication, managing bookmarks, and managing feeds.
---

# cli

A unified command-line interface for managing bookmarks (via Linkding) and feeds (via Miniflux).

## Critical Notes

1. **Pagination**: All `list` commands return `{total, items}`. Default limit is 10, default offset is 0. Use `--limit` and `--offset` for pagination.
2. **Output filtering**: Use `--jq` for inline filtering or `--json "field1,field2"` to select specific fields.
3. **Quote handling**: For values with double quotes, wrap in single quotes: `--notes 'Title: "Example"'`

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

### Check Total Results Before Processing

Before processing results, verify you have all of them:

```bash
cli list entries --status unread --jq '{total: .total, returned: (.items | length)}'
```

If `total > returned`, either increase the limit or paginate with offset:

```bash
# Increase limit to get all results
cli list entries --status unread --limit 100

# Or paginate through results
cli list entries --status unread --limit 10 --offset 0
cli list entries --status unread --limit 10 --offset 10
cli list entries --status unread --limit 10 --offset 20
```

### List Unread Entries

Get unread entries with feed context:

```bash
cli list entries --status unread --jq ".items[] | { id, url, title, published_at, status, feed_id: .feed.id, feed_title: .feed.title }"
```

Output fields:
- `id`: Entry ID (use for marking read/starred)
- `url`: Original article URL
- `feed_id`, `feed_title`: Source feed info for grouping/filtering

### List Entries by Feed

When you have a `feed_id` from a previous query, fetch more entries from that feed:

```bash
cli list entries --feed-id 42 --limit 20 --jq ".items[] | { id, url, title, published_at }"
```

### Find Starred/Read Entries by Date

Use `changed_at` to filter by when entries were starred or marked read:

```bash
cli list entries --starred --status read --limit 100 --json "id,url,title,changed_at,starred" | jq '.items[] | select(.changed_at >= "2025-12-26")'
```

Note: `changed_at` reflects when the entry was last modified (starred, read status changed), not publication date.

### Add a Feed

```bash
cli add feed <url>
```

The URL must point to a valid RSS/Atom feed.

### Add a Bookmark

Basic:
```bash
cli add bookmark <url>
```

With metadata:
```bash
cli add bookmark <url> --notes 'Title: "Some Title"' --tags "tag1 tag2"
```

Tags are space-separated within the quoted string.
