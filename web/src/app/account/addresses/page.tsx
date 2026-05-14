"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { customerApi } from "@/lib/api";
import type { Address } from "@/lib/types";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useState } from "react";
import { toast } from "sonner";

const blank: Address = {
  fullName: "",
  phone: "",
  country: "US",
  state: "",
  city: "",
  postalCode: "",
  addressLine: "",
  isDefault: false,
};

export default function AddressesPage() {
  const { user, loading } = useAuth();
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["addresses"],
    queryFn: () => customerApi.addresses(),
    enabled: !!user,
  });
  const [form, setForm] = useState<Address>(blank);
  const [editing, setEditing] = useState<string | null>(null);

  const add = useMutation({
    mutationFn: () => customerApi.postAddress(form),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["addresses"] });
      toast.success("Added");
      setForm(blank);
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const del = useMutation({
    mutationFn: (id: string) => customerApi.deleteAddress(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["addresses"] }),
  });

  const put = useMutation({
    mutationFn: () => customerApi.putAddress(editing!, { ...form, id: editing! }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["addresses"] });
      toast.success("Updated");
      setEditing(null);
      setForm(blank);
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

  return (
    <div className="mx-auto max-w-xl space-y-6 px-4 py-8">
      <h1 className="text-2xl font-semibold">Addresses</h1>
      <ul className="space-y-2">
        {(q.data ?? []).map((a) => (
          <li key={a.id} className="flex flex-wrap items-start justify-between gap-2 border p-3 text-sm">
            <div>
              <p className="font-medium">{a.fullName}</p>
              <p>{a.addressLine}</p>
              <p>
                {a.city}, {a.state} {a.postalCode} · {a.country}
              </p>
            </div>
            <div className="flex gap-2">
              <Button
                size="sm"
                variant="outline"
                onClick={() => {
                  setEditing(a.id!);
                  setForm(a);
                }}
              >
                Edit
              </Button>
              <Button size="sm" variant="ghost" onClick={() => a.id && del.mutate(a.id)}>
                Delete
              </Button>
            </div>
          </li>
        ))}
      </ul>
      <Card className="border">
        <CardHeader>
          <CardTitle className="text-base">{editing ? "Edit address" : "New address"}</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-2 sm:grid-cols-2">
          {(
            [
              ["fullName", "Full name"],
              ["phone", "Phone"],
              ["country", "Country"],
              ["state", "State"],
              ["city", "City"],
              ["postalCode", "Postal"],
              ["addressLine", "Street", true],
            ] as const
          ).map(([k, lab, full]) => (
            <div key={k} className={full ? "sm:col-span-2" : ""}>
              <Label>{lab}</Label>
              <Input
                value={String(form[k] ?? "")}
                onChange={(e) => setForm({ ...form, [k]: e.target.value })}
              />
            </div>
          ))}
          <div className="sm:col-span-2">
            <Button
              type="button"
              onClick={() => (editing ? put.mutate() : add.mutate())}
              disabled={add.isPending || put.isPending}
            >
              {editing ? "Save changes" : "Add address"}
            </Button>
            {editing && (
              <Button
                type="button"
                variant="ghost"
                className="ml-2"
                onClick={() => {
                  setEditing(null);
                  setForm(blank);
                }}
              >
                Cancel
              </Button>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
