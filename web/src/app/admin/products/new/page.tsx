"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useState } from "react";
import { toast } from "sonner";

const defaultJson = `{
  "name": "New product",
  "slug": "new-product",
  "description": "",
  "image": "",
  "categoryId": "",
  "tags": [],
  "tagIds": [],
  "variants": [
    { "sku": "SKU-1", "size": "", "color": "", "price": 9.99, "stock": 10, "weightKg": 0.5 }
  ],
  "isActive": true
}`;

export default function AdminProductNewPage() {
  const router = useRouter();
  const qc = useQueryClient();
  const [body, setBody] = useState(defaultJson);
  const create = useMutation({
    mutationFn: () => adminApi.postProduct(JSON.parse(body) as Record<string, unknown>),
    onSuccess: (p) => {
      qc.invalidateQueries({ queryKey: ["admin-products"] });
      toast.success("Created");
      router.push(`/admin/products/${p.id}`);
    },
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <div className="mx-auto max-w-3xl space-y-4">
      <p className="text-sm">
        <Link href="/admin/products" className="underline">
          ← Products
        </Link>
      </p>
      <h1 className="text-2xl font-semibold">New product</h1>
      <p className="text-xs text-muted-foreground">
        JSON body matches admin API: name, slug, description, image, categoryId, tags, tagIds,
        variants (sku, size, color, price, stock, optional weightKg), isActive.
      </p>
      <div className="space-y-2">
        <Label>JSON</Label>
        <Textarea rows={18} value={body} onChange={(e) => setBody(e.target.value)} className="font-mono text-xs" />
      </div>
      <Button disabled={create.isPending} onClick={() => create.mutate()}>
        Create
      </Button>
    </div>
  );
}
