// Dynamic OG image per intent page, with the current number (add-seo-coverage
// stretch). `/og/<slug>.png` for EN, `?lang=ar` for Arabic. Edge-cached on the
// API's 5-min rhythm. Falls back to "—" if live data is briefly unavailable.
import type { APIRoute } from "astro";
import { intentBySlug } from "../../lib/intents";
import { getFX, getGold } from "../../lib/api";
import { renderOG } from "../../lib/og";

export const prerender = false;

const fmt = (n: number) => n.toLocaleString("en-US", { maximumFractionDigits: 2 });

export const GET: APIRoute = async ({ params, url }) => {
  const intent = intentBySlug(params.slug ?? "");
  if (!intent) return new Response("Not found", { status: 404 });

  const isAr = url.searchParams.get("lang") === "ar";
  const t = isAr ? intent.ar : intent.en;
  let big = "—";
  let unit: string | undefined;
  let sub = t.sub;

  try {
    if (intent.kind === "gold") {
      const g = await getGold();
      const v = g?.data?.world_derived?.find((x: any) => x.karat === intent.karat)?.value_egp ?? null;
      if (v != null) big = fmt(v);
      unit = isAr ? "جنيه/جرام" : "EGP/g";
    } else {
      const fx = await getFX();
      unit = isAr ? "جنيه" : "EGP";
      if (intent.kind === "parallel") {
        const p = fx?.data?.parallel;
        if (p && typeof p.min === "number" && p.n > 0) {
          big = `${fmt(p.min)}–${fmt(p.max)}`;
          sub = isAr ? `الموازي · ${p.n} مصادر` : `parallel · ${p.n} sources`;
        }
      } else {
        const v = fx?.data?.official?.[intent.quote!] ?? null;
        if (v != null) big = fmt(v);
      }
    }
  } catch {
    // live data temporarily unavailable — render the card with "—"
  }

  const png = await renderOG({ headline: t.h1, big, unit, sub, isAr });
  return new Response(png, {
    headers: {
      "Content-Type": "image/png",
      "Cache-Control": "public, max-age=300, stale-while-revalidate=86400",
    },
  });
};
