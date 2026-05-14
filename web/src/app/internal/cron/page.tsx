"use client";

import { useState } from "react";
import { publicApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";

export default function InternalCronPage() {
  const [secret, setSecret] = useState("");
  const [busy, setBusy] = useState(false);

  async function run() {
    setBusy(true);
    try {
      const r = await publicApi.abandonedCron(secret);
      toast.success(`Sent: ${r.sent}`);
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : "Failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="mx-auto max-w-md px-4 py-12">
      <Card className="border">
        <CardHeader>
          <CardTitle>Abandoned cart cron</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 text-sm">
          <p className="text-muted-foreground">
            Calls <code className="rounded bg-muted px-1">POST /internal/cron/abandoned-carts</code>{" "}
            with header <code className="rounded bg-muted px-1">X-Cron-Secret</code>.
          </p>
          <div className="space-y-1">
            <Label htmlFor="s">CRON_SECRET</Label>
            <Input id="s" value={secret} onChange={(e) => setSecret(e.target.value)} />
          </div>
          <Button type="button" disabled={busy} onClick={run}>
            Run job
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
