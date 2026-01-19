import * as React from "react";
import * as PopoverPrimitive from "@radix-ui/react-popover";
import { cn } from "@/lib/utils";

const { Content, Portal, Root, Trigger } = PopoverPrimitive;

const DEFAULT_SIDE_OFFSET = 4;

const Popover = Root;
const PopoverTrigger = Trigger;

const PopoverContent = React.forwardRef<
  React.ComponentRef<typeof Content>,
  React.ComponentPropsWithoutRef<typeof Content>
>(({ className, align = "center", sideOffset = DEFAULT_SIDE_OFFSET, ...props }, ref) => (
  <Portal>
    <Content
      ref={ref}
      align={align}
      sideOffset={sideOffset}
      className={cn(
        "z-50 w-72 rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none",
        className,
      )}
      {...props}
    />
  </Portal>
));
PopoverContent.displayName = Content.displayName;

export { Popover, PopoverTrigger, PopoverContent };
