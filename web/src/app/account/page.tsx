"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useState } from "react";
import { toast } from "sonner";

export default function AccountPage() {
  const { user, loading } = useAuth();
  const qc = useQueryClient();
  const me = useQuery({
    queryKey: ["me"],
    queryFn: () => customerApi.me(),
    enabled: !!user,
  });
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");

  const patch = useMutation({
    mutationFn: () =>
      customerApi.patchMe({
        ...(name.trim() ? { name: name.trim() } : {}),
        ...(email.trim() ? { email: email.trim() } : {}),
      }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["me"] });
      toast.success("Profile updated");
      setName("");
      setEmail("");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (loading) return null;
  if (!user) {
    return (
      <p className="p-8 text-center text-sm">
        <Link href="/login" className="underline">
          Log in
        </Link>
      </p>
    );
  }

  const u = me.data;

  return (
    <div className="mx-auto max-w-lg space-y-6 px-4 py-8">
      <h1 className="text-2xl font-semibold">Account</h1>
      {u && (
        <Card className="border">
          <CardHeader>
            <CardTitle className="text-base">Current</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1 text-sm">
            <p>{u.name}</p>
            <p className="text-muted-foreground">{u.email}</p>
            <p className="text-xs">Role: {u.role}</p>
          </CardContent>
        </Card>
      )}
      <Card className="border">
        <CardHeader>
          <CardTitle className="text-base">Update profile</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-1">
            <Label htmlFor="n">New name (optional)</Label>
            <Input id="n" value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-1">
            <Label htmlFor="e">New email (optional)</Label>
            <Input id="e" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
          </div>
          <Button
            disabled={patch.isPending || (!name.trim() && !email.trim())}
            onClick={() => patch.mutate()}
          >
            Save
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
