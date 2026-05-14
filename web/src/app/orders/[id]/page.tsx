"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { fmtMoney } from "@/lib/format";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { useState } from "react";

export default function OrderDetailPage() {
  const { id } = useParams();
  const oid = String(id);
  const { user, loading } = useAuth();
  const qc = useQueryClient();
  const order = useQuery({
    queryKey: ["order", oid],
    queryFn: () => customerApi.order(oid),
    enabled: !!user && !!oid,
  });
  const [pi, setPi] = useState("");

  const cancel = useMutation({
    mutationFn: () => customerApi.cancelOrder(oid),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["order", oid] });
      toast.success("Cancelled");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const intent = useMutation({
    mutationFn: () => customerApi.stripeIntent(oid),
    onSuccess: (d) => {
      setPi(d.paymentIntentId);
      toast.message("PaymentIntent created", {
        description: `clientSecret: ${d.clientSecret.slice(0, 24)}…`,
      });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const pay = useMutation({
    mutationFn: (body: { stripePaymentIntentId?: string; stub?: boolean }) =>
      customerApi.payOrder(oid, body),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["order", oid] });
      toast.success("Paid");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  async function downloadInvoice() {
    try {
      const blob = await customerApi.orderInvoicePdf(oid);
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `invoice-${oid}.pdf`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : "PDF failed");
    }
  }

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
  if (!order.data) return <p className="p-8 text-sm">Loading…</p>;

  const o = order.data;

  return (
    <div className="mx-auto max-w-2xl space-y-6 px-4 py-8">
      <p className="text-sm">
        <Link href="/orders" className="underline">
          ← Orders
        </Link>
      </p>
      <h1 className="font-mono text-lg">Order {o.id}</h1>
      <div className="flex flex-wrap gap-2">
        <Badge>{o.status}</Badge>
        <Badge variant="secondary">{o.paymentStatus}</Badge>
      </div>
      <p className="text-sm">Total {fmtMoney(o.total)}</p>
      <ul className="space-y-2 text-sm">
        {o.items.map((it) => (
          <li key={it.variantId} className="flex justify-between border-b py-1">
            <span>
              {it.name} ×{it.quantity}
            </span>
            <span>{fmtMoney(it.price * it.quantity)}</span>
          </li>
        ))}
      </ul>
      <div className="flex flex-wrap gap-2">
        <Button type="button" variant="outline" size="sm" onClick={downloadInvoice}>
          Download invoice PDF
        </Button>
        {o.status === "created" && o.paymentStatus === "pending" && (
          <Button type="button" variant="outline" size="sm" onClick={() => cancel.mutate()}>
            Cancel order
          </Button>
        )}
      </div>
      {o.status === "created" && o.paymentStatus === "pending" && (
        <div className="space-y-3 border-t pt-4">
          <p className="text-sm font-medium">Pay</p>
          <Button type="button" size="sm" variant="outline" onClick={() => intent.mutate()}>
            Create Stripe PaymentIntent
          </Button>
          <div className="space-y-1">
            <Label htmlFor="pi">Stripe PaymentIntent ID</Label>
            <Input id="pi" value={pi} onChange={(e) => setPi(e.target.value)} />
          </div>
          <Button
            type="button"
            size="sm"
            onClick={() => pay.mutate({ stripePaymentIntentId: pi || undefined })}
          >
            Pay with Stripe ID
          </Button>
          <Button type="button" size="sm" variant="secondary" onClick={() => pay.mutate({ stub: true })}>
            Pay (dev stub)
          </Button>
        </div>
      )}
    </div>
  );
}
