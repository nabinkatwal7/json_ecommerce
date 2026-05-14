"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function AdminTagEditPage() {
  const { id } = useParams();
  const tid = String(id);
  const qc = useQueryClient();
  const q = useQuery({ queryKey: ["admin-tags"], queryFn: () => adminApi.adminTags() });
  const t = (q.data ?? []).find((x) => x.id === tid);
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");

  useEffect(() => {
    if (t) {
      setName(t.name);
      setSlug(t.slug);
    }
  }, [t]);

  const save = useMutation({
    mutationFn: () => adminApi.putTag(tid, { name, slug }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["admin-tags"] });
      toast.success("Saved");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (!t) return <p className="text-sm">Loading…</p>;

  return (
    <div className="max-w-md space-y-4">
      <p className="text-sm">
        <Link href="/admin/tags" className="underline">
          ← Tags
        </Link>
      </p>
      <h1 className="text-xl font-semibold">Edit tag</h1>
      <div className="space-y-1">
        <Label>Name</Label>
        <Input value={name} onChange={(e) => setName(e.target.value)} />
      </div>
      <div className="space-y-1">
        <Label>Slug</Label>
        <Input value={slug} onChange={(e) => setSlug(e.target.value)} />
      </div>
      <Button disabled={save.isPending} onClick={() => save.mutate()}>
        Save
      </Button>
    </div>
  );
}
