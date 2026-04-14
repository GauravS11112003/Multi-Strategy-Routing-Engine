import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  BarChart3, Truck, Package, TrendingUp, Fuel, Leaf,
  Clock, Gauge, Users, ArrowUpRight, ChevronRight, X,
  Target, Layers
} from "lucide-react";
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip,
  ResponsiveContainer, PieChart, Pie, Cell, Area, AreaChart
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { Badge } from "../components/ui/badge";
import { MetricCard } from "../components/ui/metric-card";
import { Progress } from "../components/ui/progress";
import { AnimatedNumber } from "../components/ui/animated-number";
import { SpotlightCard } from "../components/ui/spotlight-card";
import { cn, formatDuration, ROUTE_COLORS } from "../lib/utils";

const CHART_COLORS = ["#00C389", "#3b82f6", "#a855f7", "#f43f5e", "#f59e0b"];

function CustomTooltip({ active, payload, label }) {
  if (!active || !payload?.length) return null;
  return (
    <div className="rounded-lg border border-border bg-card px-3 py-2 shadow-xl">
      <p className="text-xs font-medium">{label}</p>
      {payload.map((entry, i) => (
        <p key={i} className="text-xs text-muted-foreground">
          <span style={{ color: entry.color }}>{entry.name}:</span>{" "}
          <span className="font-mono font-medium text-foreground">{entry.value}</span>
        </p>
      ))}
    </div>
  );
}

