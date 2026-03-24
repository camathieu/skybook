---
ticket: "001"
epic: webapp-foundation
milestone: v1
title: Vite + Vue + Tailwind Setup
status: done
priority: high
estimate: S
---

# Vite + Vue + Tailwind Setup

Initialize the webapp with Vite 7, Vue 3, and Tailwind CSS 4 (matching Plik's frontend stack).

## Acceptance Criteria

- [x] `webapp/package.json` with vue, vue-router, pinia, @tailwindcss/vite, @vitejs/plugin-vue
- [x] `webapp/vite.config.js` with Vue plugin, Tailwind plugin, and API proxy to `:8080`
- [x] `webapp/index.html` entry point
- [x] `webapp/src/main.js` bootstrapping Vue app with router and Pinia
- [x] `npm run dev` starts dev server on `:5173`
- [x] `npm run build` produces `webapp/dist/`

## Done

- Initialized webapp using Vite 8 + Vue 3 + Tailwind CSS 4.
- Configured plugins in `vite.config.js` and set up API proxy to the Go backend.
- Set up root Vue app, Pinia store, and basic router pointing to an empty logbook view (`JumpList.vue`).
