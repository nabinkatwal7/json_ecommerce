"use client";

import { useQuery } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";
import { useMemo, useState } from "react";
import { publicApi } from "@/lib/api";
import { ProductCard } from "@/components/product-card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import type { Product } from "@/lib/types";
import { parseRFC3339Loose } from "@/lib/dates";

type SortKey = "featured" | "price-asc" | "price-desc" | "name" | "newest";

function minVariantPrice(p: Product) {
  return Math.min(...p.variants.map((v) => v.price));
}

function sortProducts(list: Product[], key: SortKey): Product[] {
  const out = [...list];
  switch (key) {
    case "price-asc":
      return out.sort((a, b) => minVariantPrice(a) - minVariantPrice(b));
    case "price-desc":
      return out.sort((a, b) => minVariantPrice(b) - minVariantPrice(a));
    case "name":
      return out.sort((a, b) => a.name.localeCompare(b.name));
    case "newest":
      return out.sort(
        (a, b) => parseRFC3339Loose(b.createdAt) - parseRFC3339Loose(a.createdAt),
      );
    default:
      return out;
  }
}

export default function ProductsPage() {
  const sp = useSearchParams();
  const router = useRouter();
  const categoryId = sp.get("categoryId") ?? "";
  const tagId = sp.get("tagId") ?? "";
  const [sort, setSort] = useState<SortKey>("featured");

  const cats = useQuery({
    queryKey: ["categories"],
    queryFn: () => publicApi.categories(),
  });
  const tags = useQuery({
    queryKey: ["tags"],
    queryFn: () => publicApi.tags(),
  });
  const q = useQuery({
    queryKey: ["products", categoryId, tagId],
    queryFn: () => publicApi.products(categoryId || undefined, tagId || undefined),
  });

  const sorted = useMemo(() => {
    const raw = q.data ?? [];
    return sortProducts(raw, sort);
  }, [q.data, sort]);

  const setParam = (key: string, value: string) => {
    const p = new URLSearchParams(sp.toString());
    if (!value) p.delete(key);
    else p.set(key, value);
    const qs = p.toString();
    router.push(qs ? `/products?${qs}` : "/products");
  };

  return (
    <div className="mx-auto max-w-6xl px-4 pb-16 pt-10 sm:px-6">
      <header className="mb-10 max-w-2xl">
        <p className="text-xs font-semibold uppercase tracking-widest text-primary">Catalog</p>
        <h1 className="mt-2 font-display text-3xl font-semibold tracking-tight sm:text-4xl">Kits &amp; gear</h1>
        <p className="mt-3 text-sm leading-relaxed text-muted-foreground sm:text-base">
          Filter by category or tag, sort the grid, and jump straight into product pages. Inventory and prices stay in
          sync with the API.
        </p>
      </header>

      <div className="rounded-xl border border-border/80 bg-muted/20 px-4 py-5 sm:px-5">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div className="grid flex-1 gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label className="text-xs font-medium uppercase tracking-wide text-muted-foreground">Category</Label>
              <Select
                value={categoryId || "__all__"}
                onValueChange={(v) => setParam("categoryId", v === "__all__" ? "" : v)}
              >
                <SelectTrigger className="bg-background">
                  <SelectValue placeholder="All" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__all__">All categories</SelectItem>
                  {(cats.data ?? []).map((c) => (
                    <SelectItem key={c.id} value={c.id}>
                      {c.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label className="text-xs font-medium uppercase tracking-wide text-muted-foreground">Sort</Label>
              <Select value={sort} onValueChange={(v) => setSort(v as SortKey)}>
                <SelectTrigger className="bg-background">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="featured">Featured order</SelectItem>
                  <SelectItem value="newest">Newest</SelectItem>
                  <SelectItem value="price-asc">Price: low to high</SelectItem>
                  <SelectItem value="price-desc">Price: high to low</SelectItem>
                  <SelectItem value="name">Name A–Z</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <div className="text-sm text-muted-foreground lg:text-right">
            {q.isFetching ? (
              <span>Updating…</span>
            ) : (
              <span>
                <span className="font-semibold tabular-nums text-foreground">{sorted.length}</span> results
              </span>
            )}
          </div>
        </div>

        <div className="mt-4 flex flex-wrap gap-2">
          <span className="mr-1 self-center text-xs font-medium text-muted-foreground">Tags:</span>
          <Button
            type="button"
            size="sm"
            variant={!tagId ? "secondary" : "outline"}
            className="h-8 rounded-full text-xs"
            onClick={() => setParam("tagId", "")}
          >
            All
          </Button>
          {(tags.data ?? []).map((t) => (
            <Button
              key={t.id}
              type="button"
              size="sm"
              variant={tagId === t.id ? "secondary" : "outline"}
              className="h-8 rounded-full text-xs"
              onClick={() => setParam("tagId", t.id === tagId ? "" : t.id)}
            >
              {t.name}
            </Button>
          ))}
        </div>
      </div>

      {q.isLoading && (
        <div className="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <Skeleton key={i} className="aspect-[4/5] rounded-xl" />
          ))}
        </div>
      )}

      {!q.isLoading && q.isError && (
        <p className="mt-10 rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          Could not load products.
        </p>
      )}

      {!q.isLoading && !q.isError && sorted.length === 0 && (
        <div className="mt-16 rounded-2xl border border-dashed border-border/80 bg-muted/15 px-6 py-16 text-center">
          <p className="font-display text-xl font-semibold text-foreground">Nothing matches</p>
          <p className="mt-2 text-sm text-muted-foreground">Try clearing a tag or switching category.</p>
          <Button className="mt-6" variant="secondary" onClick={() => router.push("/products")}>
            Reset filters
          </Button>
        </div>
      )}

      {!q.isLoading && sorted.length > 0 && (
        <div className="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
          {sorted.map((p) => (
            <ProductCard key={p.id} product={p} />
          ))}
        </div>
      )}
    </div>
  );
}
