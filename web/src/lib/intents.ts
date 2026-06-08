// Tier-1 intent-page descriptors (add-seo-coverage). One descriptor per
// high-value query; the shared IntentPage component renders EN + ar-EG from it.

export type IntentKind = "currency" | "parallel" | "gold";

export interface IntentText {
  /** <title> — stable, no live number (number lives in H1/description). */
  title: string;
  description: string;
  /** Answer-first H1 — the query verbatim. */
  h1: string;
  /** Sub-label under the hero number. */
  sub: string;
  explainerH2: string;
  explainer: string;
  faq: { q: string; a: string }[];
  /** Breadcrumb leaf label. */
  crumb: string;
  /** Friendly name for the dataset / sparkline label. */
  datasetName: string;
}

export interface Intent {
  slug: string; // e.g. "eur-to-egp"
  kind: IntentKind;
  quote?: string; // currency code (lowercase key in /v1/fx official)
  karat?: number; // gold karat
  en: IntentText;
  ar: IntentText;
}

const cur = (
  code: string,
  enName: string,
  arName: string,
  enQ: string,
  arQ: string,
): Intent => ({
  slug: `${code.toLowerCase()}-to-egp`,
  kind: "currency",
  quote: code.toLowerCase(),
  en: {
    title: `${code} to EGP Today — Live ${enName} Rate | Naharda`,
    description: `The ${enName} to Egyptian pound today: the live official rate and 30-day history — sourced and available via API.`,
    h1: enQ,
    sub: "official",
    explainerH2: `About the ${code}/EGP rate`,
    explainer: `This is the official ${enName} (${code}) rate against the Egyptian pound, derived from public reference rates and refreshed continuously. Every value carries its source and fetch time, and the full history is available via the Naharda API.`,
    faq: [
      {
        q: enQ,
        a: `See the live number above. The official ${enName} rate is refreshed continuously and published with its source and timestamp.`,
      },
      {
        q: `Where does Naharda get the ${code} rate?`,
        a: `From public reference exchange-rate data, cross-checked and stored with provenance. The official Central Bank wiring is a production follow-up.`,
      },
    ],
    crumb: `${code} to EGP`,
    datasetName: `${enName} to EGP exchange rate`,
  },
  ar: {
    title: `${code} مقابل الجنيه النهاردة — سعر ${arName} | نهاردة`,
    description: `سعر ${arName} مقابل الجنيه المصري النهاردة: السعر الرسمي اللحظي وتاريخ 30 يوم — بمصدر ومتاح عبر API.`,
    h1: arQ,
    sub: "رسمي",
    explainerH2: `عن سعر ${code}/الجنيه`,
    explainer: `ده السعر الرسمي لـ${arName} (${code}) مقابل الجنيه المصري، محسوب من أسعار مرجعية عامة وبيتحدّث باستمرار. كل قيمة معاها المصدر ووقت الجلب، والتاريخ الكامل متاح عبر API نهاردة.`,
    faq: [
      {
        q: arQ,
        a: `بُص الرقم اللحظي فوق. السعر الرسمي بيتحدّث باستمرار وبيتنشر معاه المصدر ووقت التحديث.`,
      },
      {
        q: `نهاردة بتجيب سعر ${code} منين؟`,
        a: `من بيانات أسعار صرف مرجعية عامة، متراجَعة ومتخزّنة بمصدرها. ربط البنك المركزي الرسمي خطوة إنتاجية لاحقة.`,
      },
    ],
    crumb: `${code} مقابل الجنيه`,
    datasetName: `سعر ${arName} مقابل الجنيه`,
  },
});

const gold = (karat: number): Intent => ({
  slug: `gold-${karat}k`,
  kind: "gold",
  karat,
  en: {
    title: `Gold ${karat}k Price in Egypt Today (per gram) | Naharda`,
    description: `The price of ${karat}k gold in Egypt today, per gram (world-derived), with 30-day history — sourced and available via API.`,
    h1: `What is the price of ${karat}k gold in Egypt today?`,
    sub: "world-derived · EGP/gram",
    explainerH2: `About ${karat}k gold pricing`,
    explainer: `This is the world-derived ${karat}k gold price per gram in EGP — spot gold converted at the live USD/EGP rate and scaled to ${karat} karat. Egypt-retail prices add the local masna3eya (workmanship) premium and are tracked separately.`,
    faq: [
      {
        q: `What is the price of ${karat}k gold in Egypt today?`,
        a: `See the live per-gram number above. It is world-derived (spot × FX × karat); retail adds the local premium.`,
      },
      {
        q: `Why does the shop price differ from this?`,
        a: `Retail jewellers add a workmanship (masna3eya) premium on top of the world-derived value. Naharda separates the two so you can see both.`,
      },
    ],
    crumb: `Gold ${karat}k`,
    datasetName: `${karat}k gold price in Egypt (EGP/gram)`,
  },
  ar: {
    title: `سعر دهب عيار ${karat} في مصر النهاردة (للجرام) | نهاردة`,
    description: `سعر الدهب عيار ${karat} في مصر النهاردة للجرام (محسوب عالمياً) مع تاريخ 30 يوم — بمصدر ومتاح عبر API.`,
    h1: `سعر الدهب عيار ${karat} في مصر النهاردة كام؟`,
    sub: "محسوب عالمياً · جنيه/جرام",
    explainerH2: `عن تسعير دهب عيار ${karat}`,
    explainer: `ده سعر الدهب عيار ${karat} للجرام بالجنيه محسوب عالمياً — سعر الأونصة العالمي محوّل بسعر الدولار اللحظي ومضروب في العيار. أسعار محلات مصر بتزود عليه المصنعية وبتتسجّل لوحدها.`,
    faq: [
      {
        q: `سعر الدهب عيار ${karat} في مصر النهاردة كام؟`,
        a: `بُص رقم الجرام اللحظي فوق. ده محسوب عالمياً (أونصة × دولار × عيار)؛ سعر المحل بيزود المصنعية.`,
      },
      {
        q: `ليه سعر المحل مختلف عن ده؟`,
        a: `محلات الدهب بتضيف مصنعية فوق القيمة العالمية. نهاردة بتفصل الاتنين عشان تشوف القيمتين.`,
      },
    ],
    crumb: `دهب عيار ${karat}`,
    datasetName: `سعر دهب عيار ${karat} في مصر (جنيه/جرام)`,
  },
});

