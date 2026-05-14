"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { publicApi, adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function AdminCategoryEditPage() {
  const { id } = useParams();
  const cid = String(id);
  const qc = useQueryClient();
  const q = useQuery({ queryKey: ["categories"], queryFn: () => publicApi.categories() });
  const cat = (q.data ?? []).find((c) => c.id === cid);
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [desc, setDesc] = useState("");
  const [active, setActive] = useState(true);

  useEffect(() => {
    if (cat) {
      setName(cat.name);
      setSlug(cat.slug);
      setDesc(cat.description);
      setActive(cat.isActive);
    }
  }, [cat]);

  const save = useMutation({
    mutationFn: () =>
      adminApi.putCategory(cid, {
        name,
        slug,
        description: desc,
        isActive: active,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["categories"] });
      toast.success("Saved");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (!cat) return <p className="text-sm">Loading…</p>;

  return (
    <div className="max-w-md space-y-4">
      <p className="text-sm">
        <Link href="/admin/categories" className="underline">
          ← Categories
        </Link>
      </p>
      <h1 className="text-xl font-semibold">Edit category</h1>
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
      <div className="flex items-center gap-2">
        <Checkbox checked={active} onCheckedChange={(v) => setActive(v === true)} id="a" />
        <Label htmlFor="a">Active</Label>
      </div>
      <Button disabled={save.isPending} onClick={() => save.mutate()}>
        Save
      </Button>
    </div>
  );
}
