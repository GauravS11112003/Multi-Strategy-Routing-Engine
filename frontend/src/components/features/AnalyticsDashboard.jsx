import { motion } from 'framer-motion';
import {
    TrendingUp, Users, Package, Clock, DollarSign,
    Leaf, Target, Activity, BarChart3, Truck
} from 'lucide-react';
import { useState } from 'react';

function StatCard({ icon: Icon, label, value, unit, color = "blue", delay = 0 }) {
    const colorClasses = {
        blue: 'bg-blue-50 border-blue-100',
        green: 'bg-green-50 border-green-100',
        purple: 'bg-purple-50 border-purple-100',
        orange: 'bg-orange-50 border-orange-100',
        indigo: 'bg-indigo-50 border-indigo-100',
        pink: 'bg-pink-50 border-pink-100',
        yellow: 'bg-yellow-50 border-yellow-100',
        red: 'bg-red-50 border-red-100',
    };

    const iconColorClasses = {
        blue: 'text-blue-500',
        green: 'text-green-500',
        purple: 'text-purple-500',
        orange: 'text-orange-500',
        indigo: 'text-indigo-500',
        pink: 'text-pink-500',
        yellow: 'text-yellow-500',
        red: 'text-red-500',
    };

    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay, duration: 0.4 }}
            className={`${colorClasses[color]} rounded-lg p-4 border`}
        >
            <div className="flex items-start justify-between">
                <div className="flex-1">
                    <div className="text-xs text-gray-600 mb-1">{label}</div>
                    <div className="text-2xl font-bold text-gray-800">
                        {value}
                        {unit && <span className="text-sm ml-1 text-gray-600">{unit}</span>}
                    </div>
                </div>
                <Icon className={`w-8 h-8 ${iconColorClasses[color]} opacity-60`} />
            </div>
        </motion.div>
    );
}

function ShopperCard({ shopper, index }) {
    return (
        <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: index * 0.05 }}
            className="bg-white rounded-lg p-4 border border-gray-200 hover:border-shipt-green transition-colors"
        >
            <div className="flex items-start justify-between mb-3">
                <div>
                    <div className="font-semibold text-gray-800">{shopper.shopperId}</div>
                    <div className="text-xs text-gray-500">{shopper.ordersAssigned} orders assigned</div>
                </div>
                <div className={`px-2 py-1 rounded text-xs font-semibold ${shopper.capacityUtilization >= 90 ? 'bg-red-100 text-red-700' :
                    shopper.capacityUtilization >= 70 ? 'bg-yellow-100 text-yellow-700' :
                        'bg-green-100 text-green-700'
                    }`}>
                    {shopper.capacityUtilization.toFixed(0)}% capacity
                </div>
            </div>

            <div className="grid grid-cols-2 gap-3 text-sm">
                <div>
                    <div className="text-gray-500 text-xs">Distance</div>
                    <div className="font-semibold text-gray-800">{shopper.totalDistance.toFixed(1)} km</div>
                </div>
                <div>
                    <div className="text-gray-500 text-xs">Duration</div>
                    <div className="font-semibold text-gray-800">{shopper.totalDuration.toFixed(0)} min</div>
                </div>
                <div>
                    <div className="text-gray-500 text-xs">Efficiency</div>
                    <div className="font-semibold text-gray-800">{shopper.efficiency.toFixed(1)} ord/hr</div>
                </div>
                <div>
                    <div className="text-gray-500 text-xs">Avg Distance</div>
                    <div className="font-semibold text-gray-800">{shopper.averageOrderDistance.toFixed(1)} km</div>
                </div>
            </div>

            <div className="mt-3 pt-3 border-t border-gray-100 text-xs text-gray-600">
                <div className="flex items-center gap-2">
                    <Clock className="w-3 h-3" />
                    <span>{shopper.estimatedStartTime} - {shopper.estimatedEndTime}</span>
                </div>
            </div>
        </motion.div>
    );
}

