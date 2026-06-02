// Server-side API client. In production this points at the in-cluster service
// (project.md §5); locally it defaults to the dev API. Every call fails soft —
// a down endpoint degrades one card, never the page (§2.6).
// Resolve bases at RUNTIME (Node SSR) so container env vars actually take
// effect — Vite would otherwise inline import.meta.env at build time. Order:
// process.env (runtime) → import.meta.env (build) → default.
const runtimeEnv = (k: string): string | undefined =>
  typeof process !== "undefined" ? process.env[k] : undefined;

const API_BASE =
  runtimeEnv("API_BASE") ?? import.meta.env.API_BASE ?? "http://localhost:8080";

async function get<T = any>(path: string): Promise<T | null> {
  try {
    const res = await fetch(`${API_BASE}${path}`, {
      headers: { Accept: "application/json" },
      signal: AbortSignal.timeout(4000),
    });
    if (!res.ok) {
      console.warn(`[api] ${path} → ${res.status}`);
      return null;
    }
    return (await res.json()) as T;
  } catch (err) {
    console.warn(`[api] ${path} fetch failed:`, (err as Error).message);
    return null;
  }
}

export const getFX = () => get("/v1/fx");
export const getGold = () => get("/v1/gold");
export const getFuel = () => get("/v1/fuel");
export const getWeather = (city = "cairo") => get(`/v1/weather/${city}`);
export const getAQI = (city = "cairo") => get(`/v1/aqi/${city}`);
export const getPrayer = (city = "cairo") => get(`/v1/prayer-times/${city}`);
export const getCalendar = () => get("/v1/calendar");
export const getStats = () => get("/v1/stats");

// pingPath performs a server-side health check: "down" on error/non-200,
// "degraded" when 200 but the data validator fails, else "up". Used by /status.
export async function pingPath(
  path: string,
  hasData?: (j: any) => boolean,
): Promise<"up" | "degraded" | "down"> {
  try {
    const res = await fetch(`${API_BASE}${path}`, { signal: AbortSignal.timeout(4000) });
    if (!res.ok) return "down";
    if (!hasData) return "up";
    return hasData(await res.json()) ? "up" : "degraded";
  } catch {
    return "down";
  }
}

export const API_PUBLIC_BASE =
  runtimeEnv("PUBLIC_API_BASE") ??
  import.meta.env.PUBLIC_API_BASE ??
  "https://api.naharda.com";
