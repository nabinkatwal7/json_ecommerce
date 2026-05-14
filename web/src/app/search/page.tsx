"use client";

import { useQuery } from "@tanstack/react-query";
import { useSearchParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import Link from "next/link";
import { publicApi } from "@/lib/api";
import { ProductCard } from "@/components/product-card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Card } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

export default function SearchPage() {
  const sp = useSearchParams();
  const router = useRouter();
  const urlQ = sp.get("q") ?? "";
  const [draft, setDraft] = useState(urlQ);

  useEffect(() => {
    setDraft(urlQ);
  }, [urlQ]);

  const [debounced, setDebounced] = useState("");
  useEffect(() => {
    const id = setTimeout(() => setDebounced(draft.trim()), 280);
    return () => clearTimeout(id);
  }, [draft]);

  const suggest = useQuery({
    queryKey: ["search-suggest", debounced],
    queryFn: () => publicApi.searchSuggest(debounced, 10),
    enabled: debounced.length >= 2,
  });

  const res = useQuery({
    queryKey: ["search", urlQ],
    queryFn: () => publicApi.search(urlQ, undefined, 24),
    enabled: urlQ.length > 0,
  });

  const runSearch = (q: string) => {
    const t = q.trim();
    router.push(t ? `/search?q=${encodeURIComponent(t)}` : "/search");
  };

  return (
    <div className="mx-auto max-w-6xl px-4 pb-16 pt-10 sm:px-6">
      <header className="mb-8 max-w-2xl">
        <p className="text-xs font-semibold uppercase tracking-widest text-primary">Search</p>
        <h1 className="mt-2 font-display text-3xl font-semibold tracking-tight sm:text-4xl">Find your kit</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Live suggestions from the catalog API. Results use the same search index as the rest of the shop.
        </p>
      </header>

      <Card className="border border-border/80 p-4 shadow-sm sm:p-6">
        <form
          className="flex flex-col gap-3 sm:flex-row sm:items-end"
          onSubmit={(e) => {
            e.preventDefault();
            runSearch(draft);
          }}
        >
          <div className="min-w-0 flex-1 space-y-2">
            <Label htmlFor="q">Search the catalog</Label>
            <div className="relative">
              <Input
                id="q"
                autoComplete="off"
                placeholder="e.g. shirt, socks, mug…"
                value={draft}
                onChange={(e) => setDraft(e.target.value)}
                className="bg-background"
              />
              {debounced.length >= 2 && (suggest.data?.suggestions?.length ?? 0) > 0 && (
                <ul
                  className="absolute left-0 right-0 top-full z-40 mt-1 max-h-56 overflow-auto rounded-lg border border-border bg-popover py-1 text-sm shadow-md"
                  role="listbox"
                >
                  {suggest.data!.suggestions.map((s) => (
                    <li key={s}>
                      <button
                        type="button"
                        className="flex w-full px-3 py-2 text-left hover:bg-muted"
                        onClick={() => {
                          setDraft(s);
                          runSearch(s);
                        }}
                      >
                        {s}
                      </button>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
          <Button type="submit" className="shrink-0 sm:min-w-[7rem]">
            Search
          </Button>
        </form>
      </Card>

      {urlQ ? (
        <div className="mt-10 space-y-4">
          <div className="flex flex-wrap items-center justify-between gap-2 text-sm text-muted-foreground">
            <p>
              Results for <span className="font-medium text-foreground">&ldquo;{urlQ}&rdquo;</span>
            </p>
            {res.isFetching && <span>Searching…</span>}
            {!res.isFetching && res.data && (
              <span>
                <span className="tabular-nums text-foreground">{res.data.hits.length}</span> hits
              </span>
            )}
          </div>
          {res.isLoading && (
            <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
              {Array.from({ length: 8 }).map((_, i) => (
                <Skeleton key={i} className="aspect-[4/5] rounded-xl" />
              ))}
            </div>
          )}
          {!res.isLoading && res.data?.hits.length === 0 && (
            <p className="rounded-lg border border-dashed border-border px-4 py-10 text-center text-sm text-muted-foreground">
              No products matched. Try a shorter query or browse{" "}
              <Link href="/products" className="font-medium text-primary underline-offset-4 hover:underline">
                all kits
              </Link>
              .
            </p>
          )}
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {(res.data?.hits ?? []).map((p) => (
              <ProductCard key={p.id} product={p} />
            ))}
          </div>
        </div>
      ) : (
        <p className="mt-10 text-center text-sm text-muted-foreground">
          Type at least two characters to see suggestions, or press Search to run a full query.
        </p>
      )}
    </div>
  );
}
