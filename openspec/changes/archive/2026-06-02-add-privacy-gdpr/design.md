# Design — add-privacy-gdpr

## Privacy page
A static (long-cached) Astro page at `/privacy` + `/ar/privacy`, plain-language: collected data
(email + consent boolean + locale + timestamp — nothing else), lawful basis (consent), retention
(until deletion), no cookies / cookieless analytics, third parties (Open-Meteo, Aladhan, etc. for
data — not personal), and the deletion mechanism + contact.

## Deletion
```
  DELETE /v1/signups   body { email }
    → delete the row WHERE email = lower(email); always 200 (don't reveal existence — anti-enumeration)
    → rate-limited (existing IP limiter); honeypot not needed (idempotent delete)
  Fallback: privacy@naharda.com (manual deletion within the statutory window)
```
No confirmation email in v1 (no Resend) — a self-serve delete-by-email + the contact address
satisfy the §12 deletion-flow requirement for a pre-accounts product. Hardening (signed
confirmation link) lands with the v2 email infra.

## Cross-cutting
- Footer gains a "Privacy" / "الخصوصية" link (Base.astro, localized).
- The signup consent checkbox text links to `/privacy`.
- Page is `Cache-Control` long-lived (rarely changes).
