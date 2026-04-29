// API types mirroring backend models in backend/internal/models/*.

export type UUID = string;

export interface User {
  uid: UUID;
  username: string;
  email: string;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthUser {
  uid: UUID;
  username: string;
  email: string;
  is_admin: boolean;
}

export interface LoginResponse {
  token: string;
  user: AuthUser;
  expires: number;
}

export interface Article {
  articles_id: number;
  uid: UUID;
  title: string;
  article_text: string;
  date_created: string;
  updated_at: string;
}

export interface PaginatedArticles {
  data: Article[];
  total_count: number;
  limit: number;
  offset: number;
}

export interface Product {
  product_id: string;
  product_name: string;
  product_quantity: number;
  product_prices: string;
  product_type: string;
  created_at: string;
  created_by: UUID;
  image_path: string;
}

export interface PaginatedProducts {
  data: Product[];
  limit: number;
  offset: number;
}

export interface UploadResponse {
  path: string;
  url: string;
  filename: string;
  size: number;
}

export type LogLevel = "info" | "warn" | "error";

export interface ClientLogEntry {
  level: LogLevel;
  message: string;
  stack?: string;
  url?: string;
  user_agent?: string;
  meta?: Record<string, unknown>;
}
