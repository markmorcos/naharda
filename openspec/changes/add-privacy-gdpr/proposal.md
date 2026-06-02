# add-privacy-gdpr

> A privacy policy page and a data-deletion path — required because the dashboard collects emails.
> Cites `project.md` §3 (don't store PII beyond need), §10 (email capture), §12 (GDPR: privacy
> policy + deletion flow).

## Why

§12 states email collection requires **a privacy policy + a data-deletion flow** (standard SaaS).
v1 ships email capture (`/v1/signups`) but neither — a real compliance gap given the Berlin
jurisdiction. This closes it.

## What changes

- **`/privacy` page** (web, EN + ar-EG): what we collect (email + consent only), why, retention,
  the no-cookies stance, and how to delete. Linked from the footer.
- **Deletion path**: `DELETE /v1/signups` (by email + simple confirmation) that removes the row,
  plus a `privacy@naharda.com` contact as the human fallback.
- Footer "Privacy" link; the email-capture consent line links to `/privacy`.

## Scope

In: the privacy page (bilingual), the deletion endpoint, the footer/consent links, the contact.

## Non-goals

- Cookie consent banner (we use no cookies — none needed).
- Account-level data export (no accounts until v2).
- Double opt-in confirmation emails (no Resend until v2).

## Acceptance criteria

- [ ] `/privacy` (and `/ar/privacy`) describe data collected, basis, retention, and deletion.
- [ ] A user can request deletion of their email and it is removed from `signups`.
- [ ] Footer + signup consent link to the policy; `privacy@naharda.com` published.

## Dependencies

Builds on `public-api` (signups) and `dashboard` (web). No schema change.
