import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  Zap, Play, Square, Clock, TrendingDown, Truck, Package,
  ChevronDown, ChevronRight, Route as RouteIcon, Cpu,
  ArrowDownRight, Gauge
} from "lucide-react";
import { Button } from "../components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { Badge } from "../components/ui/badge";
import { Progress } from "../components/ui/progress";
import { MetricCard } from "../components/ui/metric-card";
import { AnimatedNumber } from "../components/ui/animated-number";
import { cn, formatNumber, formatDuration, ROUTE_COLORS } from "../lib/utils";
import MapView from "../components/features/MapView";

function OptimizePage({
  orders, shoppers, assignments, stats, analytics, routeGeometries,
  loading, onOptimize, onOptimizeHybrid, onCancelHybrid,
  hybridRunning, hybridTimeline, hybridStats,
  onLoadSampleData,
}) {
  const [expandedShopper, setExpandedShopper] = useState(null);
  const hasData = orders.length > 0;
  const hasResults = assignments.length > 0;

  const improvement =
    stats && stats.totalDistanceBefore > 0
      ? ((stats.totalDistanceBefore - stats.totalDistanceAfter) / stats.totalDistanceBefore) * 100
      : 0;

  const lastProgress = hybridTimeline.length > 0
    ? hybridTimeline[hybridTimeline.length - 1]
    : null;

  return (
    <div className="flex h-[calc(100vh-3.5rem)] overflow-hidden">
      {/* Left Panel */}
      <div className="w-[380px] flex-shrink-0 overflow-y-auto border-r border-border bg-background p-4 space-y-4">
        {/* Actions */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm flex items-center gap-2">
              <Zap className="h-4 w-4 text-primary" />
              Optimization Controls
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {!hasData && (
              <Button
                className="w-full"
                variant="outline"
                onClick={onLoadSampleData}
                disabled={loading}
              >
                <Package className="h-4 w-4" />
                Load Sample Data
              </Button>
            )}
            <Button
              className="w-full"
              variant="glow"
              onClick={onOptimize}
              disabled={loading || !hasData}
            >
              {loading && !hybridRunning ? (
                <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
              ) : (
                <Play className="h-4 w-4" />
              )}
              Run Optimization
            </Button>
            <div className="relative">
              <Button
                className="w-full"
                variant={hybridRunning ? "destructive" : "secondary"}
                onClick={hybridRunning ? onCancelHybrid : onOptimizeHybrid}
                disabled={loading && !hybridRunning || !hasData}
              >
                {hybridRunning ? (
                  <>
                    <Square className="h-4 w-4" />
                    Cancel Hybrid
                  </>
                ) : (
                  <>
                    <Cpu className="h-4 w-4" />
                    Hybrid Solver (GRASP+ALNS)
                  </>
                )}
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Hybrid Progress */}
        {(hybridRunning || hybridStats) && (
          <Card>
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Cpu className="h-4 w-4 text-purple-400" />
                  Hybrid Solver
                </CardTitle>
                {hybridRunning && (
                  <Badge variant="warning">Running</Badge>
                )}
                {!hybridRunning && hybridStats && (
                  <Badge variant="success">Complete</Badge>
                )}
              </div>
            </CardHeader>
            <CardContent className="space-y-3">
              {lastProgress && hybridRunning && (
                <>
                  <div className="flex justify-between text-xs text-muted-foreground">
                    <span>Iteration {lastProgress.iteration}</span>
                    <span>{lastProgress.exploredSolutions} explored</span>
                  </div>
                  <Progress
                    value={lastProgress.iteration ? (lastProgress.iteration / 320) * 100 : 0}
                    indicatorClassName="bg-purple-500"
                  />
                  <div className="flex justify-between text-xs">
                    <span className="text-muted-foreground">Best distance</span>
                    <span className="font-mono text-purple-400">
                      {lastProgress.bestDistance?.toFixed(2)} km
                    </span>
                  </div>
                  <div className="flex justify-between text-xs">
                    <span className="text-muted-foreground">Temperature</span>
                    <span className="font-mono">
                      {lastProgress.temperature?.toFixed(2)}
                    </span>
                  </div>
                </>
              )}
              {hybridStats && !hybridRunning && (
                <div className="space-y-2">
                  {[
                    { label: "Runtime", value: `${(hybridStats.runtime / 1e6).toFixed(0)}ms` },
                    { label: "Iterations", value: hybridStats.iterations },
                    { label: "Best at", value: `#${hybridStats.bestIteration}` },
                    { label: "Solutions explored", value: hybridStats.exploredSolutions },
                    { label: "Improvements", value: hybridStats.acceptedImprovements },
                  ].map((stat) => (
                    <div key={stat.label} className="flex justify-between text-xs">
                      <span className="text-muted-foreground">{stat.label}</span>
                      <span className="font-mono font-medium">{stat.value}</span>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        )}

        {/* Results Summary */}
        {hasResults && stats && (
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <TrendingDown className="h-4 w-4 text-emerald-400" />
                Results
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-3">
                <div className="rounded-lg bg-secondary p-3">
                  <p className="text-[10px] uppercase tracking-wider text-muted-foreground">Before</p>
                  <p className="text-lg font-bold font-mono">
                    <AnimatedNumber value={stats.totalDistanceBefore} suffix=" km" />
                  </p>
                </div>
                <div className="rounded-lg bg-primary/10 p-3">
                  <p className="text-[10px] uppercase tracking-wider text-primary">After</p>
                  <p className="text-lg font-bold font-mono text-primary">
                    <AnimatedNumber value={stats.totalDistanceAfter} suffix=" km" />
                  </p>
                </div>
              </div>

              <div className="flex items-center justify-center gap-2 py-2">
                <ArrowDownRight className="h-5 w-5 text-emerald-400" />
                <span className="text-2xl font-bold text-emerald-400">
                  <AnimatedNumber value={improvement} suffix="%" />
                </span>
                <span className="text-xs text-muted-foreground">reduction</span>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Shopper Assignments - Drill Down */}
        {hasResults && (
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-sm flex items-center gap-2">
                <Truck className="h-4 w-4 text-blue-400" />
                Shopper Routes
                <Badge variant="outline" className="ml-auto">
                  {assignments.length}
                </Badge>
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-1 p-3">
              {assignments.map((assignment, index) => {
                const isExpanded = expandedShopper === assignment.shopperId;
                const color = ROUTE_COLORS[index % ROUTE_COLORS.length];
                const shopperAnalytics = analytics?.shoppers?.find(
                  (s) => s.shopperId === assignment.shopperId
                );

                return (
                  <motion.div
                    key={assignment.shopperId}
                    initial={{ opacity: 0, x: -10 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.05 }}
                  >
                    <button
                      onClick={() =>
                        setExpandedShopper(isExpanded ? null : assignment.shopperId)
                      }
                      className="flex w-full items-center gap-3 rounded-lg p-2.5 text-left transition-colors hover:bg-accent"
                    >
                      <div
                        className="h-3 w-3 rounded-full flex-shrink-0"
                        style={{ backgroundColor: color }}
                      />
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between">
                          <span className="text-xs font-medium truncate">
                            {assignment.shopperId}
                          </span>
                          <span className="text-xs font-mono text-muted-foreground">
                            {assignment.totalDistance} km
                          </span>
                        </div>
                        <div className="flex items-center gap-2 mt-0.5">
                          <span className="text-[10px] text-muted-foreground">
                            {assignment.route.length} stops
                          </span>
                          {shopperAnalytics && (
                            <span className="text-[10px] text-muted-foreground">
                              · {shopperAnalytics.efficiency} ord/hr
                            </span>
                          )}
                        </div>
                      </div>
                      {isExpanded ? (
                        <ChevronDown className="h-3.5 w-3.5 text-muted-foreground" />
                      ) : (
                        <ChevronRight className="h-3.5 w-3.5 text-muted-foreground" />
                      )}
                    </button>

                    <AnimatePresence>
                      {isExpanded && (
                        <motion.div
                          initial={{ height: 0, opacity: 0 }}
                          animate={{ height: "auto", opacity: 1 }}
                          exit={{ height: 0, opacity: 0 }}
                          transition={{ duration: 0.2 }}
                          className="overflow-hidden"
                        >
                          <div className="ml-6 space-y-2 py-2 border-l border-border pl-3">
                            {shopperAnalytics && (
                              <div className="grid grid-cols-2 gap-2">
                                {[
                                  { label: "Distance", value: `${shopperAnalytics.totalDistance} km` },
                                  { label: "Duration", value: formatDuration(shopperAnalytics.totalDuration) },
                                  { label: "Capacity", value: `${shopperAnalytics.capacityUtilization}%` },
                                  { label: "Efficiency", value: `${shopperAnalytics.efficiency} ord/hr` },
                                ].map((stat) => (
                                  <div key={stat.label} className="rounded-md bg-secondary/50 p-2">
                                    <p className="text-[9px] uppercase tracking-wider text-muted-foreground">
                                      {stat.label}
                                    </p>
                                    <p className="text-xs font-mono font-medium">{stat.value}</p>
                                  </div>
                                ))}
                              </div>
                            )}
                            <div className="space-y-1">
                              <p className="text-[10px] uppercase tracking-wider text-muted-foreground">
                                Route sequence
                              </p>
                              {assignment.route.map((orderId, orderIndex) => (
                                <div
                                  key={orderId}
                                  className="flex items-center gap-2 text-[11px]"
                                >
                                  <span className="flex h-4 w-4 items-center justify-center rounded-full bg-secondary text-[9px] font-medium">
                                    {orderIndex + 1}
                                  </span>
                                  <span className="text-muted-foreground font-mono truncate">
                                    {orderId}
                                  </span>
                                </div>
                              ))}
                            </div>
                            {shopperAnalytics && (
                              <div className="flex items-center gap-2 text-[10px] text-muted-foreground pt-1">
                                <Clock className="h-3 w-3" />
                                {shopperAnalytics.estimatedStartTime} — {shopperAnalytics.estimatedEndTime}
                              </div>
                            )}
                          </div>
                        </motion.div>
                      )}
                    </AnimatePresence>
                  </motion.div>
                );
              })}
            </CardContent>
          </Card>
        )}
      </div>

      {/* Right Panel - Map */}
      <div className="flex-1 relative">
        <MapView
          orders={orders}
          shoppers={shoppers}
          assignments={assignments}
          routeGeometries={routeGeometries}
        />

        {/* Floating stats overlay */}
        {hasResults && (
          <div className="absolute top-4 right-4 z-[400] flex flex-col gap-2">
            <div className="glass rounded-lg px-3 py-2">
              <div className="flex items-center gap-2">
                <Gauge className="h-3.5 w-3.5 text-primary" />
                <span className="text-xs font-medium">
                  <AnimatedNumber value={improvement} suffix="%" decimals={1} />
                  {" "}saved
                </span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default OptimizePage;
