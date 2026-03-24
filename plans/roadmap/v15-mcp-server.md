---
milestone: v15
title: MCP Server
status: planned
---

# v15 — MCP Server

Model Context Protocol (MCP) integration to expose SkyBook APIs and data to AI Assistants (like Claude, Cursor, or Gemini) for automated logbook querying, analysis, and management.

## Epics

- `01-mcp-server-foundation` — Base MCP server implementation (stdio/SSE transport) tied to the `skybook` binary.
- `02-mcp-resources` — Expose logbook jumps, summary stats, and user settings as readable MCP resources.
- `03-mcp-tools` — Expose logbook actions (create jump, search jumps, generate stats) as callable MCP tools.
- `04-mcp-prompts` — Pre-defined MCP prompts for common AI assistant tasks (e.g. "Analyze my recent performance", "Summarize my canopy progression").
