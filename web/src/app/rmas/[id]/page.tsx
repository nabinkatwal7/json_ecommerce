"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Badge } from "@/components/ui/badge";
import { fmtMoney } from "@/lib/format";

export default function RmaDetailPage() {
  const { id } = useParams();
  const rid = String(id);
  const { user, loading } = useAuth();
  const q = useQuery({
    queryKey: ["rma", rid],
    queryFn: () => customerApi.rma(rid),
    enabled: !!user && !!rid,
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
  if (!q.data) return <p className="p-8 text-sm">Loading…</p>;

  const r = q.data;

  return (
    <div className="mx-auto max-w-2xl space-y-4 px-4 py-8">
      <p className="text-sm">
        <Link href="/rmas" className="underline">
          ← RMAs
        </Link>
      </p>
      <h1 className="font-mono text-lg">RMA {r.id}</h1>
      <Badge>{r.status}</Badge>
      <p className="text-sm">{r.reason}</p>
      <ul className="text-sm">
        {r.items.map((it) => (
          <li key={it.variantId} className="border-b py-1">
            {it.name} ×{it.quantity} — {fmtMoney(it.price)}
          </li>
        ))}
      </ul>
      {r.refundAmount ? <p>Refund: {fmtMoney(r.refundAmount)}</p> : null}
      {r.adminNote ? <p className="text-sm text-muted-foreground">Note: {r.adminNote}</p> : null}
    </div>
  );
}
