"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { fmtMoney } from "@/lib/format";
import { toast } from "sonner";

export default function CartPage() {
  const { user, loading } = useAuth();
  const qc = useQueryClient();
  const cart = useQuery({
    queryKey: ["cart"],
    queryFn: () => customerApi.cart(),
    enabled: !!user,
  });

  const patch = useMutation({
    mutationFn: ({ id, qty }: { id: string; qty: number }) =>
      customerApi.patchCartItem(id, qty),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["cart"] }),
  });

  const del = useMutation({
    mutationFn: (id: string) => customerApi.deleteCartItem(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["cart"] }),
  });

  if (loading) return null;
  if (!user) {
    return (
      <p className="p-8 text-center text-sm">
        <Link href="/login" className="underline">
          Log in
        </Link>{" "}
        to view your cart.
      </p>
    );
  }

  const lines = cart.data?.items ?? [];
  const subtotal = lines.reduce((s, l) => s + l.price * l.quantity, 0);

  return (
    <div className="mx-auto max-w-3xl space-y-6 px-4 py-8">
      <h1 className="text-2xl font-semibold">Cart</h1>
      <div className="flex flex-wrap gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={async () => {
            try {
              const r = await customerApi.cartValidate();
              toast.message(r.ok ? "Cart OK" : "Issues found", {
                description: r.lines
                  .filter((l) => !l.ok)
                  .map((l) => `${l.name}: ${l.issues.join(", ")}`)
                  .join(" · "),
              });
            } catch (e: unknown) {
              toast.error(e instanceof Error ? e.message : "Validate failed");
            }
          }}
        >
          Validate cart
        </Button>
        <Button size="sm" asChild>
          <Link href="/checkout">Checkout</Link>
        </Button>
      </div>
      {lines.length === 0 ? (
        <p className="text-sm text-muted-foreground">Your cart is empty.</p>
      ) : (
        <ul className="space-y-4 border-t pt-4">
          {lines.map((l) => (
            <li key={l.id} className="flex flex-wrap items-center gap-3 border-b pb-4">
              <div className="min-w-0 flex-1">
                <p className="font-medium">{l.name}</p>
                <p className="text-xs text-muted-foreground">{l.sku}</p>
                <p className="text-sm">{fmtMoney(l.price)} each</p>
              </div>
              <Input
                className="w-20"
                type="number"
                min={1}
                defaultValue={l.quantity}
                onBlur={(e) => {
                  const n = Number(e.target.value);
                  if (n > 0) patch.mutate({ id: l.id, qty: n });
                }}
              />
              <Button variant="ghost" size="sm" onClick={() => del.mutate(l.id)}>
                Remove
              </Button>
            </li>
          ))}
        </ul>
      )}
      <p className="text-right text-lg font-medium">Subtotal {fmtMoney(subtotal)}</p>
    </div>
  );
}
