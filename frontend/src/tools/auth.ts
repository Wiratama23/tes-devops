import { decodeJwt, type JWTPayload } from "jose";

export const AUTH_COOKIE_NAME = "auth_token";

export interface AuthClaims extends JWTPayload {
  uid?: string;
  username?: string;
  is_admin?: boolean;
}

// Decodes the JWT without verifying its signature. Verification stays on the
// Go backend; the frontend only inspects claims to make routing decisions
// (e.g. show admin nav, redirect if non-admin), so a forged token here cannot
// access protected data — it would still be rejected by the API.
export function decodeAuthToken(token: string | undefined | null): AuthClaims | null {
  if (!token) return null;
  try {
    return decodeJwt(token) as AuthClaims;
  } catch {
    return null;
  }
}

export function isTokenExpired(claims: AuthClaims | null): boolean {
  if (!claims?.exp) return true;
  const nowSeconds = Math.floor(Date.now() / 1000);
  return nowSeconds >= claims.exp;
}

export function isAdminFromClaims(claims: AuthClaims | null): boolean {
  return Boolean(claims?.is_admin) && !isTokenExpired(claims);
}
