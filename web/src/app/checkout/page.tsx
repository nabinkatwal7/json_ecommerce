"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { customerApi } from "@/lib/api";
import type { Address } from "@/lib/types";
import { useAuth } from "@/contexts/auth-context";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useState } from "react";
import { toast } from "sonner";
import { fmtMoney } from "@/lib/format";

const emptyAddr: Address = {
  fullName: "",
  phone: "",
  country: "US",
  state: "",
  city: "",
  postalCode: "",
  addressLine: "",
  isDefault: false,
};

export default function CheckoutPage() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const qc = useQueryClient();
  const cart = useQuery({
    queryKey: ["cart"],
    queryFn: () => customerApi.cart(),
    enabled: !!user,
  });
  const addresses = useQuery({
    queryKey: ["addresses"],
    queryFn: () => customerApi.addresses(),
    enabled: !!user,
  });

  const [addr, setAddr] = useState<Address>(emptyAddr);
  const [discountCode, setDiscountCode] = useState("");
  const [carrier, setCarrier] = useState("flat");
  const [rates, setRates] = useState<{ code: string; label: string; amount: number }[]>([]);

  const checkout = useMutation({
    mutationFn: () =>
      customerApi.checkout({
        shippingAddress: addr,
        discountCode: discountCode || undefined,
        shippingCarrier: carrier,
      }),
    onSuccess: (order) => {
      qc.invalidateQueries({ queryKey: ["cart"] });
      qc.invalidateQueries({ queryKey: ["orders"] });
      toast.success("Order created");
      router.push(`/orders/${order.id}`);
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (loading) return null;
  if (!user) {
    return (
      <p className="p-8 text-center text-sm">
        <Link href="/login" className="underline">
          Log in
        </Link>{" "}
        to checkout.
      </p>
    );
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6 px-4 py-8">
      <h1 className="text-2xl font-semibold">Checkout</h1>
      <Card className="border">
        <CardHeader>
          <CardTitle className="text-base">Coupon</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-2">
          <Input
            placeholder="MATCHDAY20"
            value={discountCode}
            onChange={(e) => setDiscountCode(e.target.value)}
          />
          <Button
            type="button"
            variant="outline"
            onClick={async () => {
              try {
                const r = await customerApi.couponValidate(discountCode);
                toast.message(r.valid ? "Coupon OK" : "Coupon not applied", {
                  description: `${r.message} · discount ${fmtMoney(r.discountAmount)} on subtotal ${fmtMoney(r.subtotal)}`,
                });
              } catch (e: unknown) {
                toast.error(e instanceof Error ? e.message : "Failed");
              }
            }}
          >
            Validate coupon
          </Button>
        </CardContent>
      </Card>
      <Card className="border">
        <CardHeader>
          <CardTitle className="text-base">Shipping quote</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={async () => {
              try {
                const r = await customerApi.shippingQuote({
                  ...addr,
                  country: addr.country || "US",
                });
                setRates(r.rates);
                toast.success("Rates loaded");
              } catch (e: unknown) {
                toast.error(e instanceof Error ? e.message : "Quote failed");
              }
            }}
          >
            Get rates (uses address country)
          </Button>
          {rates.length > 0 && (
            <ul className="text-sm">
              {rates.map((x) => (
                <li key={x.code}>
                  {x.label}: {fmtMoney(x.amount)} ({x.code})
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
      <Card className="border">
        <CardHeader>
          <CardTitle className="text-base">Saved addresses</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          {(addresses.data ?? []).map((a) => (
            <Button
              key={a.id}
              type="button"
              variant="outline"
              size="sm"
              className="mr-2"
              onClick={() => setAddr({ ...a, id: a.id })}
            >
              {a.fullName} — {a.city}
            </Button>
          ))}
        </CardContent>
      </Card>
      <div className="grid gap-3 sm:grid-cols-2">
        {(
          [
            ["fullName", "Full name"],
            ["phone", "Phone"],
            ["country", "Country"],
            ["state", "State"],
            ["city", "City"],
            ["postalCode", "Postal code"],
            ["addressLine", "Address line"],
          ] as const
        ).map(([k, lab]) => (
          <div key={k} className="space-y-1 sm:col-span-2">
            <Label htmlFor={k}>{lab}</Label>
            <Input
              id={k}
              value={String(addr[k as keyof Address] ?? "")}
              onChange={(e) => setAddr({ ...addr, [k]: e.target.value })}
            />
          </div>
        ))}
      </div>
      <div className="space-y-2">
        <Label>Shipping carrier</Label>
        <Select value={carrier} onValueChange={setCarrier}>
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="flat">Flat</SelectItem>
            <SelectItem value="fedex_ground">FedEx Ground</SelectItem>
            <SelectItem value="ups_ground">UPS Ground</SelectItem>
            <SelectItem value="dhl_express">DHL Express</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <Button
        disabled={checkout.isPending || !(cart.data?.items ?? []).length}
        onClick={() => checkout.mutate()}
      >
        Place order
      </Button>
    </div>
  );
}
