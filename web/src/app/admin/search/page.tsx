"use client";

import { useMutation } from "@tanstack/react-query";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function AdminSearchPage() {
  const reindex = useMutation({
    mutationFn: () => adminApi.searchReindex(),
    onSuccess: () => toast.success("Reindex requested"),
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Search</h1>
      <p className="text-sm text-muted-foreground">
        POST /admin/search/reindex — requires Algolia env on the API server.
      </p>
      <Button disabled={reindex.isPending} onClick={() => reindex.mutate()}>
        Reindex active products
      </Button>
    </div>
  );
}
