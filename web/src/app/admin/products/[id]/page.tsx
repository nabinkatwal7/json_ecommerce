"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function AdminProductEditPage() {
  const { id } = useParams();
  const pid = String(id);
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["admin-products"],
    queryFn: () => adminApi.products(),
  });
  const p = (q.data ?? []).find((x) => x.id === pid);
  const [body, setBody] = useState("");

  useEffect(() => {
    if (p) setBody(JSON.stringify(p, null, 2));
  }, [p]);

  const save = useMutation({
    mutationFn: () => {
      const o = JSON.parse(body) as Record<string, unknown>;
      delete o.id;
      delete o.createdAt;
      delete o.updatedAt;
      return adminApi.putProduct(pid, o);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["admin-products"] });
      toast.success("Saved");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const del = useMutation({
    mutationFn: () => adminApi.deleteProduct(pid),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["admin-products"] });
      toast.success("Deleted");
      window.location.href = "/admin/products";
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (q.isLoading) return <p className="text-sm">Loading…</p>;
  if (!p) return <p className="text-sm">Not found</p>;

  return (
    <div className="mx-auto max-w-3xl space-y-4">
      <p className="text-sm">
        <Link href="/admin/products" className="underline">
          ← Products
        </Link>
      </p>
      <h1 className="text-2xl font-semibold">Edit product</h1>
      <div className="space-y-2">
        <Label>JSON (omit id / timestamps or they are stripped on save)</Label>
        <Textarea rows={22} value={body} onChange={(e) => setBody(e.target.value)} className="font-mono text-xs" />
      </div>
      <div className="flex gap-2">
        <Button disabled={save.isPending} onClick={() => save.mutate()}>
          Save
        </Button>
        <Button variant="destructive" disabled={del.isPending} onClick={() => del.mutate()}>
          Delete
        </Button>
      </div>
    </div>
  );
}
