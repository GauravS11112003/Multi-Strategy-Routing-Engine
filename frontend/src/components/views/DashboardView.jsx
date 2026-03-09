import { motion } from "framer-motion";
import {
  Package, Truck, Zap, Map, ArrowRight, Database, Cpu, Globe,
  Activity, Layers
} from "lucide-react";
import { Button } from "../ui/button";
import { MovingBorderCard } from "../ui/moving-border";
import { SpotlightCard } from "../ui/spotlight-card";
import { Badge } from "../ui/badge";

function DashboardView({ orders, shoppers, onLoadSampleData, onTabChange, loading }) {
  const hasData = orders.length > 0;

  return (
    <div className="mx-auto max-w-6xl space-y-8 p-6">
      {/* Hero Section */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="relative overflow-hidden rounded-2xl border border-border bg-card"
      >
        <div className="absolute inset-0 dot-pattern opacity-40" />
        <div className="absolute -top-24 -right-24 h-64 w-64 rounded-full bg-primary/5 blur-3xl" />
        <div className="absolute -bottom-24 -left-24 h-64 w-64 rounded-full bg-blue-500/5 blur-3xl" />

        <div className="relative p-10">
          <div className="flex items-start justify-between">
            <div className="space-y-4">
              <div className="flex items-center gap-2">
                <Badge variant="default">Multi-Strategy Engine</Badge>
                <Badge variant="outline">v2.0</Badge>
              </div>
              <h1 className="text-4xl font-bold tracking-tight">
                <span className="gradient-text">Intelligent</span> Route
                <br />
                Optimization
              </h1>
              <p className="max-w-md text-base text-muted-foreground leading-relaxed">
                Three optimization algorithms. Real-time analytics.
                Built for Shipt-scale delivery logistics.
              </p>
              <div className="flex items-center gap-3 pt-2">
                <Button
                  variant="glow"
                  size="lg"
                  onClick={hasData ? () => onTabChange("optimize") : onLoadSampleData}
                  disabled={loading}
                >
                  {loading ? (
                    <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                  ) : hasData ? (
                    <>
                      <Zap className="h-4 w-4" />
                      Start Optimizing
                    </>
                  ) : (
                    <>
                      <Database className="h-4 w-4" />
                      Load Sample Data
                    </>
                  )}
                </Button>
                {hasData && (
                  <Button variant="outline" size="lg" onClick={() => onTabChange("analytics")}>
                    <Activity className="h-4 w-4" />
                    View Analytics
                  </Button>
                )}
              </div>
            </div>

            {/* Architecture diagram */}
            <div className="hidden lg:block">
              <div className="grid grid-cols-2 gap-3">
                {[
                  { icon: Cpu, label: "3 Algorithms", desc: "Greedy, A*, GRASP+ALNS" },
                  { icon: Database, label: "PostgreSQL", desc: "Persistent storage" },
                  { icon: Globe, label: "Redis Cache", desc: "Distance matrix" },
                  { icon: Layers, label: "Kafka Events", desc: "Async processing" },
                ].map((item, i) => (
                  <motion.div
                    key={item.label}
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ delay: 0.1 * i + 0.3 }}
                    className="flex items-center gap-3 rounded-lg border border-border/50 bg-background/50 p-3"
                  >
                    <div className="rounded-md bg-primary/10 p-2">
                      <item.icon className="h-4 w-4 text-primary" />
                    </div>
                    <div>
                      <p className="text-xs font-medium">{item.label}</p>
                      <p className="text-[10px] text-muted-foreground">{item.desc}</p>
                    </div>
                  </motion.div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </motion.div>

      {/* Quick Stats */}
      {hasData && (
        <div className="grid gap-4 md:grid-cols-3">
          <MovingBorderCard className="p-5" duration={4000}>
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-primary/10 p-2.5">
                <Package className="h-5 w-5 text-primary" />
              </div>
              <div>
                <p className="text-2xl font-bold">{orders.length}</p>
                <p className="text-xs text-muted-foreground">Active Orders</p>
              </div>
            </div>
          </MovingBorderCard>

          <MovingBorderCard className="p-5" duration={5000}>
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-blue-500/10 p-2.5">
                <Truck className="h-5 w-5 text-blue-400" />
              </div>
              <div>
                <p className="text-2xl font-bold">{shoppers.length}</p>
                <p className="text-xs text-muted-foreground">Available Shoppers</p>
              </div>
            </div>
          </MovingBorderCard>

          <MovingBorderCard className="p-5" duration={6000}>
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-purple-500/10 p-2.5">
                <Map className="h-5 w-5 text-purple-400" />
              </div>
              <div>
                <p className="text-2xl font-bold">Birmingham</p>
                <p className="text-xs text-muted-foreground">Service Area</p>
              </div>
            </div>
          </MovingBorderCard>
        </div>
      )}

      {/* Feature Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        {[
          {
            title: "Nearest Neighbor",
            desc: "Greedy assignment with fast O(n²) routing. Best for small batches.",
            icon: Package,
            color: "text-emerald-400",
            bg: "bg-emerald-500/10",
          },
          {
            title: "A* Search",
            desc: "Optimal TSP routing with MST heuristic. Up to 8 stops per shopper.",
            icon: Zap,
            color: "text-blue-400",
            bg: "bg-blue-500/10",
          },
          {
            title: "GRASP + ALNS",
            desc: "Metaheuristic solver with parallel workers, simulated annealing, and real-time streaming.",
            icon: Cpu,
            color: "text-purple-400",
            bg: "bg-purple-500/10",
          },
        ].map((feature, i) => (
          <SpotlightCard
            key={feature.title}
            className="cursor-pointer"
            spotlightColor={
              i === 0
                ? "rgba(0,195,137,0.06)"
                : i === 1
                ? "rgba(59,130,246,0.06)"
                : "rgba(168,85,247,0.06)"
            }
          >
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 * i + 0.2 }}
              onClick={() => onTabChange("optimize")}
            >
              <div className={`rounded-lg ${feature.bg} p-2.5 w-fit`}>
                <feature.icon className={`h-5 w-5 ${feature.color}`} />
              </div>
              <h3 className="mt-3 text-sm font-semibold">{feature.title}</h3>
              <p className="mt-1 text-xs text-muted-foreground leading-relaxed">
                {feature.desc}
              </p>
              <div className="mt-3 flex items-center gap-1 text-xs font-medium text-primary">
                Try it <ArrowRight className="h-3 w-3" />
              </div>
            </motion.div>
          </SpotlightCard>
        ))}
      </div>
    </div>
  );
}

export default DashboardView;
