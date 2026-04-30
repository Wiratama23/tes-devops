import { NextResponse, type NextRequest } from "next/server";

import {
  AUTH_COOKIE_NAME,
  decodeAuthToken,
  isAdminFromClaims,
} from "./tools/auth";

// Matches /admin and everything below it. /admin/login is allowed through
// inside the proxy body so the login page can render while logged out.
export const config = {
  matcher: ["/admin/:path*"],
};

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;

  if (pathname === "/admin/login" || pathname.startsWith("/admin/login/")) {
    return NextResponse.next();
  }

  const token = request.cookies.get(AUTH_COOKIE_NAME)?.value;
  const claims = decodeAuthToken(token);

  if (!isAdminFromClaims(claims)) {
    const loginUrl = request.nextUrl.clone();
    loginUrl.pathname = "/admin/login";
    loginUrl.searchParams.set("next", pathname);
    return NextResponse.redirect(loginUrl);
  }

  return NextResponse.next();
}
