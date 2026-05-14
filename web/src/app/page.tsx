"use client";

import { forwardRef, type ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { publicApi } from "@/lib/api";
import { ProductCard } from "@/components/product-card";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { SITE_NAME, SITE_TAGLINE } from "@/lib/site";
import { fmtMoney } from "@/lib/format";

const BannerHref = forwardRef<
  HTMLAnchorElement,
  { href: string; className?: string; children: ReactNode }
>(function BannerHref({ href, className, children }, ref) {
  if (href.startsWith("/")) {
    return (
      <Link ref={ref} href={href} className={className}>
        {children}
      </Link>
    );
  }
  return (
    <a ref={ref} href={href} className={className} target="_blank" rel="noopener noreferrer">
      {children}
    </a>
  );
});

export default function HomePage() {
  const home = useQuery({
    queryKey: ["storefront-feed"],
    queryFn: () => publicApi.storefrontFeed(),
    staleTime: 45_000,
  });

  const d = home.data;

  return (
    <div>
      <section className="border-b border-border/60 bg-muted/20">
        <div className="mx-auto max-w-6xl px-4 py-14 sm:px-6 sm:py-20">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">{SITE_NAME}</p>
          <h1 className="mt-3 max-w-2xl font-display text-4xl font-semibold leading-[1.1] tracking-tight text-foreground sm:text-5xl">
            Wear the colours. Live matchday.
          </h1>
          <p className="mt-4 max-w-lg text-base leading-relaxed text-muted-foreground sm:text-lg">
            {SITE_TAGLINE} New kits, terrace layers, and fan favourites — all in one place.
          </p>
          {home.isLoading ? (
            <div className="mt-8 flex flex-wrap gap-2">
              <Skeleton className="h-6 w-28 rounded-full" />
              <Skeleton className="h-6 w-36 rounded-full" />
              <Skeleton className="h-6 w-40 rounded-full" />
            </div>
          ) : d?.stats ? (
            <div className="mt-8 flex flex-wrap gap-2 text-xs font-medium text-muted-foreground">
              <span className="rounded-full border border-border/80 bg-background/80 px-3 py-1.5 shadow-sm">
                {d.stats.productCount} products live
              </span>
              <span className="rounded-full border border-border/80 bg-background/80 px-3 py-1.5 shadow-sm">
                {d.stats.categoryCount} categories
              </span>
              <span className="rounded-full border border-primary/20 bg-primary/[0.07] px-3 py-1.5 text-foreground shadow-sm">
                Free shipping over {fmtMoney(d.stats.freeShippingAtUsd)}
              </span>
            </div>
          ) : null}
          <div className="mt-8 flex flex-wrap gap-3">
            <Button size="lg" className="font-semibold shadow-sm" asChild>
              <Link href="/products">Shop all</Link>
            </Button>
            <Button size="lg" variant="outline" className="border-primary/20 font-semibold" asChild>
              <Link href="/categories">Browse categories</Link>
            </Button>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-6xl space-y-14 px-4 py-12 sm:px-6 sm:py-16">
        {home.isLoading && (
          <div className="space-y-4">
            <Skeleton className="h-24 w-full rounded-xl" />
            <div className="grid gap-5 sm:grid-cols-2">
              <Skeleton className="h-56 rounded-xl" />
              <Skeleton className="h-56 rounded-xl" />
            </div>
          </div>
        )}

        {home.isError && (
          <p className="rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
            We couldn&apos;t load the storefront. Check the API and try again.
          </p>
        )}

        {d?.announcementBanners?.map((b) => (
          <div
            key={b.id}
            className="flex flex-col gap-1 rounded-xl border border-primary/15 bg-primary/[0.06] px-5 py-4 sm:flex-row sm:items-center sm:justify-between sm:gap-6"
          >
            <div>
              <p className="font-display text-lg font-semibold text-foreground">{b.title}</p>
              <p className="mt-1 text-sm leading-relaxed text-muted-foreground">{b.body}</p>
            </div>
            {b.linkUrl ? (
              <Button variant="secondary" size="sm" className="shrink-0 self-start sm:self-center" asChild>
                <BannerHref href={b.linkUrl}>Shop kits</BannerHref>
              </Button>
            ) : null}
          </div>
        ))}

        {d?.carouselBanners && d.carouselBanners.length > 0 && (
          <section className="space-y-5">
            <div className="flex items-end justify-between gap-4">
              <h2 className="font-display text-2xl font-semibold tracking-tight sm:text-3xl">Matchday picks</h2>
            </div>
            <div className="grid gap-5 sm:grid-cols-2">
              {d.carouselBanners.map((b) => (
                <Card
                  key={b.id}
                  className="overflow-hidden border border-border/80 bg-card shadow-sm transition-shadow hover:shadow-md"
                >
                  {b.imageUrl ? (
                    <div className="aspect-[21/9] w-full overflow-hidden bg-muted sm:aspect-[2/1]">
                      {/* eslint-disable-next-line @next/next/no-img-element */}
                      <img
                        src={b.imageUrl}
                        alt={b.title || "Promo banner"}
                        className="h-full w-full object-cover"
                      />
                    </div>
                  ) : null}
                  <CardContent className="space-y-2 p-5">
                    <h3 className="font-display text-lg font-semibold">{b.title}</h3>
                    <p className="text-sm leading-relaxed text-muted-foreground">{b.body}</p>
                    {b.linkUrl ? (
                      <Button variant="link" className="h-auto px-0 font-semibold text-primary" asChild>
                        <BannerHref href={b.linkUrl}>View collection</BannerHref>
                      </Button>
                    ) : null}
                  </CardContent>
                </Card>
              ))}
            </div>
          </section>
        )}

        {d?.salePicks && d.salePicks.length > 0 && (
          <section className="space-y-6">
            <div className="flex flex-col gap-1 border-b border-border/60 pb-3 sm:flex-row sm:items-end sm:justify-between">
              <div>
                <h2 className="font-display text-2xl font-semibold tracking-tight sm:text-3xl">Sale rack</h2>
                <p className="text-sm text-muted-foreground">Tagged deals — same quality, sharper prices.</p>
              </div>
              <Link
                href={d.saleTagId ? `/products?tagId=${encodeURIComponent(d.saleTagId)}` : "/products"}
                className="text-sm font-semibold text-primary hover:underline"
              >
                View all sale
              </Link>
            </div>
            <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
              {d.salePicks.map((p) => (
                <ProductCard key={p.id} product={p} />
              ))}
            </div>
          </section>
        )}

        {home.isLoading && (
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="aspect-[4/5] rounded-xl" />
            ))}
          </div>
        )}

        {d?.featured && (
          <>
            <section className="space-y-6">
              <div className="flex items-end justify-between gap-4 border-b border-border/60 pb-3">
                <h2 className="font-display text-2xl font-semibold tracking-tight sm:text-3xl">New arrivals</h2>
                <Link href="/products" className="hidden text-sm font-semibold text-primary hover:underline sm:inline">
                  View all
                </Link>
              </div>
              <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
                {d.featured.newArrivals.map((p) => (
                  <ProductCard key={p.id} product={p} />
                ))}
              </div>
            </section>
            <section className="space-y-6">
              <div className="flex items-end justify-between gap-4 border-b border-border/60 pb-3">
                <h2 className="font-display text-2xl font-semibold tracking-tight sm:text-3xl">Best sellers</h2>
              </div>
              <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
                {d.featured.bestSellers.map((p) => (
                  <ProductCard key={p.id} product={p} />
                ))}
              </div>
            </section>
            <section className="space-y-5">
              <h2 className="border-b border-border/60 pb-3 font-display text-2xl font-semibold tracking-tight sm:text-3xl">
                Shop by category
              </h2>
              <div className="flex flex-wrap gap-2.5">
                {d.featured.featuredCategories.map((c) => (
                  <Link key={c.id} href={`/categories/${c.id}`}>
                    <Badge
                      variant="outline"
                      className="border-border/90 px-4 py-2 text-sm font-medium transition-colors hover:border-primary/40 hover:bg-muted/50"
                    >
                      {c.name}
                      <span className="ml-1.5 tabular-nums text-muted-foreground">({c.productCount})</span>
                    </Badge>
                  </Link>
                ))}
              </div>
            </section>
          </>
        )}
      </div>
    </div>
  );
}