const parallelDollar: Intent = {
  slug: "black-market-dollar",
  kind: "parallel",
  quote: "usd",
  en: {
    title: "Black Market Dollar in Egypt Today — Parallel USD/EGP | Naharda",
    description:
      "The parallel (black-market) US dollar rate in Egypt today, published as an honest range with the number of sources — plus the official rate and API access.",
    h1: "What is the black-market dollar rate in Egypt today?",
    sub: "parallel market",
    explainerH2: "Official vs parallel",
    explainer:
      "Egypt has two dollar rates. The official rate is published by the Central Bank of Egypt. The parallel (black-market) rate trades outside official channels — Naharda publishes it as an honest range with the number of sources, never a fake-precise single value.",
    faq: [
      {
        q: "What is the black-market dollar rate in Egypt today?",
        a: "See the range above. The parallel rate is published as a min–max range with the number of contributing sources, not a single fake-precise number.",
      },
      {
        q: "Why are there two dollar rates in Egypt?",
        a: "The official rate is set/published by the Central Bank of Egypt; the parallel rate reflects supply and demand outside official channels.",
      },
    ],
    crumb: "Black-market dollar",
    datasetName: "Parallel-market USD/EGP rate",
  },
  ar: {
    title: "سعر الدولار في السوق السودا النهاردة — موازي USD/EGP | نهاردة",
    description:
      "سعر الدولار في السوق الموازية (السودا) في مصر النهاردة كنطاق أمين مع عدد المصادر — بالإضافة للسعر الرسمي والوصول عبر API.",
    h1: "سعر الدولار في السوق السودا في مصر النهاردة كام؟",
    sub: "السوق الموازية",
    explainerH2: "الرسمي مقابل الموازي",
    explainer:
      "مصر عندها سعرين للدولار. السعر الرسمي بينشره البنك المركزي. السعر الموازي (السوق السودا) بيتداول بره القنوات الرسمية — نهاردة بتنشره كنطاق أمين مع عدد المصادر، مش رقم واحد بدقة وهمية.",
    faq: [
      {
        q: "سعر الدولار في السوق السودا في مصر النهاردة كام؟",
        a: "بُص النطاق فوق. السعر الموازي بيتنشر كنطاق من الأدنى للأعلى مع عدد المصادر، مش رقم واحد بدقة وهمية.",
      },
      {
        q: "ليه فيه سعرين للدولار في مصر؟",
        a: "السعر الرسمي بيحدده/بينشره البنك المركزي؛ السعر الموازي بيعكس العرض والطلب بره القنوات الرسمية.",
      },
    ],
    crumb: "الدولار السوق السودا",
    datasetName: "سعر الدولار الموازي مقابل الجنيه",
  },
};

export const INTENTS: Intent[] = [
  cur("EUR", "Euro", "اليورو", "How much is the Euro in Egypt today?", "اليورو في مصر النهاردة بكام؟"),
  cur("SAR", "Saudi Riyal", "الريال السعودي", "How much is the Saudi Riyal in Egypt today?", "الريال السعودي في مصر النهاردة بكام؟"),
  cur("AED", "UAE Dirham", "الدرهم الإماراتي", "How much is the UAE Dirham in Egypt today?", "الدرهم الإماراتي في مصر النهاردة بكام؟"),
  cur("KWD", "Kuwaiti Dinar", "الدينار الكويتي", "How much is the Kuwaiti Dinar in Egypt today?", "الدينار الكويتي في مصر النهاردة بكام؟"),
  cur("GBP", "British Pound", "الجنيه الإسترليني", "How much is the British Pound in Egypt today?", "الجنيه الإسترليني في مصر النهاردة بكام؟"),
  parallelDollar,
  gold(21),
  gold(18),
  gold(24),
];

export const intentBySlug = (slug: string): Intent | undefined =>
  INTENTS.find((i) => i.slug === slug);
