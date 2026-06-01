// Server-side API client. In production this points at the in-cluster service
// (project.md §5); locally it defaults to the dev API. Every call fails soft —
// a down endpoint degrades one card, never the page (§2.6).
const API_BASE = import.meta.env.API_BASE ?? "http://localhost:8080";

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

export const API_PUBLIC_BASE =
  import.meta.env.PUBLIC_API_BASE ?? "https://api.naharda.com";
