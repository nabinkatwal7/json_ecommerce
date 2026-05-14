"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { fmtMoney } from "@/lib/format";
import { Badge } from "@/components/ui/badge";

export default function OrdersPage() {
  const { user, loading } = useAuth();
  const q = useQuery({
    queryKey: ["orders"],
    queryFn: () => customerApi.orders(),
    enabled: !!user,
  });

  if (loading) return null;
  if (!user) {
    return (
      <p className="p-8 text-center text-sm">
        <Link href="/login" className="underline">
          Log in
        </Link>
      </p>
    );
  }

  return (
    <div className="mx-auto max-w-4xl space-y-4 px-4 py-8">
      <h1 className="text-2xl font-semibold">Orders</h1>
      <ul className="space-y-2">
        {(q.data ?? []).map((o) => (
          <li key={o.id} className="flex flex-wrap items-center justify-between gap-2 border p-3">
            <div>
              <Link href={`/orders/${o.id}`} className="font-mono text-sm underline">
                {o.id.slice(0, 8)}…
              </Link>
              <p className="text-xs text-muted-foreground">{o.createdAt}</p>
            </div>
            <div className="text-right text-sm">
              <Badge variant="outline">{o.status}</Badge>{" "}
              <Badge variant="secondary">{o.paymentStatus}</Badge>
              <p className="font-medium">{fmtMoney(o.total)}</p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
