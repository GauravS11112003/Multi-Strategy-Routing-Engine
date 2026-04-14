import { motion } from "framer-motion";
import { Settings, Package, Truck, BarChart3, History, Zap, Route as RouteIcon, ChevronDown } from "lucide-react";
import { cn } from "../../lib/utils";
import { Badge } from "../ui/badge";

function Navbar({
  activeTab,
  onTabChange,
  algorithm,
  onAlgorithmChange,
  useRealRoutes,
  onRealRoutesChange,
  onSettingsOpen,
  apiKey,
  ordersCount,
  shoppersCount,
}) {
  const tabs = [
    { id: "dashboard", label: "Dashboard", icon: Package },
    { id: "optimize", label: "Optimize", icon: Zap },
    { id: "analytics", label: "Analytics", icon: BarChart3 },
    { id: "history", label: "History", icon: History },
  ];

  return (
    <motion.header
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.3 }}
      className="sticky top-0 z-50 border-b border-border bg-background/80 backdrop-blur-xl"
    >
      <div className="flex h-14 items-center justify-between px-6">
        {/* Left: Logo + Nav */}
        <div className="flex items-center gap-8">
          <div className="flex items-center gap-2.5">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
              <Truck className="h-4 w-4 text-white" />
            </div>
            <span className="text-sm font-semibold tracking-tight">
              Multi-Strat<span className="text-primary"> Routing</span>
            </span>
          </div>

          <nav className="flex items-center">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;
              return (
                <button
                  key={tab.id}
                  onClick={() => onTabChange(tab.id)}
                  className={cn(
                    "relative flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium transition-colors",
                    isActive
                      ? "text-foreground"
                      : "text-muted-foreground hover:text-foreground"
                  )}
                >
                  <Icon className="h-3.5 w-3.5" />
                  {tab.label}
                  {isActive && (
                    <motion.div
                      layoutId="activeTab"
                      className="absolute inset-x-0 -bottom-[15px] h-[2px] bg-primary"
                      transition={{ type: "spring", stiffness: 500, damping: 30 }}
                    />
                  )}
                </button>
              );
            })}
          </nav>
        </div>

        {/* Right: Controls */}
        <div className="flex items-center gap-3">
          {(ordersCount > 0 || shoppersCount > 0) && (
            <div className="flex items-center gap-2">
              <Badge variant="default">{ordersCount} orders</Badge>
              <Badge variant="success">{shoppersCount} shoppers</Badge>
            </div>
          )}

          {/* Algorithm Toggle */}
          <div className="flex items-center rounded-lg border border-border bg-secondary p-0.5">
            <button
              onClick={() => onAlgorithmChange("nearest-neighbor")}
              className={cn(
                "flex items-center gap-1 rounded-md px-2.5 py-1 text-xs font-medium transition-all",
                algorithm === "nearest-neighbor"
                  ? "bg-primary text-white shadow-sm"
                  : "text-muted-foreground hover:text-foreground"
              )}
            >
              <RouteIcon className="h-3 w-3" />
              Greedy
            </button>
            <button
              onClick={() => onAlgorithmChange("astar")}
              className={cn(
                "flex items-center gap-1 rounded-md px-2.5 py-1 text-xs font-medium transition-all",
                algorithm === "astar"
                  ? "bg-primary text-white shadow-sm"
                  : "text-muted-foreground hover:text-foreground"
              )}
            >
              <Zap className="h-3 w-3" />
              A*
            </button>
          </div>

          {/* Real Routes Toggle */}
          <button
            onClick={() => onRealRoutesChange(!useRealRoutes)}
            className={cn(
              "flex items-center gap-1.5 rounded-lg border px-3 py-1.5 text-xs font-medium transition-all",
              useRealRoutes
                ? "border-primary/40 bg-primary/10 text-primary"
                : "border-border text-muted-foreground hover:text-foreground"
            )}
          >
            <RouteIcon className="h-3 w-3" />
            Roads
            <div
              className={cn(
                "h-2 w-2 rounded-full transition-colors",
                useRealRoutes ? "bg-primary" : "bg-muted-foreground/30"
              )}
            />
          </button>

          {/* Settings */}
          <button
            onClick={onSettingsOpen}
            className={cn(
              "flex h-8 w-8 items-center justify-center rounded-lg border transition-colors",
              apiKey
                ? "border-primary/30 text-primary hover:bg-primary/10"
                : "border-border text-muted-foreground hover:text-foreground hover:bg-accent"
            )}
          >
            <Settings className="h-4 w-4" />
          </button>
        </div>
      </div>
    </motion.header>
  );
}

export default Navbar;
