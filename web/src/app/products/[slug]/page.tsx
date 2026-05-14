"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import Link from "next/link";
import { publicApi, customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ProductCard } from "@/components/product-card";
import { fmtMoney } from "@/lib/format";
import { useState, useEffect } from "react";
import { toast } from "sonner";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { ChevronRight } from "lucide-react";

export default function ProductDetailPage() {
  const params = useParams();
  const slug = String(params.slug);
  const { user } = useAuth();
  const qc = useQueryClient();
  const [variantId, setVariantId] = useState("");
  const [qty, setQty] = useState(1);

  const product = useQuery({
    queryKey: ["product", "slug", slug],
    queryFn: () => publicApi.productBySlug(slug),
  });

  useEffect(() => {
    if (product.data?.variants[0] && !variantId) {
      setVariantId(product.data.variants[0].id);
    }
  }, [product.data, variantId]);

  const category = useQuery({
    queryKey: ["category", product.data?.categoryId],
    queryFn: () => publicApi.category(product.data!.categoryId),
    enabled: !!product.data?.categoryId,
  });

  const related = useQuery({
    queryKey: ["related", product.data?.id],
    queryFn: () => publicApi.related(product.data!.id, 8),
    enabled: !!product.data?.id,
  });

  const add = useMutation({
    mutationFn: () =>
      customerApi.addCartItem({
        productId: product.data!.id,
        variantId,
        quantity: qty,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["cart"] });
      toast.success("Added to cart");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (product.isLoading) {
    return (
      <div className="mx-auto max-w-6xl px-4 py-10 sm:px-6">
        <Skeleton className="mb-6 h-4 w-48" />
        <div className="grid gap-10 lg:grid-cols-2">
          <Skeleton className="aspect-square rounded-2xl" />
          <div className="space-y-4">
            <Skeleton className="h-10 w-3/4" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-5/6" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-40" />
          </div>
        </div>
      </div>
    );
  }

  if (!product.data) {
    return (
      <div className="mx-auto max-w-6xl px-4 py-20 text-center sm:px-6">
        <p className="font-display text-xl font-semibold">Product not found</p>
        <p className="mt-2 text-sm text-muted-foreground">This slug may have changed. Browse the catalog instead.</p>
        <Button className="mt-6" asChild>
          <Link href="/products">Back to kits</Link>
        </Button>
      </div>
    );
  }

  const p = product.data;
  const v = p.variants.find((x) => x.id === variantId) ?? p.variants[0];
  const stockLabel =
    !v || v.stock <= 0 ? "Out of stock" : v.stock < 8 ? `Only ${v.stock} left` : "In stock";

  return (
    <div className="mx-auto max-w-6xl px-4 pb-16 pt-8 sm:px-6 sm:pt-10">
      <nav className="mb-8 flex flex-wrap items-center gap-1 text-xs text-muted-foreground sm:text-sm">
        <Link href="/" className="hover:text-foreground">
          Home
        </Link>
        <ChevronRight className="h-3.5 w-3.5 shrink-0 opacity-60" aria-hidden />
        <Link href="/products" className="hover:text-foreground">
          Kits
        </Link>
        {category.data ? (
          <>
            <ChevronRight className="h-3.5 w-3.5 shrink-0 opacity-60" aria-hidden />
            <Link href={`/products?categoryId=${encodeURIComponent(category.data.id)}`} className="hover:text-foreground">
              {category.data.name}
            </Link>
          </>
        ) : null}
        <ChevronRight className="h-3.5 w-3.5 shrink-0 opacity-60" aria-hidden />
        <span className="line-clamp-1 font-medium text-foreground">{p.name}</span>
      </nav>

      <div className="grid gap-10 lg:grid-cols-2 lg:gap-14">
        <div className="space-y-4">
          <div className="overflow-hidden rounded-2xl border border-border/80 bg-muted shadow-sm">
            {p.image ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img src={p.image} alt={p.name} className="aspect-square w-full object-cover" />
            ) : (
              <div className="flex aspect-square items-center justify-center text-sm text-muted-foreground">
                No image
              </div>
            )}
          </div>
        </div>

        <div className="space-y-6">
          <div>
            <h1 className="font-display text-3xl font-semibold leading-tight tracking-tight sm:text-4xl">{p.name}</h1>
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground sm:text-base">{p.description}</p>
          </div>

          {v && (
            <div className="flex flex-wrap items-end gap-3">
              <p className="text-3xl font-semibold tabular-nums tracking-tight text-foreground">{fmtMoney(v.price)}</p>
              <Badge
                variant={v.stock <= 0 ? "destructive" : v.stock < 8 ? "secondary" : "outline"}
                className="mb-1 font-medium"
              >
                {stockLabel}
              </Badge>
            </div>
          )}

          <Separator />

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label>Size / SKU</Label>
              <Select value={variantId} onValueChange={setVariantId}>
                <SelectTrigger className="bg-background">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {p.variants.map((x) => (
                    <SelectItem key={x.id} value={x.id}>
                      {x.size} · {x.color} ({x.sku})
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2">
              <Label>Quantity</Label>
              <Input
                type="number"
                min={1}
                max={v ? Math.max(1, v.stock) : 1}
                value={qty}
                onChange={(e) => setQty(Math.max(1, Number(e.target.value) || 1))}
                className="bg-background"
              />
            </div>
          </div>

          <div className="flex flex-wrap gap-3">
            {user ? (
              <Button
                size="lg"
                className="min-w-[10rem] font-semibold shadow-sm"
                disabled={!v || v.stock <= 0 || add.isPending}
                onClick={() => {
                  if (!v) return;
                  add.mutate();
                }}
              >
                Add to cart
              </Button>
            ) : (
              <Button size="lg" variant="secondary" className="font-semibold" asChild>
                <Link href="/login">Log in to purchase</Link>
              </Button>
            )}
            {user && (
              <>
                <Button
                  type="button"
                  variant="outline"
                  size="lg"
                  onClick={async () => {
                    try {
                      await customerApi.postWishlist({
                        productId: p.id,
                        variantId: variantId || p.variants[0]?.id,
                      });
                      toast.success("Saved to wishlist");
                      qc.invalidateQueries({ queryKey: ["wishlist"] });
                    } catch (e: unknown) {
                      toast.error(e instanceof Error ? e.message : "Failed");
                    }
                  }}
                >
                  Wishlist
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="lg"
                  onClick={async () => {
                    try {
                      await customerApi.postSaveLater({
                        productId: p.id,
                        variantId: variantId || p.variants[0]?.id,
                      });
                      toast.success("Saved for later");
                    } catch (e: unknown) {
                      toast.error(e instanceof Error ? e.message : "Failed");
                    }
                  }}
                >
                  Save for later
                </Button>
              </>
            )}
          </div>
        </div>
      </div>

      <section className="mt-16 border-t border-border/60 pt-12">
        <h2 className="font-display text-2xl font-semibold tracking-tight">You may also like</h2>
        <p className="mt-1 text-sm text-muted-foreground">Picked from related categories and tags via the API.</p>
        {related.isLoading ? (
          <div className="mt-8 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="aspect-[4/5] rounded-xl" />
            ))}
          </div>
        ) : (
          <div className="mt-8 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
            {(related.data ?? []).map((rp) => (
              <ProductCard key={rp.id} product={rp} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
