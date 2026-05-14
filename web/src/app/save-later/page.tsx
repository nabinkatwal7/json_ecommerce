"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { fmtMoney } from "@/lib/format";
import { toast } from "sonner";

export default function SaveLaterPage() {
  const { user, loading } = useAuth();
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["save-later"],
    queryFn: () => customerApi.saveLater(),
    enabled: !!user,
  });

  const remove = useMutation({
    mutationFn: (it: { productId: string; variantId: string }) =>
      customerApi.deleteSaveLater(it.productId, it.variantId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["save-later"] }),
  });

  const move = useMutation({
    mutationFn: (it: { productId: string; variantId: string }) =>
      customerApi.saveLaterMoveWish(it),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["save-later"] });
      qc.invalidateQueries({ queryKey: ["wishlist"] });
      toast.success("Moved to wishlist");
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
      <h1 className="text-2xl font-semibold">Save for later</h1>
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
                Move to wishlist
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
