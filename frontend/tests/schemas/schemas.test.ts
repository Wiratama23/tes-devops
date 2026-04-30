import { describe, expect, it } from "vitest";

import {
  articleCreateSchema,
  contactSchema,
  loginSchema,
  productCreateSchema,
} from "@/schemas";

describe("loginSchema", () => {
  it("accepts a valid login", () => {
    const r = loginSchema.safeParse({ username: "alice", password: "secret123" });
    expect(r.success).toBe(true);
  });

  it("rejects short passwords", () => {
    const r = loginSchema.safeParse({ username: "alice", password: "x" });
    expect(r.success).toBe(false);
  });
});

describe("productCreateSchema", () => {
  it("accepts a valid product", () => {
    const r = productCreateSchema.safeParse({
      product_id: "SKU10001",
      product_name: "Coffee",
      product_quantity: 10,
      product_prices: "29.99",
      product_type: "10",
      image_path: "assets/x.jpg",
    });
    expect(r.success).toBe(true);
  });

  it("rejects bad SKU format", () => {
    const r = productCreateSchema.safeParse({
      product_id: "ABC",
      product_name: "Coffee",
      product_quantity: 10,
      product_prices: "29.99",
      product_type: "10",
      image_path: "assets/x.jpg",
    });
    expect(r.success).toBe(false);
  });

  it("rejects malformed prices", () => {
    const r = productCreateSchema.safeParse({
      product_id: "SKU10001",
      product_name: "Coffee",
      product_quantity: 10,
      product_prices: "29.999",
      product_type: "10",
      image_path: "assets/x.jpg",
    });
    expect(r.success).toBe(false);
  });
});

describe("articleCreateSchema", () => {
  it("requires title and body", () => {
    const r = articleCreateSchema.safeParse({ title: "", article_text: "" });
    expect(r.success).toBe(false);
  });
});

describe("contactSchema", () => {
  it("requires a 10+ char message", () => {
    const r = contactSchema.safeParse({
      name: "A",
      email: "a@b.io",
      message: "short",
    });
    expect(r.success).toBe(false);
  });
});
