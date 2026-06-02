// Arabic names for the 13 canonical cities, keyed by slug (the API returns
// English names). Used by the Arabic detail pages + city picker.
export const cityNameAr: Record<string, string> = {
  cairo: "القاهرة",
  giza: "الجيزة",
  alexandria: "الإسكندرية",
  hurghada: "الغردقة",
  "sharm-el-sheikh": "شرم الشيخ",
  aswan: "أسوان",
  luxor: "الأقصر",
  mansoura: "المنصورة",
  tanta: "طنطا",
  asyut: "أسيوط",
  "port-said": "بورسعيد",
  suez: "السويس",
  ismailia: "الإسماعيلية",
};

// arCities maps an API city list to Arabic display names for the selector.
export function arCities(cities: { slug: string; name: string }[]) {
  return cities.map((c) => ({ slug: c.slug, name: cityNameAr[c.slug] ?? c.name }));
}
