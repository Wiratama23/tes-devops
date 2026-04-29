"use client";

import { useEffect } from "react";

import { installGlobalErrorReporter } from "@/tools/logger";

export function GlobalErrorReporter() {
  useEffect(() => {
    installGlobalErrorReporter();
  }, []);
  return null;
}
