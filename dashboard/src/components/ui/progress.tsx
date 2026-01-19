"use client"

import * as React from "react"
import * as ProgressPrimitive from "@radix-ui/react-progress"

import { cn } from "@/lib/utils"

const { Indicator, Root } = ProgressPrimitive

const FULL_PERCENT = 100
const EMPTY_VALUE = 0

const Progress = React.forwardRef<
  React.ComponentRef<typeof Root>,
  React.ComponentPropsWithoutRef<typeof Root>
>(({ className, value, ...props }, ref) => (
  <Root
    ref={ref}
    className={cn(
      "relative h-4 w-full overflow-hidden rounded-full bg-secondary",
      className
    )}
    {...props}
  >
    <Indicator
      className="h-full w-full flex-1 bg-primary transition-all"
      style={{
        transform: `translateX(-${String(FULL_PERCENT - (value ?? EMPTY_VALUE))}%)`,
      }}
    />
  </Root>
))
Progress.displayName = Root.displayName

export { Progress }
