import * as React from "react";

import { cn } from "@/lib/utils";

const Spinner = React.forwardRef<
  SVGSVGElement,
  React.SVGAttributes<SVGSVGElement>
>(({ className, ...props }, ref) => {
  const id = React.useId();

  return (
    <svg
      ref={ref}
      viewBox="0 0 24 24"
      height="1rem"
      width="1rem"
      className={cn("animate-spin", className)}
      {...props}
    >
      <defs>
        <linearGradient
          x1="28.154%"
          y1="63.74%"
          x2="74.629%"
          y2="17.783%"
          id={id}
        >
          <stop stopColor="currentColor" offset="0%" />
          <stop stopColor="#fff" stopOpacity="0" offset="100%" />
        </linearGradient>
      </defs>
      <g transform="translate(2)" fill="none">
        <circle stroke={`url(#${id})`} strokeWidth="4" cx="10" cy="12" r="10" />
        <path
          d="M10 2C4.477 2 0 6.477 0 12"
          stroke="currentColor"
          strokeWidth="4"
        />
        <rect fill="currentColor" x="8" width="4" height="4" rx="8" />
      </g>
    </svg>
  );
});

Spinner.displayName = "Spinner";

export { Spinner };
