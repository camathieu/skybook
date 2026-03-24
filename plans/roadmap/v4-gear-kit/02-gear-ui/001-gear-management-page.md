---
ticket: "001"
epic: gear-ui
milestone: v4
title: Gear Management Page
status: planned
priority: high
estimate: L
---

# Gear Management Page

## Acceptance Criteria

- [ ] New "Equipment" page accessible from the main navigation
- [ ] Gear list table showing: Type, Manufacturer, Model, Size, Active status
- [ ] Filter by gear type and active/inactive
- [ ] Add gear form with all fields: Type, Manufacturer, Model, Size, Serial, DOM, PurchaseDate, PurchasePrice, NextMaintenanceAt, Notes, Active
- [ ] Edit gear item inline or via modal
- [ ] Delete gear item with confirmation dialog
- [ ] Visual indicator for items approaching `NextMaintenanceAt` date (amber \< 30 days, red = overdue)
- [ ] Pinia store in `webapp/src/stores/gear.js`
- [ ] API client methods in `webapp/src/api.js`
- [ ] Dark-first styling consistent with existing design system
