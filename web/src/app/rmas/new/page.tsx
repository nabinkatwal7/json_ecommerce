"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useState } from "react";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

export default function RmaNewPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const qc = useQueryClient();
  const orders = useQuery({
    queryKey: ["orders"],
    queryFn: () => customerApi.orders(),
    enabled: !!user,
  });
  const [orderId, setOrderId] = useState("");
  const [reason, setReason] = useState("");
  const [lines, setLines] = useState(
    "[]" as string, // JSON array of {productId, variantId, quantity}
  );

  const create = useMutation({
    mutationFn: () =>
      customerApi.postRma({
        orderId,
        reason,
        items: JSON.parse(lines) as {
          productId: string;
          variantId: string;
          quantity: number;
        }[],
      }),
    onSuccess: (r) => {
      qc.invalidateQueries({ queryKey: ["rmas"] });
      toast.success("RMA created");
      router.push(`/rmas/${r.id}`);
    },
    onError: (e: Error) => toast.error(e.message),
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
    <div className="mx-auto max-w-xl space-y-4 px-4 py-8">
      <h1 className="text-2xl font-semibold">New RMA</h1>
      <p className="text-xs text-muted-foreground">
        Pick a paid order that is fulfilled or shipped. Items must be JSON array of objects with
        productId, variantId, and quantity.
      </p>
      <div className="space-y-2">
        <Label>Order</Label>
        <select
          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
          value={orderId}
          onChange={(e) => setOrderId(e.target.value)}
        >
          <option value="">Select…</option>
          {(orders.data ?? [])
            .filter(
              (o) =>
                o.paymentStatus === "paid" &&
                (o.status === "fulfilled" || o.status === "shipped" || o.status === "paid"),
            )
            .map((o) => (
              <option key={o.id} value={o.id}>
                {o.id.slice(0, 8)} — {o.status}
              </option>
            ))}
        </select>
      </div>
      <div className="space-y-2">
        <Label>Reason</Label>
        <Input value={reason} onChange={(e) => setReason(e.target.value)} />
      </div>
      <div className="space-y-2">
        <Label>Items JSON</Label>
        <Textarea rows={6} value={lines} onChange={(e) => setLines(e.target.value)} />
      </div>
      <Button disabled={create.isPending} onClick={() => create.mutate()}>
        Submit
      </Button>
    </div>
  );
}
