"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { fmtMoney } from "@/lib/format";

export default function InsightsPage() {
  const { user, loading } = useAuth();
  const q = useQuery({
    queryKey: ["insights"],
    queryFn: () => customerApi.insights(),
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
  if (!q.data) return <p className="p-8 text-sm">Loading…</p>;

  return (
    <div className="mx-auto max-w-lg space-y-4 px-4 py-8">
      <h1 className="text-2xl font-semibold">Insights</h1>
      <Card className="border">
        <CardHeader>
          <CardTitle className="text-base">Spend</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <p>Paid orders: {q.data.paidOrders}</p>
          <p>Lifetime spend: {fmtMoney(q.data.lifetimeSpend)}</p>
          <p>Big spender threshold: {fmtMoney(q.data.bigSpenderThresholdUsd)}</p>
        </CardContent>
      </Card>
      <Card className="border">
        <CardHeader>
          <CardTitle className="text-base">Segments</CardTitle>
        </CardHeader>
        <CardContent className="text-sm">
          {(q.data.user.segments ?? []).length === 0 ? (
            <p className="text-muted-foreground">No segments</p>
          ) : (
            <ul className="list-inside list-disc">
              {(q.data.user.segments ?? []).map((s) => (
                <li key={s}>{s}</li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
