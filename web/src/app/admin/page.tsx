"use client";

import { useQuery } from "@tanstack/react-query";
import { adminApi } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { fmtMoney } from "@/lib/format";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";

export default function AdminDashboardPage() {
  const [days, setDays] = useState(30);
  const q = useQuery({
    queryKey: ["admin-stats", days],
    queryFn: () => adminApi.dashboardStats(days),
  });

  const s = q.data;
  const maxRev = Math.max(
    s?.totalRevenue ?? 0,
    s?.previous.totalRevenue ?? 0,
    1,
  );

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Dashboard</h1>
      <div className="max-w-xs space-y-2">
        <Label>Days window</Label>
        <Input
          type="number"
          min={1}
          max={366}
          value={days}
          onChange={(e) => setDays(Number(e.target.value) || 30)}
        />
      </div>
      {q.isLoading && <p className="text-sm text-muted-foreground">Loading…</p>}
      {s && (
        <div className="grid gap-4 md:grid-cols-2">
          <Card className="border">
            <CardHeader>
              <CardTitle className="text-base">Revenue (paid)</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <p className="text-2xl font-semibold">{fmtMoney(s.totalRevenue)}</p>
              <p className="text-xs text-muted-foreground">Previous: {fmtMoney(s.previous.totalRevenue)}</p>
              <div className="flex h-4 gap-1 overflow-hidden rounded border bg-muted">
                <div
                  className="bg-primary"
                  style={{ width: `${(s.totalRevenue / maxRev) * 100}%` }}
                  title="Current"
                />
              </div>
            </CardContent>
          </Card>
          <Card className="border">
            <CardHeader>
              <CardTitle className="text-base">Orders & customers</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-2 text-sm">
              <div>
                <p className="text-muted-foreground">Orders placed</p>
                <p className="text-xl font-medium">{s.ordersPlaced}</p>
                <p className="text-xs">Prev {s.previous.ordersPlaced}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Paid orders</p>
                <p className="text-xl font-medium">{s.paidOrders}</p>
                <p className="text-xs">Prev {s.previous.paidOrders}</p>
              </div>
              <div className="col-span-2">
                <p className="text-muted-foreground">New customers</p>
                <p className="text-xl font-medium">{s.newCustomers}</p>
                <p className="text-xs">Prev {s.previous.newCustomers}</p>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
