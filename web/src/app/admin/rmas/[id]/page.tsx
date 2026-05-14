"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useParams } from "next/navigation";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";
import { toast } from "sonner";

export default function AdminRmaDetailPage() {
  const { id } = useParams();
  const rid = String(id);
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["admin-rma", rid],
    queryFn: () => adminApi.rma(rid),
    enabled: !!rid,
  });
  const [note, setNote] = useState("");

  const ok = (fn: () => Promise<unknown>, msg: string) => ({
    mutationFn: fn,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["admin-rma", rid] });
      qc.invalidateQueries({ queryKey: ["admin-rmas"] });
      toast.success(msg);
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const approve = useMutation(ok(() => adminApi.rmaApprove(rid, note || undefined), "Approved"));
  const reject = useMutation(ok(() => adminApi.rmaReject(rid, note || undefined), "Rejected"));
  const receive = useMutation(ok(() => adminApi.rmaReceive(rid, note || undefined), "Received"));
  const refund = useMutation(ok(() => adminApi.rmaRefund(rid, note || undefined), "Refunded"));

  if (!q.data) return <p className="text-sm">Loading…</p>;
  const r = q.data;

  return (
    <div className="space-y-4">
      <p className="text-sm">
        <Link href="/admin/rmas" className="underline">
          ← RMAs
        </Link>
      </p>
      <h1 className="font-mono text-lg">{r.id}</h1>
      <p className="text-sm">{r.reason}</p>
      <div className="space-y-1 max-w-md">
        <Label>Note (optional)</Label>
        <Input value={note} onChange={(e) => setNote(e.target.value)} />
      </div>
      <div className="flex flex-wrap gap-2">
        <Button size="sm" variant="outline" onClick={() => approve.mutate()}>
          Approve
        </Button>
        <Button size="sm" variant="outline" onClick={() => reject.mutate()}>
          Reject
        </Button>
        <Button size="sm" variant="outline" onClick={() => receive.mutate()}>
          Receive
        </Button>
        <Button size="sm" variant="outline" onClick={() => refund.mutate()}>
          Refund
        </Button>
      </div>
    </div>
  );
}
