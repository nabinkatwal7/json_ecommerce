"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { fmtMoney } from "@/lib/format";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";

export default function AdminOrderDetailPage() {
  const { id } = useParams();
  const oid = String(id);
  const qc = useQueryClient();
  const orders = useQuery({ queryKey: ["admin-orders"], queryFn: () => adminApi.orders() });
  const o = (orders.data ?? []).find((x) => x.id === oid);
  const tl = useQuery({
    queryKey: ["admin-order-timeline", oid],
    queryFn: () => adminApi.orderTimeline(oid),
    enabled: !!oid,
  });

  const cancel = useMutation({
    mutationFn: () => adminApi.cancelOrder(oid),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["admin-orders"] });
      toast.success("Cancelled");
    },
    onError: (e: Error) => toast.error(e.message),
  });
  const fulfill = useMutation({
    mutationFn: () => adminApi.fulfillOrder(oid),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["admin-orders"] }),
    onError: (e: Error) => toast.error(e.message),
  });
  const ship = useMutation({
    mutationFn: () => adminApi.shipOrder(oid),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["admin-orders"] }),
    onError: (e: Error) => toast.error(e.message),
  });

  if (orders.isLoading) return <p className="text-sm">Loading…</p>;
  if (!o) return <p className="text-sm">Not found</p>;

  return (
    <div className="space-y-6">
      <p className="text-sm">
        <Link href="/admin/orders" className="underline">
          ← Orders
        </Link>
      </p>
      <h1 className="font-mono text-lg">Order {o.id}</h1>
      <div className="flex flex-wrap gap-2">
        <Badge>{o.status}</Badge>
        <Badge variant="secondary">{o.paymentStatus}</Badge>
        <span className="text-sm">{fmtMoney(o.total)}</span>
      </div>
      <div className="flex flex-wrap gap-2">
        <Button size="sm" variant="outline" onClick={() => cancel.mutate()}>
          Cancel
        </Button>
        <Button size="sm" variant="outline" onClick={() => fulfill.mutate()}>
          Fulfill
        </Button>
        <Button size="sm" variant="outline" onClick={() => ship.mutate()}>
          Ship
        </Button>
      </div>
      <div className="rounded-md border p-4">
        <h2 className="mb-2 font-medium">Timeline</h2>
        <ul className="space-y-2 text-sm">
          {(tl.data?.events ?? []).map((e, i) => (
            <li key={i} className="border-b py-1">
              <span className="text-muted-foreground">{e.at}</span> — {e.label}{" "}
              <span className="text-xs">({e.kind})</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
