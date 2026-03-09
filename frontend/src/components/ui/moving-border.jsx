import { useRef } from "react";
import { motion } from "framer-motion";
import { cn } from "../../lib/utils";

function MovingBorderCard({ children, className, containerClassName, borderRadius = "0.75rem", duration = 3000 }) {
  return (
    <div
      className={cn("relative p-[1px] overflow-hidden", containerClassName)}
      style={{ borderRadius }}
    >
      <div
        className="absolute inset-0"
        style={{ borderRadius }}
      >
        <motion.div
          className="absolute h-[200%] w-[200%]"
          style={{
            background: "conic-gradient(from 0deg, transparent 0 340deg, #00C389 360deg)",
            top: "-50%",
            left: "-50%",
          }}
          animate={{ rotate: 360 }}
          transition={{
            duration: duration / 1000,
            repeat: Infinity,
            ease: "linear",
          }}
        />
      </div>
      <div
        className={cn(
          "relative bg-card backdrop-blur-xl",
          className
        )}
        style={{ borderRadius: `calc(${borderRadius} - 1px)` }}
      >
        {children}
      </div>
    </div>
  );
}

export { MovingBorderCard };
