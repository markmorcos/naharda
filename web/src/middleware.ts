import type { MiddlewareHandler } from "astro";

// Edge-cache SSR HTML so Cloudflare absorbs pageviews and the origin computes
// rarely (project.md §2.3). Data pages get the API's 5-min rhythm + a long
// stale-while-revalidate; rarely-changing pages get an hour. Pages may override
// by setting their own Cache-Control (e.g. /status uses max-age=60). Requires a
// Cloudflare Cache Rule ("Cache Everything" / eligible) for HTML to be honored.

// Rarely-changing pages: docs, the hijri date, fuel prices (quarterly), privacy.
const LONG_LIVED =
  /^\/(ar\/|en\/)?(docs|docs\/reference|hijri-date|fuel-prices-egypt|privacy)\/?$/;

export const onRequest: MiddlewareHandler = async (context, next) => {
  const res = await next();
  const contentType = res.headers.get("Content-Type") ?? "";

  // Only HTML, and only when the page didn't set its own policy.
  if (contentType.includes("text/html") && !res.headers.has("Cache-Control")) {
    res.headers.set(
      "Cache-Control",
      LONG_LIVED.test(context.url.pathname)
        ? "public, max-age=3600, stale-while-revalidate=86400"
        : "public, max-age=300, stale-while-revalidate=86400",
    );
  }
  return res;
};
