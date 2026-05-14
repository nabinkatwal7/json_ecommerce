"use client";

import { useQuery } from "@tanstack/react-query";
import { adminApi } from "@/lib/api";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";

type LowResp = {
  threshold: number;
  lines: {
    productId: string;
    productName: string;
    variantId: string;
    sku: string;
    stock: number;
  }[];
};

export default function AdminLowStockPage() {
  const [th, setTh] = useState("");
  const q = useQuery({
    queryKey: ["admin-low-stock", th],
    queryFn: () => adminApi.lowStock(th ? Number(th) : undefined) as Promise<LowResp>,
  });

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Low stock</h1>
      <div className="max-w-xs space-y-1">
        <Label>Threshold override</Label>
        <Input value={th} onChange={(e) => setTh(e.target.value)} placeholder="default from server" />
      </div>
      <p className="text-sm text-muted-foreground">Server threshold: {q.data?.threshold}</p>
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="border-b">
            <th className="py-2">Product</th>
            <th>SKU</th>
            <th>Stock</th>
          </tr>
        </thead>
        <tbody>
          {(q.data?.lines ?? []).map((l) => (
            <tr key={`${l.productId}-${l.variantId}`} className="border-b">
              <td className="py-2">{l.productName}</td>
              <td>{l.sku}</td>
              <td>{l.stock}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
