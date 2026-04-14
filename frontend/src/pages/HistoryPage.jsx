import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { History, Clock, TrendingDown, Zap, RefreshCw } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { Badge } from "../components/ui/badge";
import { Button } from "../components/ui/button";

const BASE_URL = "http://localhost:8080/api";

function HistoryPage() {
  const [runs, setRuns] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchHistory = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${BASE_URL}/optimizations?limit=20`);
      if (res.status === 503) {
        setRuns([]);
        return;
      }
      if (!res.ok) throw new Error("Failed to fetch");
      const data = await res.json();
      setRuns(data.optimizations || []);
    } catch (err) {
      setError("Could not load history. Database may not be connected.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHistory();
  }, []);

  return (
    <div className="mx-auto max-w-4xl space-y-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <History className="h-5 w-5 text-primary" />
            Optimization History
          </h2>
          <p className="text-xs text-muted-foreground mt-1">
            Past optimization runs persisted to PostgreSQL
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={fetchHistory} disabled={loading}>
          <RefreshCw className={`h-3.5 w-3.5 ${loading ? "animate-spin" : ""}`} />
          Refresh
        </Button>
      </div>

      {error && (
        <div className="rounded-lg border border-border bg-secondary/50 p-4 text-center">
          <p className="text-sm text-muted-foreground">{error}</p>
          <p className="text-xs text-muted-foreground mt-1">
            Connect PostgreSQL via docker-compose to enable persistence.
          </p>
        </div>
      )}

      {!error && runs.length === 0 && !loading && (
        <div className="rounded-lg border border-border bg-card p-8 text-center">
          <History className="h-10 w-10 text-muted-foreground/30 mx-auto mb-3" />
          <p className="text-sm text-muted-foreground">No optimization runs yet</p>
          <p className="text-xs text-muted-foreground mt-1">
            Run an optimization — results will be saved here automatically.
          </p>
        </div>
      )}

      <div className="space-y-3">
        {runs.map((run, index) => (
          <motion.div
            key={run.id}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.03 }}
          >
            <Card className="hover:border-border/80 transition-colors">
              <CardContent className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="rounded-lg bg-primary/10 p-2">
                      <Zap className="h-4 w-4 text-primary" />
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <p className="text-sm font-medium">{run.algorithm}</p>
                        <Badge variant="outline" className="text-[10px]">
                          {run.totalOrders} orders · {run.totalShoppers} shoppers
                        </Badge>
                      </div>
                      <div className="flex items-center gap-3 mt-1 text-xs text-muted-foreground">
                        <span className="flex items-center gap-1">
                          <Clock className="h-3 w-3" />
                          {run.durationMs}ms
                        </span>
                        <span className="flex items-center gap-1">
                          <TrendingDown className="h-3 w-3 text-emerald-400" />
                          {run.improvementPct?.toFixed(1)}% improvement
                        </span>
                      </div>
                    </div>
                  </div>

                  <div className="text-right">
                    <p className="text-xs font-mono text-muted-foreground">
                      {run.distanceBefore?.toFixed(1)} → {run.distanceAfter?.toFixed(1)} km
                    </p>
                    <p className="text-[10px] text-muted-foreground mt-1">
                      {run.createdAt ? new Date(run.createdAt).toLocaleString() : ""}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </motion.div>
        ))}
      </div>
    </div>
  );
}

export default HistoryPage;
