# Tasks — add-privacy-gdpr

## Slice 1 — Privacy page
- [ ] `/privacy` + `/ar/privacy` (bilingual): data collected, basis, retention, no-cookies, deletion.
- [ ] Footer "Privacy"/"الخصوصية" link (Base.astro, localized); long cache.
- [ ] Signup consent text links to `/privacy`.

## Slice 2 — Deletion
- [ ] `DELETE /v1/signups` (by email, anti-enumeration 200, rate-limited) → removes the row.
- [ ] Publish `privacy@naharda.com` as the human fallback on the page.
- [ ] Verify: a captured email can be deleted end-to-end.
