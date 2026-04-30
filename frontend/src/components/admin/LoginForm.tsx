"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { toast } from "sonner";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ApiError } from "@/tools/api";
import { logger } from "@/tools/logger";
import { loginSchema, type LoginInput } from "@/schemas";
import { login } from "@/services/auth";

export function LoginForm() {
  const router = useRouter();
  const params = useSearchParams();
  const next = params.get("next") ?? "/admin/dashboard";

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginInput>({
    resolver: zodResolver(loginSchema),
    defaultValues: { username: "", password: "" },
  });

  const mutation = useMutation({
    mutationFn: login,
    onSuccess: (data) => {
      if (!data.user.is_admin) {
        toast.error("This account is not an admin.");
        return;
      }
      toast.success(`Welcome, ${data.user.username}.`);
      router.replace(next);
      router.refresh();
    },
    onError: (err) => {
      logger.warn("login failed", {
        kind: "auth",
        status: err instanceof ApiError ? err.status : 0,
      });
    },
  });

  const errorMessage =
    mutation.error instanceof ApiError && mutation.error.status === 401
      ? "Invalid username or password."
      : mutation.error
      ? "Something went wrong. Try again."
      : null;

  return (
    <form
      onSubmit={handleSubmit((values) => mutation.mutate(values))}
      className="space-y-4"
    >
      {errorMessage ? (
        <Alert variant="destructive">
          <AlertTitle>Sign in failed</AlertTitle>
          <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
      ) : null}

      <div className="space-y-2">
        <Label htmlFor="username">Username</Label>
        <Input
          id="username"
          autoComplete="username"
          autoFocus
          {...register("username")}
        />
        {errors.username ? (
          <p className="text-xs text-destructive">{errors.username.message}</p>
        ) : null}
      </div>

      <div className="space-y-2">
        <Label htmlFor="password">Password</Label>
        <Input
          id="password"
          type="password"
          autoComplete="current-password"
          {...register("password")}
        />
        {errors.password ? (
          <p className="text-xs text-destructive">{errors.password.message}</p>
        ) : null}
      </div>

      <div className="flex gap-2">
        <Button asChild variant="outline" className="flex-1">
          <Link href="/">Back to home</Link>
        </Button>
        <Button
          type="submit"
          className="flex-1"
          disabled={mutation.isPending}
        >
          {mutation.isPending ? "Signing in…" : "Sign in"}
        </Button>
      </div>
    </form>
  );
}
