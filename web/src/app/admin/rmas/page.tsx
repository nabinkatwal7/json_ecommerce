"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { adminApi } from "@/lib/api";
import { Badge } from "@/components/ui/badge";

export default function AdminRmasPage() {
  const q = useQuery({ queryKey: ["admin-rmas"], queryFn: () => adminApi.rmas() });

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">RMAs</h1>
      <ul className="divide-y rounded-md border">
        {(q.data ?? []).map((r) => (
          <li key={r.id} className="flex items-center justify-between p-3 text-sm">
            <Link href={`/admin/rmas/${r.id}`} className="font-mono underline">
              {r.id.slice(0, 8)}…
            </Link>
            <Badge variant="outline">{r.status}</Badge>
          </li>
        ))}
      </ul>
    </div>
  );
}
