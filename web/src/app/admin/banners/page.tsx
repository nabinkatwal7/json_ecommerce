"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useState } from "react";
import { toast } from "sonner";

export default function AdminBannersPage() {
  const qc = useQueryClient();
  const q = useQuery({ queryKey: ["admin-banners"], queryFn: () => adminApi.banners() });
  const [slot, setSlot] = useState("home_carousel");
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [imageUrl, setImageUrl] = useState("");
  const [linkUrl, setLinkUrl] = useState("");
  const [sortOrder, setSortOrder] = useState("0");

  const create = useMutation({
    mutationFn: () =>
      adminApi.postBanner({
        slot,
        title,
        body,
        imageUrl,
        linkUrl,
        sortOrder: Number(sortOrder),
        isActive: true,
        startsAt: "",
        endsAt: "",
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["admin-banners"] });
      toast.success("Banner created");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const del = useMutation({
    mutationFn: (id: string) => adminApi.deleteBanner(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["admin-banners"] }),
  });

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Banners</h1>
      <div className="max-w-md space-y-2 rounded-md border p-4">
        <div className="space-y-1">
          <Label>Slot</Label>
          <Input value={slot} onChange={(e) => setSlot(e.target.value)} placeholder="home_carousel | announcement" />
        </div>
        <div className="space-y-1">
          <Label>Title</Label>
          <Input value={title} onChange={(e) => setTitle(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Body</Label>
          <Textarea value={body} onChange={(e) => setBody(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Image URL</Label>
          <Input value={imageUrl} onChange={(e) => setImageUrl(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Link URL</Label>
          <Input value={linkUrl} onChange={(e) => setLinkUrl(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Sort order</Label>
          <Input value={sortOrder} onChange={(e) => setSortOrder(e.target.value)} />
        </div>
        <Button disabled={create.isPending} onClick={() => create.mutate()}>
          Create
        </Button>
      </div>
      <ul className="divide-y rounded-md border">
        {(q.data ?? []).map((b) => (
          <li key={b.id} className="flex flex-wrap items-center justify-between gap-2 p-3 text-sm">
            <span>
              {b.title} <span className="text-muted-foreground">({b.slot})</span>
            </span>
            <div className="flex gap-2">
              <Button size="sm" variant="outline" asChild>
                <Link href={`/admin/banners/${b.id}`}>Edit</Link>
              </Button>
              <Button size="sm" variant="ghost" onClick={() => del.mutate(b.id)}>
                Delete
              </Button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
