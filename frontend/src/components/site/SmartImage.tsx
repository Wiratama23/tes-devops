"use client";

import Image, { type ImageProps } from "next/image";
import { useState } from "react";

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

  return (
    <Image
      {...rest}
      src={src}
      alt={alt}
      quality={85}
      onError={() => {
        if (src !== FALLBACK) setSrc(FALLBACK);
      }}
    />
  );
}
