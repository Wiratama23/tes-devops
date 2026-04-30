import { z } from "zod";

// Mirrors handlers.LoginRequest / LoginResponse on the backend.
export const loginSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(6, "Password must be at least 6 characters"),
});
export type LoginInput = z.infer<typeof loginSchema>;

// Mirrors handlers.CreateProductRequest. Decimal is sent as a string to keep
// precision parity with the Go decimal.Decimal column.
export const productCreateSchema = z.object({
  product_id: z
    .string()
    .regex(/^SKU\d{2}\d+$/, "SKU must look like SKU<2-digit-type><digits>"),
  product_name: z.string().min(1, "Product name is required").max(255),
  product_quantity: z.coerce.number().int().nonnegative(),
  product_prices: z
    .string()
    .regex(/^\d+(\.\d{1,2})?$/, "Price must be a number with up to 2 decimals"),
  product_type: z.string().min(1, "Product type is required").max(8),
  image_path: z
    .string()
    .trim()
    .optional()
    .or(z.literal("")),
});
export type ProductCreateInput = z.infer<typeof productCreateSchema>;

// Mirrors handlers.UpdateProductRequest (no product_id).
export const productUpdateSchema = productCreateSchema.omit({
  product_id: true,
});
export type ProductUpdateInput = z.infer<typeof productUpdateSchema>;

// Mirrors handlers.CreateArticleRequest. uid is supplied by the auth context.
export const articleCreateSchema = z.object({
  title: z.string().min(1, "Title is required").max(255),
  article_text: z
    .string()
    .min(1, "Article body is required"),
});
export type ArticleCreateInput = z.infer<typeof articleCreateSchema>;

export const articleUpdateSchema = articleCreateSchema;
export type ArticleUpdateInput = z.infer<typeof articleUpdateSchema>;

// Contact form (client-only; no backend endpoint).
export const contactSchema = z.object({
  name: z.string().min(1, "Name is required").max(120),
  email: z.string().email("A valid email is required"),
  message: z.string().min(10, "Please write at least 10 characters").max(4000),
});
export type ContactInput = z.infer<typeof contactSchema>;
