"use client";

import { useEffect } from "react";

import {
  installGlobalErrorReporter,
  uninstallGlobalErrorReporter,
} from "@/tools/logger";

export function GlobalErrorReporter() {
  useEffect(() => {
    installGlobalErrorReporter();
    return () => uninstallGlobalErrorReporter();
  }, []);
  return null;
}
