"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { SiteHeader } from "@/components/site-header";
import { SITE_NAME, SITE_TAGLINE } from "@/lib/site";

export function AppChrome({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const admin = pathname?.startsWith("/admin");

  if (admin) {
    return <div className="min-h-screen">{children}</div>;
  }

  const year = new Date().getFullYear();

  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader />
      <main className="flex-1 pb-10">{children}</main>
      <footer className="mt-auto border-t border-border/80 bg-muted/25">
        <div className="mx-auto grid max-w-6xl gap-8 px-4 py-12 sm:grid-cols-2 sm:px-6">
          <div>
            <p className="font-display text-xl font-semibold text-foreground">{SITE_NAME}</p>
            <p className="mt-1 max-w-xs text-sm leading-relaxed text-muted-foreground">{SITE_TAGLINE}</p>
          </div>
          <div className="flex flex-col gap-3 text-sm text-muted-foreground sm:items-end sm:text-right">
            <div className="flex flex-wrap gap-x-4 gap-y-1 sm:justify-end">
              <Link href="/products" className="hover:text-foreground">
                Kits
              </Link>
              <Link href="/categories" className="hover:text-foreground">
                Categories
              </Link>
              <Link href="/search" className="hover:text-foreground">
                Search
              </Link>
            </div>
            <div className="flex flex-wrap gap-x-4 gap-y-1 sm:justify-end">
              <Link href="/status" className="hover:text-foreground">
                Status
              </Link>
              <Link href="/internal/cron" className="hover:text-foreground">
                Tools
              </Link>
            </div>
            <p className="text-xs text-muted-foreground/80">© {year} {SITE_NAME}</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
