// @ts-check
import { defineConfig } from "astro/config";
import node from "@astrojs/node";
import sitemap from "@astrojs/sitemap";

// Hybrid SSR in a k3s container (project.md §5, §11; add-dashboard design.md).
// Data pages render on demand and are edge-cached by Cloudflare; static pages
// opt in with `export const prerender = true`.
export default defineConfig({
  site: "https://naharda.com",
  output: "server",
  adapter: node({ mode: "standalone" }),
  i18n: {
    defaultLocale: "en",
    locales: ["en", "ar"],
    routing: { prefixDefaultLocale: false },
  },
  integrations: [
    // Per-language sitemap with daily lastmod; en ⇄ ar-EG alternates (add-seo-coverage).
    sitemap({
      i18n: { defaultLocale: "en", locales: { en: "en", ar: "ar-EG" } },
      changefreq: "daily",
      lastmod: new Date(),
    }),
  ],
  server: { port: 4321, host: true },
});
