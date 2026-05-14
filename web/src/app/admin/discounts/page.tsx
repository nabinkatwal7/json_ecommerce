"use client";

import { useMutation } from "@tanstack/react-query";
import { adminApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useState } from "react";
import { toast } from "sonner";

export default function AdminDiscountsPage() {
  const [code, setCode] = useState("");
  const [type, setType] = useState("percent");
  const [value, setValue] = useState("10");
  const [minAmt, setMinAmt] = useState("0");
  const [expires, setExpires] = useState("");
  const [buy, setBuy] = useState("1");
  const [get, setGet] = useState("1");
  const [productId, setProductId] = useState("");
  const [categoryId, setCategoryId] = useState("");

  const create = useMutation({
    mutationFn: () =>
      adminApi.postDiscount({
        code,
        type,
        value: Number(value),
        minimumAmount: Number(minAmt),
        isActive: true,
        expiresAt: expires || undefined,
        buyQty: Number(buy),
        getQty: Number(get),
        productId: productId || undefined,
        categoryId: categoryId || undefined,
      }),
    onSuccess: () => toast.success("Discount created"),
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <div className="max-w-md space-y-4">
      <h1 className="text-2xl font-semibold">Discounts</h1>
      <div className="space-y-2 rounded-md border p-4">
        <div className="space-y-1">
          <Label>Code</Label>
          <Input value={code} onChange={(e) => setCode(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Type</Label>
          <Select value={type} onValueChange={setType}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="percent">percent</SelectItem>
              <SelectItem value="fixed">fixed</SelectItem>
              <SelectItem value="bogo">bogo</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-1">
          <Label>Value (percent or fixed amount)</Label>
          <Input value={value} onChange={(e) => setValue(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Minimum amount</Label>
          <Input value={minAmt} onChange={(e) => setMinAmt(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Expires at (RFC3339, optional)</Label>
          <Input value={expires} onChange={(e) => setExpires(e.target.value)} />
        </div>
        <div className="grid grid-cols-2 gap-2">
          <div className="space-y-1">
            <Label>BOGO buy qty</Label>
            <Input value={buy} onChange={(e) => setBuy(e.target.value)} />
          </div>
          <div className="space-y-1">
            <Label>BOGO get qty</Label>
            <Input value={get} onChange={(e) => setGet(e.target.value)} />
          </div>
        </div>
        <div className="space-y-1">
          <Label>Product ID (BOGO scope, optional)</Label>
          <Input value={productId} onChange={(e) => setProductId(e.target.value)} />
        </div>
        <div className="space-y-1">
          <Label>Category ID (BOGO scope, optional)</Label>
          <Input value={categoryId} onChange={(e) => setCategoryId(e.target.value)} />
        </div>
        <Button disabled={create.isPending} onClick={() => create.mutate()}>
          Create discount
        </Button>
      </div>
    </div>
  );
}
