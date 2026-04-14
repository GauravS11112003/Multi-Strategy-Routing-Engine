import { motion, useSpring, useTransform, animate } from 'framer-motion';
import { useEffect, useState } from 'react';
import { TrendingDown, CheckCircle2, Route } from 'lucide-react';

function AnimatedNumber({ value, decimals = 1 }) {
    const [displayValue, setDisplayValue] = useState(0);

    useEffect(() => {
        const controls = animate(displayValue, value, {
            duration: 1.5,
            ease: "easeOut",
            onUpdate: (latest) => setDisplayValue(latest),
        });

        return controls.stop;
    }, [value]);

    return <span>{displayValue.toFixed(decimals)}</span>;
}

export default function SummaryPanel({ stats, visible }) {
    if (!visible || !stats) return null;

    const improvement = stats.totalDistanceBefore > 0
        ? ((stats.totalDistanceBefore - stats.totalDistanceAfter) / stats.totalDistanceBefore) * 100
        : 0;

    const totalOrders = stats.assignments.reduce((sum, a) => sum + a.route.length, 0);

    return (
        <motion.div
            initial={{ y: 100, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ duration: 0.6, ease: "easeOut" }}
            className="absolute bottom-6 left-6 bg-white rounded-xl shadow-2xl p-6 z-[1001] min-w-[400px]"
        >
            <div className="flex items-center gap-2 mb-4">
                <CheckCircle2 className="w-6 h-6 text-shipt-green" />
                <h3 className="text-lg font-bold text-gray-800">Optimization Complete</h3>
            </div>

            <div className="grid grid-cols-2 gap-4">
                {/* Total Distance Before */}
                <motion.div
                    initial={{ scale: 0.9, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    transition={{ delay: 0.2 }}
                    className="bg-red-50 rounded-lg p-4 border border-red-100"
                >
                    <div className="text-xs text-gray-600 mb-1">Distance Before</div>
                    <div className="text-2xl font-bold text-red-600">
                        <AnimatedNumber value={stats.totalDistanceBefore} />
                        <span className="text-sm ml-1">km</span>
                    </div>
                </motion.div>

                {/* Total Distance After */}
                <motion.div
                    initial={{ scale: 0.9, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    transition={{ delay: 0.3 }}
                    className="bg-emerald-50 rounded-lg p-4 border border-emerald-100"
                >
                    <div className="text-xs text-gray-600 mb-1">Distance After</div>
                    <div className="text-2xl font-bold text-shipt-green">
                        <AnimatedNumber value={stats.totalDistanceAfter} />
                        <span className="text-sm ml-1">km</span>
                    </div>
                </motion.div>

                {/* Improvement Percentage */}
                <motion.div
                    initial={{ scale: 0.9, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    transition={{ delay: 0.4 }}
                    className="bg-blue-50 rounded-lg p-4 border border-blue-100"
                >
                    <div className="text-xs text-gray-600 mb-1 flex items-center gap-1">
                        <TrendingDown className="w-3 h-3" />
                        Improvement
                    </div>
                    <div className="text-2xl font-bold text-blue-600">
                        <AnimatedNumber value={improvement} decimals={0} />
                        <span className="text-sm ml-1">%</span>
                    </div>
                </motion.div>

                {/* Orders Optimized */}
                <motion.div
                    initial={{ scale: 0.9, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                    transition={{ delay: 0.5 }}
                    className="bg-purple-50 rounded-lg p-4 border border-purple-100"
                >
                    <div className="text-xs text-gray-600 mb-1 flex items-center gap-1">
                        <Route className="w-3 h-3" />
                        Orders
                    </div>
                    <div className="text-2xl font-bold text-purple-600">
                        {totalOrders}
                    </div>
                </motion.div>
            </div>

            {/* Route Breakdown */}
            <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.6 }}
                className="mt-4 pt-4 border-t border-gray-200"
            >
                <div className="text-xs text-gray-600 mb-2">Routes Assigned</div>
                <div className="flex flex-wrap gap-2">
                    {stats.assignments.map((assignment, idx) => (
                        <div
                            key={assignment.shopperId}
                            className="bg-gray-100 rounded px-3 py-1 text-xs font-medium text-gray-700"
                        >
                            {assignment.shopperId}: {assignment.route.length} orders
                        </div>
                    ))}
                </div>
            </motion.div>
        </motion.div>
    );
}

