"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { adminApi } from "@/lib/api";
import { fmtMoney } from "@/lib/format";
import { Badge } from "@/components/ui/badge";

export default function AdminOrdersPage() {
  const q = useQuery({ queryKey: ["admin-orders"], queryFn: () => adminApi.orders() });

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Orders</h1>
      <ul className="divide-y rounded-md border">
        {(q.data ?? []).map((o) => (
          <li key={o.id} className="flex flex-wrap items-center justify-between gap-2 p-3 text-sm">
            <Link href={`/admin/orders/${o.id}`} className="font-mono underline">
              {o.id.slice(0, 10)}…
            </Link>
            <div className="flex items-center gap-2">
              <Badge variant="outline">{o.status}</Badge>
              <span>{fmtMoney(o.total)}</span>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
