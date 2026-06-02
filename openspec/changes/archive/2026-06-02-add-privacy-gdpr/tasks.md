# Tasks — add-privacy-gdpr

## Slice 1 — Privacy page
- [x] `/privacy` + `/ar/privacy` (bilingual): data collected, basis, retention, no-cookies, deletion.
- [x] Footer "Privacy"/"الخصوصية" link (Base.astro, localized); long cache.
- [x] Signup consent text links to `/privacy`.

## Slice 2 — Deletion
- [x] `DELETE /v1/signups` (by email, anti-enumeration 200, rate-limited) → removes the row.
- [x] Publish `privacy@naharda.com` as the human fallback on the page.
- [x] Verify: a captured email can be deleted end-to-end.
