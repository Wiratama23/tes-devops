"use client";

import Image from "next/image";
import { useEffect, useState } from "react";
import { useDropzone } from "react-dropzone";
import { ImagePlus, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import { cn } from "@/tools/utils";

const MAX_BYTES = 10 * 1024 * 1024; // mirrors backend
const ACCEPT = {
  "image/jpeg": [".jpg", ".jpeg"],
  "image/png": [".png"],
  "image/webp": [".webp"],
  "image/gif": [".gif"],
};

interface ImageDropzoneProps {
  value: File | null;
  onChange: (file: File | null) => void;
  initialPreviewUrl?: string | null;
  disabled?: boolean;
  error?: string;
}

export function ImageDropzone({
  value,
  onChange,
  initialPreviewUrl,
  disabled,
  error,
}: ImageDropzoneProps) {
  const [previewUrl, setPreviewUrl] = useState<string | null>(
    initialPreviewUrl ?? null
  );
  const [rejection, setRejection] = useState<string | null>(null);

  useEffect(() => {
    if (value) {
      const url = URL.createObjectURL(value);
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setPreviewUrl(url);
      return () => URL.revokeObjectURL(url);
    }
    // Reset to initial preview when value is cleared
    setPreviewUrl(initialPreviewUrl ?? null);
  }, [value, initialPreviewUrl]);

  const { getRootProps, getInputProps, isDragActive, isFocused, open } =
    useDropzone({
      multiple: false,
      maxSize: MAX_BYTES,
      accept: ACCEPT,
      disabled,
      onDrop: (accepted, rejected) => {
        if (rejected.length > 0) {
          const first = rejected[0]?.errors?.[0];
          setRejection(first?.message ?? "File rejected");
          onChange(null);
          return;
        }
        setRejection(null);
        const file = accepted[0];
        if (file) {
          onChange(file);
        }
      },
      noClick: true,
    });

  return (
    <div className="space-y-2">
      <div
        {...getRootProps()}
        className={cn(
          "relative flex h-48 w-full cursor-pointer flex-col items-center justify-center gap-2 overflow-hidden rounded-md border border-dashed border-input bg-muted/30 text-sm transition-colors",
          isDragActive && "border-primary bg-primary/5",
          isFocused && "ring-1 ring-ring",
          error && "border-destructive"
        )}
        role="button"
        aria-label="Upload product image"
        onClick={() => !disabled && open()}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            if (!disabled) open();
          }
        }}
        tabIndex={0}
      >
        <input data-testid="dropzone-input" {...getInputProps()} />

        {previewUrl ? (
          <Image
            src={previewUrl}
            alt={value?.name ?? "Selected image"}
            fill
            sizes="(min-width: 768px) 50vw, 100vw"
            className="object-cover"
            unoptimized
          />
        ) : (
          <div className="flex flex-col items-center gap-2 text-muted-foreground">
            <ImagePlus className="h-6 w-6" />
            <p className="text-sm">
              {isDragActive
                ? "Drop the image here…"
                : "Drag & drop or click to upload"}
            </p>
            <p className="text-xs">JPG, PNG, WEBP, GIF · max 10 MB</p>
          </div>
        )}

        {previewUrl && value ? (
          <Button
            type="button"
            size="icon"
            variant="secondary"
            className="absolute right-2 top-2 h-7 w-7"
            onClick={(e) => {
              e.stopPropagation();
              onChange(null);
              setRejection(null);
            }}
            aria-label="Remove image"
          >
            <X className="h-4 w-4" />
          </Button>
        ) : null}
      </div>
      {(rejection || error) && (
        <p className="text-xs text-destructive">{rejection ?? error}</p>
      )}
    </div>
  );
}
