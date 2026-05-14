"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { toast } from "sonner";

export default function AdminTagsPage() {
  const qc = useQueryClient();
  const q = useQuery({ queryKey: ["admin-tags"], queryFn: () => adminApi.adminTags() });
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");

  const create = useMutation({
    mutationFn: () => adminApi.postTag({ name, slug }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["admin-tags"] });
      toast.success("Tag created");
      setName("");
      setSlug("");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const del = useMutation({
    mutationFn: (id: string) => adminApi.deleteTag(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["admin-tags"] }),
  });

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Tags</h1>
      <div className="max-w-md flex flex-wrap gap-2 rounded-md border p-4">
        <div className="space-y-1">
          <Label>Name</Label>
          <Input value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Slug</Label>
          <Input value={slug} onChange={(e) => setSlug(e.target.value)} />
        </div>
        <Button className="self-end" disabled={create.isPending} onClick={() => create.mutate()}>
          Add
        </Button>
      </div>
      <ul className="divide-y rounded-md border">
        {(q.data ?? []).map((t) => (
          <li key={t.id} className="flex flex-wrap items-center justify-between gap-2 p-3 text-sm">
            <span>
              {t.name} <span className="text-muted-foreground">({t.slug})</span>
            </span>
            <div className="flex gap-2">
              <Button size="sm" variant="outline" asChild>
                <Link href={`/admin/tags/${t.id}`}>Edit</Link>
              </Button>
              <Button size="sm" variant="ghost" onClick={() => del.mutate(t.id)}>
                Delete
              </Button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
