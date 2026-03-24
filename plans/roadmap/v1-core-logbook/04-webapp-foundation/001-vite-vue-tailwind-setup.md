---
ticket: "001"
epic: webapp-foundation
milestone: v1
title: Vite + Vue + Tailwind Setup
status: planned
priority: high
estimate: S
---

# Vite + Vue + Tailwind Setup

Initialize the webapp with Vite 7, Vue 3, and Tailwind CSS 4 (matching Plik's frontend stack).

## Acceptance Criteria

- [ ] `webapp/package.json` with vue, vue-router, pinia, @tailwindcss/vite, @vitejs/plugin-vue
- [ ] `webapp/vite.config.js` with Vue plugin, Tailwind plugin, and API proxy to `:8080`
- [ ] `webapp/index.html` entry point
- [ ] `webapp/src/main.js` bootstrapping Vue app with router and Pinia
- [ ] `npm run dev` starts dev server on `:5173`
- [ ] `npm run build` produces `webapp/dist/`
