"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export default function AdminProductsPage() {
  const q = useQuery({ queryKey: ["admin-products"], queryFn: () => adminApi.products() });

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Products</h1>
        <Button asChild size="sm">
          <Link href="/admin/products/new">New product</Link>
        </Button>
      </div>
      <ul className="divide-y rounded-md border">
        {(q.data ?? []).map((p) => (
          <li key={p.id} className="flex flex-wrap items-center justify-between gap-2 p-3 text-sm">
            <div>
              <Link href={`/admin/products/${p.id}`} className="font-medium underline">
                {p.name}
              </Link>
              <p className="text-xs text-muted-foreground">{p.slug}</p>
            </div>
            <Badge variant={p.isActive ? "default" : "secondary"}>
              {p.isActive ? "active" : "inactive"}
            </Badge>
          </li>
        ))}
      </ul>
    </div>
  );
}
