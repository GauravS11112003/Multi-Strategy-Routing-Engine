import { motion } from "framer-motion";
import { cn } from "../../lib/utils";
import { AnimatedNumber } from "./animated-number";

function MetricCard({
  title,
  value,
  suffix = "",
  prefix = "",
  decimals = 1,
  description,
  icon: Icon,
  trend,
  onClick,
  className,
  accentColor = "text-primary",
  delay = 0,
}) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay }}
      onClick={onClick}
      className={cn(
        "group relative overflow-hidden rounded-xl border border-border bg-card p-5 transition-all duration-300",
        onClick && "cursor-pointer hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5",
        className
      )}
    >
      <div className="flex items-start justify-between">
        <div className="space-y-1">
          <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
            {title}
          </p>
          <div className="flex items-baseline gap-1">
            <span className="text-2xl font-bold tracking-tight">
              <AnimatedNumber
                value={value}
                prefix={prefix}
                suffix={suffix}
                decimals={decimals}
              />
            </span>
          </div>
          {description && (
            <p className="text-xs text-muted-foreground">{description}</p>
          )}
          {trend !== undefined && (
            <div
              className={cn(
                "inline-flex items-center gap-1 text-xs font-medium",
                trend >= 0 ? "text-emerald-400" : "text-red-400"
              )}
            >
              {trend >= 0 ? "↑" : "↓"} {Math.abs(trend).toFixed(1)}%
            </div>
          )}
        </div>
        {Icon && (
          <div className={cn("rounded-lg bg-primary/10 p-2.5", accentColor)}>
            <Icon className="h-5 w-5" />
          </div>
        )}
      </div>
      {onClick && (
        <div className="absolute bottom-0 left-0 h-[2px] w-0 bg-primary transition-all duration-300 group-hover:w-full" />
      )}
    </motion.div>
  );
}

export { MetricCard };
