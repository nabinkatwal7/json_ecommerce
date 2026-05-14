import { Suspense } from "react";

export default function SearchLayout({ children }: { children: React.ReactNode }) {
  return <Suspense fallback={<p className="p-8 text-sm">Loading…</p>}>{children}</Suspense>;
}
