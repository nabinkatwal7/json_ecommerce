"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { publicApi, adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { toast } from "sonner";

export default function AdminCategoriesPage() {
  const qc = useQueryClient();
  const q = useQuery({ queryKey: ["categories"], queryFn: () => publicApi.categories() });
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [desc, setDesc] = useState("");

  const create = useMutation({
    mutationFn: () =>
      adminApi.postCategory({
        name,
        slug,
        description: desc,
        isActive: true,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["categories"] });
      toast.success("Created");
      setName("");
      setSlug("");
      setDesc("");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const del = useMutation({
    mutationFn: (id: string) => adminApi.deleteCategory(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["categories"] }),
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Categories</h1>
      <div className="max-w-md space-y-2 rounded-md border p-4">
        <h2 className="font-medium">Create</h2>
        <div className="space-y-1">
          <Label>Name</Label>
          <Input value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Slug</Label>
          <Input value={slug} onChange={(e) => setSlug(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Description</Label>
          <Input value={desc} onChange={(e) => setDesc(e.target.value)} />
        </div>
        <Button disabled={create.isPending} onClick={() => create.mutate()}>
          Create category
        </Button>
      </div>
      <ul className="divide-y rounded-md border">
        {(q.data ?? []).map((c) => (
          <li key={c.id} className="flex flex-wrap items-center justify-between gap-2 p-3 text-sm">
            <span>
              {c.name} <span className="text-muted-foreground">({c.slug})</span>
            </span>
            <div className="flex gap-2">
              <Button size="sm" variant="outline" asChild>
                <Link href={`/admin/categories/${c.id}`}>Edit</Link>
              </Button>
              <Button size="sm" variant="ghost" onClick={() => del.mutate(c.id)}>
                Delete
              </Button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
