"use client";

import Image, { type ImageProps } from "next/image";
import { useState, useCallback } from "react";

import { resolveImageUrl } from "@/services/uploads";

const FALLBACK =
  process.env.NEXT_PUBLIC_DEFAULT_IMAGE ?? "/assets/default_image.jpg";

interface SmartImageProps extends Omit<ImageProps, "src"> {
  // Backend `image_path` value (e.g. "assets/coffee.jpg") OR a fully
  // qualified URL OR a relative `/foo.jpg` path. Empty/null falls back to the
  // bundled default.
  imagePath: string | null | undefined;
}

export function SmartImage({ imagePath, alt, ...rest }: SmartImageProps) {
  const initial = resolveImageUrl(imagePath);
  const [src, setSrc] = useState(initial);
  const [failureCount, setFailureCount] = useState(0);

  const handleImageError = useCallback(() => {
    // Only try fallback once to avoid infinite loops
    if (src !== FALLBACK && failureCount === 0) {
      setSrc(FALLBACK);
      setFailureCount(1);
    }
  }, [src, failureCount]);

  return (
    <Image
      {...rest}
      src={src}
      alt={alt}
      quality={85}
      onError={handleImageError}
    />
  );
}