function AnalyticsPage({ analytics, assignments }) {
  const [drillDown, setDrillDown] = useState(null);

  if (!analytics) {
    return (
      <div className="flex h-[calc(100vh-3.5rem)] items-center justify-center">
        <div className="text-center space-y-3">
          <BarChart3 className="h-12 w-12 text-muted-foreground/30 mx-auto" />
          <p className="text-sm text-muted-foreground">
            Run an optimization to see analytics
          </p>
        </div>
      </div>
    );
  }

  const { system, shoppers, orders } = analytics;

  const shopperChartData = (shoppers || []).map((s, i) => ({
    name: s.shopperId?.slice(-4) || `S${i}`,
    orders: s.ordersAssigned,
    distance: s.totalDistance,
    efficiency: s.efficiency,
    capacity: s.capacityUtilization,
  }));

  const timeWindowData = orders?.timeWindowBreakdown
    ? Object.entries(orders.timeWindowBreakdown).map(([key, val]) => ({
        name: key,
        count: val,
      }))
    : [];

  const efficiencyData = (shoppers || []).map((s, i) => ({
    name: s.shopperId?.slice(-4) || `S${i}`,
    efficiency: s.efficiency,
    capacity: s.capacityUtilization,
  }));

  return (
    <div className="relative">
      <div className={cn(
        "mx-auto max-w-7xl space-y-6 p-6 transition-all duration-300",
        drillDown && "blur-sm pointer-events-none"
      )}>
        {/* Top Metrics */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <MetricCard
            title="Optimization Score"
            value={system.optimizationScore}
            suffix="/100"
            decimals={0}
            icon={Target}
            description="Overall efficiency rating"
            delay={0}
          />
          <MetricCard
            title="Average Efficiency"
            value={system.averageEfficiency}
            suffix=" ord/hr"
            icon={Gauge}
            description="Orders per hour per shopper"
            delay={0.05}
          />
          <MetricCard
            title="Fuel Cost"
            value={system.estimatedFuelCost}
            prefix="$"
            icon={Fuel}
            description="Estimated total fuel cost"
            delay={0.1}
          />
          <MetricCard
            title="CO₂ Saved"
            value={system.co2Saved}
            suffix=" kg"
            icon={Leaf}
            accentColor="text-emerald-400"
            description="vs unoptimized routing"
            delay={0.15}
          />
        </div>

        {/* System Overview */}
        <div className="grid gap-4 md:grid-cols-3">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-2">
                <Users className="h-3.5 w-3.5" />
                Shoppers
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold">
                  <AnimatedNumber value={system.activeShoppers} decimals={0} />
                </span>
                <span className="text-sm text-muted-foreground">
                  / {system.totalShoppers} active
                </span>
              </div>
              <Progress
                value={(system.activeShoppers / system.totalShoppers) * 100}
                className="mt-3"
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-2">
                <Package className="h-3.5 w-3.5" />
                Orders
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold">
                  <AnimatedNumber value={system.assignedOrders} decimals={0} />
                </span>
                <span className="text-sm text-muted-foreground">
                  / {system.totalOrders} assigned
                </span>
              </div>
              <Progress
                value={(system.assignedOrders / system.totalOrders) * 100}
                className="mt-3"
                indicatorClassName="bg-blue-500"
              />
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-xs uppercase tracking-wider text-muted-foreground flex items-center gap-2">
                <Clock className="h-3.5 w-3.5" />
                Total Time
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold">
                  {formatDuration(system.totalDuration)}
                </span>
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                {system.totalDistance} km total distance
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Charts Row */}
        <div className="grid gap-4 md:grid-cols-2">
          {/* Orders per Shopper */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm">Orders per Shopper</CardTitle>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={220}>
                <BarChart data={shopperChartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#27272a" />
                  <XAxis dataKey="name" tick={{ fill: "#71717a", fontSize: 11 }} />
                  <YAxis tick={{ fill: "#71717a", fontSize: 11 }} />
                  <Tooltip content={<CustomTooltip />} />
                  <Bar dataKey="orders" fill="#00C389" radius={[4, 4, 0, 0]} name="Orders" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>

          {/* Delivery Windows */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm">Delivery Windows</CardTitle>
            </CardHeader>
            <CardContent>
              {timeWindowData.length > 0 ? (
                <ResponsiveContainer width="100%" height={220}>
                  <PieChart>
                    <Pie
                      data={timeWindowData}
                      cx="50%"
                      cy="50%"
                      innerRadius={55}
                      outerRadius={85}
                      paddingAngle={3}
                      dataKey="count"
                    >
                      {timeWindowData.map((_, index) => (
                        <Cell
                          key={index}
                          fill={CHART_COLORS[index % CHART_COLORS.length]}
                        />
                      ))}
                    </Pie>
                    <Tooltip content={<CustomTooltip />} />
                  </PieChart>
                </ResponsiveContainer>
              ) : (
                <div className="flex h-[220px] items-center justify-center text-sm text-muted-foreground">
                  No window data
                </div>
              )}
              {timeWindowData.length > 0 && (
                <div className="flex flex-wrap gap-2 justify-center">
                  {timeWindowData.map((entry, i) => (
                    <div key={entry.name} className="flex items-center gap-1.5">
                      <div
                        className="h-2 w-2 rounded-full"
                        style={{ backgroundColor: CHART_COLORS[i % CHART_COLORS.length] }}
                      />
                      <span className="text-[10px] text-muted-foreground">{entry.name}: {entry.count}</span>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Efficiency Chart */}
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm">Shopper Efficiency & Capacity</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={200}>
              <AreaChart data={efficiencyData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#27272a" />
                <XAxis dataKey="name" tick={{ fill: "#71717a", fontSize: 11 }} />
                <YAxis tick={{ fill: "#71717a", fontSize: 11 }} />
                <Tooltip content={<CustomTooltip />} />
                <Area
                  type="monotone"
                  dataKey="efficiency"
                  stroke="#00C389"
                  fill="#00C389"
                  fillOpacity={0.1}
                  name="Efficiency (ord/hr)"
                />
                <Area
                  type="monotone"
                  dataKey="capacity"
                  stroke="#3b82f6"
                  fill="#3b82f6"
                  fillOpacity={0.1}
                  name="Capacity (%)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Shopper Cards - Drill Down */}
        <div>
          <h3 className="text-sm font-semibold mb-3 flex items-center gap-2">
            <Truck className="h-4 w-4 text-primary" />
            Shopper Performance
            <span className="text-xs text-muted-foreground font-normal">
              Click for details
            </span>
          </h3>
          <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
            {(shoppers || []).map((shopper, i) => (
              <SpotlightCard
                key={shopper.shopperId}
                className="cursor-pointer"
                spotlightColor="rgba(0,195,137,0.05)"
              >
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: i * 0.03 }}
                  onClick={() => setDrillDown(shopper)}
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center gap-2">
                      <div
                        className="h-2.5 w-2.5 rounded-full"
                        style={{ backgroundColor: ROUTE_COLORS[i % ROUTE_COLORS.length] }}
                      />
                      <span className="text-xs font-medium font-mono truncate max-w-[120px]">
                        {shopper.shopperId}
                      </span>
                    </div>
                    <ChevronRight className="h-3.5 w-3.5 text-muted-foreground" />
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <p className="text-[10px] text-muted-foreground uppercase tracking-wider">Orders</p>
                      <p className="text-lg font-bold">{shopper.ordersAssigned}</p>
                    </div>
                    <div>
                      <p className="text-[10px] text-muted-foreground uppercase tracking-wider">Distance</p>
                      <p className="text-lg font-bold">{shopper.totalDistance} km</p>
                    </div>
                  </div>

                  <div className="mt-3">
                    <div className="flex items-center justify-between text-[10px] mb-1">
                      <span className="text-muted-foreground">Capacity</span>
                      <span className="font-medium">{shopper.capacityUtilization}%</span>
                    </div>
                    <Progress value={shopper.capacityUtilization} />
                  </div>
                </motion.div>
              </SpotlightCard>
            ))}
          </div>
        </div>

        {/* Order Stats */}
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm flex items-center gap-2">
              <Package className="h-4 w-4 text-blue-400" />
              Order Insights
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[
                { label: "Total Orders", value: orders.totalOrders },
                { label: "Total Items", value: orders.totalItems },
                { label: "Avg Items/Order", value: orders.averageItemCount?.toFixed(1) },
                { label: "Avg Distance", value: `${orders.averageDistance} km` },
                { label: "Order Density", value: `${orders.orderDensity}/km²` },
                { label: "Unassigned", value: orders.unassignedOrders },
              ].map((stat) => (
                <div key={stat.label} className="rounded-lg bg-secondary/50 p-3">
                  <p className="text-[10px] uppercase tracking-wider text-muted-foreground">{stat.label}</p>
                  <p className="text-sm font-bold font-mono mt-1">{stat.value}</p>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Drill-Down Panel */}
      <AnimatePresence>
        {drillDown && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-50 flex items-center justify-center p-4"
            onClick={() => setDrillDown(null)}
          >
            <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" />
            <motion.div
              initial={{ scale: 0.95, opacity: 0, y: 20 }}
              animate={{ scale: 1, opacity: 1, y: 0 }}
              exit={{ scale: 0.95, opacity: 0, y: 20 }}
              onClick={(e) => e.stopPropagation()}
              className="relative z-10 w-full max-w-lg rounded-xl border border-border bg-card p-6 shadow-2xl"
            >
              <div className="flex items-center justify-between mb-5">
                <div>
                  <h3 className="text-base font-semibold">Shopper Details</h3>
                  <p className="text-xs text-muted-foreground font-mono mt-0.5">
                    {drillDown.shopperId}
                  </p>
                </div>
                <button
                  onClick={() => setDrillDown(null)}
                  className="rounded-md p-1 hover:bg-accent transition-colors"
                >
                  <X className="h-4 w-4 text-muted-foreground" />
                </button>
              </div>

              <div className="grid grid-cols-2 gap-3 mb-4">
                {[
                  { label: "Orders Assigned", value: drillDown.ordersAssigned, icon: Package },
                  { label: "Total Distance", value: `${drillDown.totalDistance} km`, icon: TrendingUp },
                  { label: "Duration", value: formatDuration(drillDown.totalDuration), icon: Clock },
                  { label: "Efficiency", value: `${drillDown.efficiency} ord/hr`, icon: Gauge },
                ].map((stat) => (
                  <div key={stat.label} className="rounded-lg border border-border bg-secondary/30 p-3">
                    <div className="flex items-center gap-2 mb-1">
                      <stat.icon className="h-3 w-3 text-muted-foreground" />
                      <p className="text-[10px] uppercase tracking-wider text-muted-foreground">
                        {stat.label}
                      </p>
                    </div>
                    <p className="text-base font-bold font-mono">{stat.value}</p>
                  </div>
                ))}
              </div>

              <div className="space-y-3">
                <div>
                  <div className="flex items-center justify-between text-xs mb-1.5">
                    <span className="text-muted-foreground">Capacity Utilization</span>
                    <span className="font-medium">{drillDown.capacityUtilization}%</span>
                  </div>
                  <Progress value={drillDown.capacityUtilization} />
                </div>

                <div className="rounded-lg bg-secondary/30 p-3 border border-border">
                  <p className="text-[10px] uppercase tracking-wider text-muted-foreground mb-2">
                    Time Window
                  </p>
                  <div className="flex items-center gap-3">
                    <Clock className="h-4 w-4 text-primary" />
                    <div>
                      <p className="text-sm font-medium">
                        {drillDown.estimatedStartTime} — {drillDown.estimatedEndTime}
                      </p>
                      <p className="text-[10px] text-muted-foreground">
                        Avg {drillDown.averageOrderDistance} km between stops
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

export default AnalyticsPage;
