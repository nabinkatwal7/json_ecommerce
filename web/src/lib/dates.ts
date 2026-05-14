/** Parse RFC3339-ish timestamps for client-side sorting; returns ms or 0. */
export function parseRFC3339Loose(s: string): number {
  const t = Date.parse(s);
  return Number.isFinite(t) ? t : 0;
}
