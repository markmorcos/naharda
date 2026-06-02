// Dynamic OG card renderer (add-seo-coverage stretch): satori (HTML→SVG) +
// resvg (SVG→PNG). Text is vectorized by satori, so resvg needs no fonts. Fonts
// are the project's IBM Plex woff files read from node_modules (present at
// runtime per the web Dockerfile). Loaded once at module init.
import { createRequire } from "node:module";
import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import satori from "satori";
import { Resvg } from "@resvg/resvg-js";

const require = createRequire(import.meta.url);
const fontFile = (pkg: string, file: string) =>
  readFileSync(join(dirname(require.resolve(`${pkg}/package.json`)), "files", file));

// satori supports woff (not woff2). Latin + Arabic, 400/600.
const FONTS = [
  { name: "IBM Plex Sans", weight: 400 as const, style: "normal" as const,
    data: fontFile("@fontsource/ibm-plex-sans", "ibm-plex-sans-latin-400-normal.woff") },
  { name: "IBM Plex Sans", weight: 600 as const, style: "normal" as const,
    data: fontFile("@fontsource/ibm-plex-sans", "ibm-plex-sans-latin-600-normal.woff") },
  { name: "IBM Plex Sans Arabic", weight: 400 as const, style: "normal" as const,
    data: fontFile("@fontsource/ibm-plex-sans-arabic", "ibm-plex-sans-arabic-arabic-400-normal.woff") },
  { name: "IBM Plex Sans Arabic", weight: 600 as const, style: "normal" as const,
    data: fontFile("@fontsource/ibm-plex-sans-arabic", "ibm-plex-sans-arabic-arabic-600-normal.woff") },
];

// Brand colors (tokens.css dark theme).
const C = { bg: "#0b2a33", ink: "#f4f1ea", brand: "#7fb2be", muted: "#8a9ba0", accent: "#d4a017" };

export interface OGCard {
  headline: string;
  big: string;
  unit?: string;
  sub: string;
  updated?: string;
  isAr?: boolean;
}

// Minimal hyperscript for satori (no JSX in a .ts lib).
function h(type: string, style: Record<string, unknown>, children?: unknown): unknown {
  return { type, props: { style, children } };
}

export async function renderOG(card: OGCard): Promise<Buffer> {
  const family = card.isAr ? "IBM Plex Sans Arabic" : "IBM Plex Sans";
  const bigSize = card.big.length > 7 ? 100 : 150;

  const tree = h(
    "div",
    {
      width: "1200px", height: "630px", display: "flex", flexDirection: "column",
      justifyContent: "space-between", padding: "64px 72px", backgroundColor: C.bg,
      color: C.ink, fontFamily: family, direction: card.isAr ? "rtl" : "ltr",
    },
    [
      // Header: wordmark + dot · domain
      h("div", { display: "flex", alignItems: "center", justifyContent: "space-between" }, [
        h("div", { display: "flex", alignItems: "center", gap: "14px" }, [
          h("div", { fontSize: "36px", fontWeight: 600, color: C.brand, fontFamily: "IBM Plex Sans" }, "naharda"),
          h("div", { width: "18px", height: "18px", borderRadius: "9px", backgroundColor: C.accent }),
        ]),
        h("div", { fontSize: "24px", color: C.muted, direction: "ltr", fontFamily: "IBM Plex Sans" }, "naharda.com"),
      ]),
      // Main: headline + big value
      h("div", { display: "flex", flexDirection: "column" }, [
        h("div", { fontSize: "40px", fontWeight: 600, lineHeight: 1.15, maxWidth: "1000px" }, card.headline),
        h("div", { display: "flex", alignItems: "flex-end", gap: "18px", marginTop: "20px" }, [
          h("div", { fontSize: `${bigSize}px`, fontWeight: 600, lineHeight: 1, direction: "ltr", fontFamily: "IBM Plex Sans" }, card.big),
          card.unit ? h("div", { fontSize: "40px", color: C.muted, paddingBottom: "16px" }, card.unit) : null,
        ].filter(Boolean)),
      ]),
      // Footer: sub-label · updated
      h("div", { display: "flex", alignItems: "center", justifyContent: "space-between" }, [
        h("div", { fontSize: "30px", fontWeight: 600, color: C.accent }, card.sub),
        h("div", { fontSize: "22px", color: C.muted, direction: "ltr", fontFamily: "IBM Plex Sans" }, card.updated ?? ""),
      ]),
    ],
  );

  const svg = await satori(tree as Parameters<typeof satori>[0], { width: 1200, height: 630, fonts: FONTS });
  return Buffer.from(new Resvg(svg, { fitTo: { mode: "width", value: 1200 } }).render().asPng());
}
