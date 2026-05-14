"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Badge } from "@/components/ui/badge";

export default function RmasListPage() {
  const { user, loading } = useAuth();
  const q = useQuery({
    queryKey: ["rmas"],
    queryFn: () => customerApi.rmas(),
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

  return (
    <div className="mx-auto max-w-3xl space-y-4 px-4 py-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Returns (RMA)</h1>
        <Link href="/rmas/new" className="text-sm underline">
          New request
        </Link>
      </div>
      <ul className="space-y-2">
        {(q.data ?? []).map((r) => (
          <li key={r.id} className="flex items-center justify-between border p-3 text-sm">
            <Link href={`/rmas/${r.id}`} className="font-mono underline">
              {r.id.slice(0, 8)}…
            </Link>
            <Badge variant="outline">{r.status}</Badge>
          </li>
        ))}
      </ul>
    </div>
  );
}
