"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { publicApi } from "@/lib/api";
import { ProductCard } from "@/components/product-card";

export default function CategoryDetailPage() {
  const { id } = useParams();
  const cat = useQuery({
    queryKey: ["category", id],
    queryFn: () => publicApi.category(String(id)),
  });
  const products = useQuery({
    queryKey: ["products", id],
    queryFn: () => publicApi.products(String(id)),
    enabled: !!id,
  });

  if (cat.isLoading) return <p className="p-8 text-sm">Loading…</p>;
  if (!cat.data) return <p className="p-8 text-sm">Not found</p>;

  return (
    <div className="mx-auto max-w-6xl space-y-6 px-4 py-8">
      <div>
        <h1 className="text-2xl font-semibold">{cat.data.name}</h1>
        <p className="text-sm text-muted-foreground">{cat.data.description}</p>
        <p className="pt-2 text-sm">
          <Link href={`/products?categoryId=${cat.data.id}`} className="underline">
            View as product grid
          </Link>
        </p>
      </div>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {(products.data ?? []).map((p) => (
          <ProductCard key={p.id} product={p} />
        ))}
      </div>
    </div>
  );
}
