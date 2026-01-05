---
name: cli
description: Unified command-line interface for managing links (linkding), feeds (miniflux), and pages (wallabag). Use for authentication, managing links, and managing feeds.
---

# cli

A unified command-line interface for managing links (via Linkding), feeds (via Miniflux), and pages (via Wallabag).

## Instructions

When using this CLI tool, follow these guidelines:

### Command Structure

Available commands:
```bash
cli link add <url>    # Add link to Linkding
cli link list         # List links
cli feed add <url>    # Add feed to Miniflux
cli feed list         # List feeds
cli entry list        # List feed entries
cli entry save <id>   # Save entry to third-party service
cli page add <url>    # Add page to Wallabag
cli page list         # List pages
```

Use `--help` on any command for options.

### Critical Guidelines

1. **Pagination**: All `list` commands return `{total, items}` structure:
   - `link list`, `feed list`, `entry list`: Use `--limit` and `--offset` for pagination (default: limit=10, offset=0)
   - `page list`: Use `--page` and `--per-page` for pagination

2. **Output Filtering**:
   - Use `--jq` for inline filtering with jq expressions
   - Use `--json "field1,field2"` to select specific fields
   - All list commands return structured JSON that can be piped to jq

3. **Quote Handling**:
   - For values with double quotes, wrap in single quotes: `--notes 'Title: "Example"'`
   - Tags are space-separated within a quoted string: `--tags "tag1 tag2"`

### Workflow Steps

1. **Before processing results**: Always check if you have all results by comparing `total` vs. returned items count
2. **When paginating**: Use appropriate pagination flags for the command type
3. **For targeted queries**: Use `--jq` to filter and transform output inline
4. **When adding content**: Include relevant metadata (notes, tags) for better organization

## Examples

### Check total results before processing

Before processing results, verify you have all of them:

```bash
cli entry list --status unread --jq '{total: .total, returned: (.items | length)}'
```

If `total > returned`, either increase the limit or paginate with offset:

```bash
# Increase limit to get all results
cli entry list --status unread --limit 100

# Or paginate through results
cli entry list --status unread --limit 10 --offset 0
cli entry list --status unread --limit 10 --offset 10
cli entry list --status unread --limit 10 --offset 20
```

### List unread entries

Get unread entries with feed context:

```bash
cli entry list --status unread --jq ".items[] | { id, url, title, published_at, status, feed_id: .feed.id, feed_title: .feed.title }"
```

Output fields:
- `id`: Entry ID (use for marking read/starred)
- `url`: Original article URL
- `feed_id`, `feed_title`: Source feed info for grouping/filtering

### List entries by feed

First, find the feed ID:

```bash
cli feed list --jq ".items[] | { id, title, site_url }"
```

Then fetch entries from that feed:

```bash
cli entry list --feed-id 42 --limit 20 --jq ".items[] | { id, url, title, published_at }"
```

### Find starred/read entries by date

Use `changed_at` to filter by when entries were starred or marked read:

```bash
cli entry list --starred --status read --limit 100 --json "id,url,title,changed_at,starred" | jq '.items[] | select(.changed_at >= "2025-12-26")'
```

Note: `changed_at` reflects when the entry was last modified (starred, read status changed), not publication date.

### Save an entry to third-party services

First, find the entry you want to save by listing entries:

```bash
cli entry list --status unread --jq ".items[] | { id, url, title }"
```

Then save it using the entry ID:

```bash
cli entry save 42
```

This saves the entry to Miniflux's third-party integration (e.g., Wallabag, Pocket, etc.), which must be configured in Miniflux settings.

### Add a Feed

```bash
cli feed add <url>
```

The URL must point to a valid RSS/Atom feed.

### Add a feed to a category

First, find the category ID by listing feeds with category information:

```bash
cli feed list --jq ".items[] | { id, title, site_url, category_id: .category.id, category_title: .category.title }"
```

Then add the feed with the category:

```bash
cli feed add <url> --category-id <category_id>
```

The `--category-id` parameter defaults to 1 (All category) if not specified.

### Add a link

Basic:
```bash
cli link add <url>
```

With metadata:
```bash
cli link add <url> --notes 'Title: "Some Title"' --tags "tag1 tag2"
```

Tags are space-separated within the quoted string.

### Add a page

Basic:
```bash
cli page add <url>
```

With metadata:
```bash
cli page add <url> --tags "tag1 tag2" --archive
```

### List pages

Get pages with filtering:
```bash
cli page list --starred --per-page 20 --jq ".items[] | { id, url, title, domain_name }"
```

Filter by domain or tags:
```bash
cli page list --domain example.com
cli page list --tags "tech news"
```

Tags are space-separated within the quoted string.
