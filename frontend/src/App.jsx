import { useState, useEffect, useRef } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { X, AlertCircle } from "lucide-react";
import Navbar from "./components/views/Navbar";
import DashboardView from "./components/views/DashboardView";
import OptimizeView from "./components/views/OptimizeView";
import AnalyticsView from "./components/views/AnalyticsView";
import HistoryView from "./components/views/HistoryView";
import SettingsDialog from "./components/views/SettingsDialog";
import { getSampleData, optimizeWithAnalytics, runHybridOptimization } from "./api/optimizer";

function App() {
  const [activeTab, setActiveTab] = useState("dashboard");
  const [orders, setOrders] = useState([]);
  const [shoppers, setShoppers] = useState([]);
  const [assignments, setAssignments] = useState([]);
  const [stats, setStats] = useState(null);
  const [analytics, setAnalytics] = useState(null);
  const [routeGeometries, setRouteGeometries] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showSettings, setShowSettings] = useState(false);
  const [useRealRoutes, setUseRealRoutes] = useState(true);
  const [algorithm, setAlgorithm] = useState("astar");
  const [apiKey, setApiKey] = useState("");
  const [apiKeyInput, setApiKeyInput] = useState("");
  const [hybridTimeline, setHybridTimeline] = useState([]);
  const [hybridRunning, setHybridRunning] = useState(false);
  const [hybridStats, setHybridStats] = useState(null);
  const hybridAbortRef = useRef(null);

  useEffect(() => {
    const savedKey = localStorage.getItem("openroute_api_key");
    if (savedKey) {
      setApiKey(savedKey);
      setApiKeyInput(savedKey);
    }
  }, []);

  const handleLoadSampleData = async () => {
    setLoading(true);
    setError(null);
    setAssignments([]);
    setStats(null);
    setAnalytics(null);
    setRouteGeometries([]);

    try {
      const data = await getSampleData();
      setOrders(data.orders);
      setShoppers(data.shoppers);
    } catch (err) {
      setError("Failed to load sample data. Make sure the backend is running on port 8080.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleOptimize = async () => {
    if (orders.length === 0 || shoppers.length === 0) {
      setError("Please load sample data first.");
      return;
    }

    if (useRealRoutes && !apiKey) {
      setError("API key required for real routes. Open settings to add your key.");
      setShowSettings(true);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const result = await optimizeWithAnalytics({ orders, shoppers }, useRealRoutes, algorithm, apiKey);
      setAssignments(result.optimization.assignments);
      setStats(result.optimization);
      setAnalytics(result.analytics);
      const geometries = result.analytics.routeGeometries || [];
      if (geometries.length > 0) {
        const firstRoutePoints = geometries[0].points?.length || 0;
        if (firstRoutePoints < 10 && useRealRoutes) {
          setError("Real routes unavailable — API key may be invalid.");
        }
      }
      setRouteGeometries(geometries);
    } catch (err) {
      setError("Failed to optimize routes. Check if the backend is running.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleHybridOptimize = async () => {
    if (orders.length === 0 || shoppers.length === 0) {
      setError("Please load sample data first.");
      return;
    }

    if (useRealRoutes && !apiKey) {
      setError("API key required for real routes. Open settings to add your key.");
      setShowSettings(true);
      return;
    }

    const controller = new AbortController();
    hybridAbortRef.current = controller;
    setHybridRunning(true);
    setLoading(true);
    setError(null);
    setHybridTimeline([]);
    setHybridStats(null);

    try {
      const availableCores =
        typeof navigator !== "undefined" && navigator.hardwareConcurrency
          ? navigator.hardwareConcurrency
          : 4;
      const result = await runHybridOptimization({
        orders,
        shoppers,
        options: {
          iterations: 320,
          workers: Math.min(4, availableCores),
          candidatePool: 240,
          randomizedListSize: 3,
          destroyRate: 0.35,
          localSearchIterations: 60,
          emitIntervalMillis: 200,
          randomSeed: Date.now(),
          useRealRoutes,
          apiKey,
        },
        onProgress: (progress) => {
          setHybridTimeline((prev) => {
            const next = [...prev, progress];
            return next.length > 200 ? next.slice(next.length - 200) : next;
          });
        },
        signal: controller.signal,
      });

      setAssignments(result.optimization.assignments || []);
      setStats(result.optimization);
      setAnalytics(result.analytics);
      setHybridStats(result.stats);
      const geometries = result.analytics?.routeGeometries || [];
      setRouteGeometries(geometries);
    } catch (err) {
      if (err.name !== "AbortError") {
        setError(err.message || "Hybrid solver failed.");
        console.error(err);
      }
    } finally {
      setHybridRunning(false);
      setLoading(false);
      hybridAbortRef.current = null;
    }
  };

  const handleHybridCancel = () => {
    if (hybridAbortRef.current) {
      hybridAbortRef.current.abort();
    }
  };

  const handleSaveApiKey = () => {
    setApiKey(apiKeyInput);
    localStorage.setItem("openroute_api_key", apiKeyInput);
    setError(null);
    setShowSettings(false);
  };

  const handleClearApiKey = () => {
    setApiKey("");
    setApiKeyInput("");
    localStorage.removeItem("openroute_api_key");
  };

  return (
    <div className="h-screen w-screen flex flex-col overflow-hidden bg-background text-foreground">
      <Navbar
        activeTab={activeTab}
        onTabChange={setActiveTab}
        algorithm={algorithm}
        onAlgorithmChange={setAlgorithm}
        useRealRoutes={useRealRoutes}
        onRealRoutesChange={setUseRealRoutes}
        onSettingsOpen={() => setShowSettings(true)}
        apiKey={apiKey}
        ordersCount={orders.length}
        shoppersCount={shoppers.length}
      />

      {/* Error Banner */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            className="border-b border-red-900/30 bg-red-950/40 overflow-hidden"
          >
            <div className="flex items-center justify-between px-6 py-2.5">
              <div className="flex items-center gap-2 text-red-400">
                <AlertCircle className="h-3.5 w-3.5 flex-shrink-0" />
                <span className="text-xs">{error}</span>
              </div>
              <button
                onClick={() => setError(null)}
                className="text-red-400/60 hover:text-red-400 transition-colors"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Tab Content */}
      <div className="flex-1 overflow-y-auto">
        {activeTab === "dashboard" && (
          <DashboardView
            orders={orders}
            shoppers={shoppers}
            onLoadSampleData={handleLoadSampleData}
            onTabChange={setActiveTab}
            loading={loading}
          />
        )}

        {activeTab === "optimize" && (
          <OptimizeView
            orders={orders}
            shoppers={shoppers}
            assignments={assignments}
            stats={stats}
            analytics={analytics}
            routeGeometries={routeGeometries}
            loading={loading}
            onOptimize={handleOptimize}
            onOptimizeHybrid={handleHybridOptimize}
            onCancelHybrid={handleHybridCancel}
            hybridRunning={hybridRunning}
            hybridTimeline={hybridTimeline}
            hybridStats={hybridStats}
            onLoadSampleData={handleLoadSampleData}
          />
        )}

        {activeTab === "analytics" && (
          <AnalyticsView
            analytics={analytics}
            assignments={assignments}
          />
        )}

        {activeTab === "history" && (
          <HistoryView />
        )}
      </div>

      {/* Settings Dialog */}
      <SettingsDialog
        open={showSettings}
        onOpenChange={setShowSettings}
        apiKey={apiKey}
        apiKeyInput={apiKeyInput}
        onApiKeyInputChange={setApiKeyInput}
        onSave={handleSaveApiKey}
        onClear={handleClearApiKey}
      />
    </div>
  );
}

export default App;
