"use client";

import { useQuery } from "@tanstack/react-query";
import { publicApi } from "@/lib/api";

export default function StatusPage() {
  const q = useQuery({
    queryKey: ["health"],
    queryFn: () => publicApi.health(),
    refetchInterval: 30_000,
  });

  return (
    <div className="mx-auto max-w-md px-4 py-12 text-sm">
      <h1 className="text-xl font-semibold">API health</h1>
      <p className="mt-4">
        GET /health:{" "}
        {q.isLoading ? "…" : q.data?.ok ? "ok" : "unreachable"}
      </p>
    </div>
  );
}
