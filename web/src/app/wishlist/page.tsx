"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { fmtMoney } from "@/lib/format";
import { toast } from "sonner";

export default function WishlistPage() {
  const { user, loading } = useAuth();
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["wishlist"],
    queryFn: () => customerApi.wishlist(),
    enabled: !!user,
  });

  const remove = useMutation({
    mutationFn: (it: { productId: string; variantId: string }) =>
      customerApi.deleteWishlist(it.productId, it.variantId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["wishlist"] }),
  });

  const move = useMutation({
    mutationFn: (it: { productId: string; variantId: string }) =>
      customerApi.wishlistMoveSave(it),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["wishlist"] });
      qc.invalidateQueries({ queryKey: ["save-later"] });
      toast.success("Moved to save for later");
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
    <div className="mx-auto max-w-3xl space-y-4 px-4 py-8">
      <h1 className="text-2xl font-semibold">Wishlist</h1>
      <ul className="space-y-2">
        {(q.data ?? []).map((it) => (
          <li key={`${it.productId}-${it.variantId}`} className="flex flex-wrap items-center justify-between gap-2 border p-3 text-sm">
            <div>
              <p className="font-medium">{it.name}</p>
              <p className="text-muted-foreground">{it.sku}</p>
              <p>{fmtMoney(it.price)}</p>
            </div>
            <div className="flex gap-2">
              <Button size="sm" variant="outline" onClick={() => move.mutate(it)}>
                Save for later
              </Button>
              <Button size="sm" variant="ghost" onClick={() => remove.mutate(it)}>
                Remove
              </Button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
