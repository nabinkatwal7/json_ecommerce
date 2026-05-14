"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { cn } from "@/lib/utils";

export function CartNavLink({ className }: { className?: string }) {
  const { user } = useAuth();
  const cart = useQuery({
    queryKey: ["cart"],
    queryFn: () => customerApi.cart(),
    enabled: !!user,
    staleTime: 15_000,
  });
  const n = cart.data?.items.reduce((s, i) => s + i.quantity, 0) ?? 0;
  return (
    <Link href="/cart" className={cn("relative inline-flex items-center gap-1.5", className)}>
      <span>Cart</span>
      {n > 0 ? (
        <span className="rounded-full bg-primary px-1.5 py-0.5 text-[10px] font-bold tabular-nums leading-none text-primary-foreground">
          {n > 99 ? "99+" : n}
        </span>
      ) : null}
    </Link>
  );
}