export default function AnalyticsDashboard({ analytics, visible }) {
    const [activeTab, setActiveTab] = useState('overview');

    if (!visible || !analytics) return null;

    const { system, shoppers, orders } = analytics;

    return (
        <motion.div
            initial={{ x: 400, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            exit={{ x: 400, opacity: 0 }}
            transition={{ duration: 0.5, ease: "easeOut" }}
            className="absolute right-0 top-0 bottom-0 w-[450px] bg-gray-50 shadow-2xl z-30 overflow-y-auto"
        >
            {/* Header */}
            <div className="bg-gradient-to-r from-shipt-green to-green-600 text-white p-6 sticky top-0 z-10">
                <div className="flex items-center gap-2 mb-2">
                    <BarChart3 className="w-6 h-6" />
                    <h2 className="text-xl font-bold">Analytics Dashboard</h2>
                </div>
                <p className="text-sm text-green-100">Real-time performance insights</p>
            </div>

            {/* Tabs */}
            <div className="flex bg-white border-b border-gray-200 sticky top-[88px] z-10">
                <button
                    onClick={() => setActiveTab('overview')}
                    className={`flex-1 px-4 py-3 text-sm font-semibold transition-colors ${activeTab === 'overview'
                        ? 'text-shipt-green border-b-2 border-shipt-green'
                        : 'text-gray-600 hover:text-gray-800'
                        }`}
                >
                    Overview
                </button>
                <button
                    onClick={() => setActiveTab('shoppers')}
                    className={`flex-1 px-4 py-3 text-sm font-semibold transition-colors ${activeTab === 'shoppers'
                        ? 'text-shipt-green border-b-2 border-shipt-green'
                        : 'text-gray-600 hover:text-gray-800'
                        }`}
                >
                    Shoppers
                </button>
                <button
                    onClick={() => setActiveTab('orders')}
                    className={`flex-1 px-4 py-3 text-sm font-semibold transition-colors ${activeTab === 'orders'
                        ? 'text-shipt-green border-b-2 border-shipt-green'
                        : 'text-gray-600 hover:text-gray-800'
                        }`}
                >
                    Orders
                </button>
            </div>

            {/* Content */}
            <div className="p-6">
                {activeTab === 'overview' && (
                    <div className="space-y-4">
                        {/* System Stats */}
                        <div>
                            <h3 className="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
                                <Target className="w-4 h-4" />
                                System Performance
                            </h3>
                            <div className="grid grid-cols-2 gap-3">
                                <StatCard
                                    icon={Activity}
                                    label="Optimization Score"
                                    value={system.optimizationScore.toFixed(0)}
                                    unit="/100"
                                    color="purple"
                                    delay={0.1}
                                />
                                <StatCard
                                    icon={TrendingUp}
                                    label="Avg Efficiency"
                                    value={system.averageEfficiency.toFixed(1)}
                                    unit="ord/hr"
                                    color="blue"
                                    delay={0.15}
                                />
                            </div>
                        </div>

                        {/* Resource Utilization */}
                        <div>
                            <h3 className="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
                                <Users className="w-4 h-4" />
                                Resources
                            </h3>
                            <div className="grid grid-cols-2 gap-3">
                                <StatCard
                                    icon={Truck}
                                    label="Active Shoppers"
                                    value={system.activeShoppers}
                                    unit={`of ${system.totalShoppers}`}
                                    color="green"
                                    delay={0.2}
                                />
                                <StatCard
                                    icon={Package}
                                    label="Assigned Orders"
                                    value={system.assignedOrders}
                                    unit={`of ${system.totalOrders}`}
                                    color="orange"
                                    delay={0.25}
                                />
                            </div>
                        </div>

                        {/* Distance & Time */}
                        <div>
                            <h3 className="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
                                <Clock className="w-4 h-4" />
                                Logistics
                            </h3>
                            <div className="grid grid-cols-2 gap-3">
                                <StatCard
                                    icon={Activity}
                                    label="Total Distance"
                                    value={system.totalDistance.toFixed(1)}
                                    unit="km"
                                    color="indigo"
                                    delay={0.3}
                                />
                                <StatCard
                                    icon={Clock}
                                    label="Total Duration"
                                    value={system.totalDuration.toFixed(0)}
                                    unit="min"
                                    color="pink"
                                    delay={0.35}
                                />
                            </div>
                        </div>

                        {/* Cost & Environmental */}
                        <div>
                            <h3 className="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
                                <DollarSign className="w-4 h-4" />
                                Impact
                            </h3>
                            <div className="grid grid-cols-2 gap-3">
                                <StatCard
                                    icon={DollarSign}
                                    label="Estimated Fuel Cost"
                                    value={system.estimatedFuelCost.toFixed(2)}
                                    unit="USD"
                                    color="yellow"
                                    delay={0.4}
                                />
                                <StatCard
                                    icon={Leaf}
                                    label="CO₂ Saved"
                                    value={system.co2Saved.toFixed(1)}
                                    unit="kg"
                                    color="green"
                                    delay={0.45}
                                />
                            </div>
                        </div>
                    </div>
                )}

                {activeTab === 'shoppers' && (
                    <div className="space-y-3">
                        <div className="text-sm text-gray-600 mb-4">
                            Detailed performance metrics for each shopper
                        </div>
                        {shoppers.map((shopper, idx) => (
                            <ShopperCard key={shopper.shopperId} shopper={shopper} index={idx} />
                        ))}
                    </div>
                )}

                {activeTab === 'orders' && (
                    <div className="space-y-4">
                        <div className="grid grid-cols-2 gap-3">
                            <StatCard
                                icon={Package}
                                label="Total Orders"
                                value={orders.totalOrders}
                                color="orange"
                                delay={0.1}
                            />
                            <StatCard
                                icon={Activity}
                                label="Avg Item Count"
                                value={orders.averageItemCount.toFixed(1)}
                                color="blue"
                                delay={0.15}
                            />
                            <StatCard
                                icon={Package}
                                label="Total Items"
                                value={orders.totalItems}
                                color="purple"
                                delay={0.2}
                            />
                            <StatCard
                                icon={Activity}
                                label="Order Density"
                                value={orders.orderDensity.toFixed(1)}
                                unit="per km²"
                                color="indigo"
                                delay={0.25}
                            />
                            <StatCard
                                icon={TrendingUp}
                                label="Avg Distance"
                                value={orders.averageDistance.toFixed(1)}
                                unit="km"
                                color="green"
                                delay={0.3}
                            />
                            <StatCard
                                icon={Package}
                                label="Unassigned"
                                value={orders.unassignedOrders}
                                color={orders.unassignedOrders > 0 ? "red" : "green"}
                                delay={0.35}
                            />
                        </div>

                        {/* Time Window Breakdown */}
                        {orders.timeWindowBreakdown && Object.keys(orders.timeWindowBreakdown).length > 0 && (
                            <div className="mt-6">
                                <h3 className="text-sm font-semibold text-gray-700 mb-3">Delivery Windows</h3>
                                <div className="bg-white rounded-lg border border-gray-200 p-4">
                                    {Object.entries(orders.timeWindowBreakdown).map(([window, count], idx) => (
                                        <div
                                            key={window}
                                            className="flex justify-between items-center py-2 border-b border-gray-100 last:border-0"
                                        >
                                            <span className="text-sm text-gray-700">{window}</span>
                                            <span className="font-semibold text-gray-800">{count} orders</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </motion.div>
    );
}

