"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function AdminBannerEditPage() {
  const { id } = useParams();
  const bid = String(id);
  const qc = useQueryClient();
  const q = useQuery({ queryKey: ["admin-banners"], queryFn: () => adminApi.banners() });
  const b = (q.data ?? []).find((x) => x.id === bid);
  const [slot, setSlot] = useState("");
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [imageUrl, setImageUrl] = useState("");
  const [linkUrl, setLinkUrl] = useState("");
  const [sortOrder, setSortOrder] = useState("0");
  const [active, setActive] = useState(true);
  const [startsAt, setStartsAt] = useState("");
  const [endsAt, setEndsAt] = useState("");

  useEffect(() => {
    if (b) {
      setSlot(b.slot);
      setTitle(b.title);
      setBody(b.body);
      setImageUrl(b.imageUrl);
      setLinkUrl(b.linkUrl);
      setSortOrder(String(b.sortOrder));
      setActive(b.isActive);
      setStartsAt(b.startsAt ?? "");
      setEndsAt(b.endsAt ?? "");
    }
  }, [b]);

  const save = useMutation({
    mutationFn: () =>
      adminApi.putBanner(bid, {
        slot,
        title,
        body,
        imageUrl,
        linkUrl,
        sortOrder: Number(sortOrder),
        isActive: active,
        startsAt,
        endsAt,
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["admin-banners"] });
      toast.success("Saved");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (!b) return <p className="text-sm">Loading…</p>;

  return (
    <div className="max-w-md space-y-3">
      <p className="text-sm">
        <Link href="/admin/banners" className="underline">
          ← Banners
        </Link>
      </p>
      <h1 className="text-xl font-semibold">Edit banner</h1>
      <div className="space-y-1">
        <Label>Slot</Label>
        <Input value={slot} onChange={(e) => setSlot(e.target.value)} />
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
      <div className="space-y-1">
        <Label>Starts at</Label>
        <Input value={startsAt} onChange={(e) => setStartsAt(e.target.value)} />
      </div>
      <div className="space-y-1">
        <Label>Ends at</Label>
        <Input value={endsAt} onChange={(e) => setEndsAt(e.target.value)} />
      </div>
      <div className="flex items-center gap-2">
        <Checkbox checked={active} onCheckedChange={(v) => setActive(v === true)} id="ba" />
        <Label htmlFor="ba">Active</Label>
      </div>
      <Button disabled={save.isPending} onClick={() => save.mutate()}>
        Save
      </Button>
    </div>
  );
}
