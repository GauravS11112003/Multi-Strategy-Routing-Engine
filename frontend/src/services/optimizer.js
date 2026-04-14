const BASE_URL = "http://localhost:8080/api";

export async function getSampleData() {
    try {
        const res = await fetch(`${BASE_URL}/sample-data`);
        if (!res.ok) {
            throw new Error(`HTTP error! status: ${res.status}`);
        }
        return await res.json();
    } catch (error) {
        console.error("Error fetching sample data:", error);
        throw error;
    }
}

export async function optimizeRoutes(data) {
    try {
        const res = await fetch(`${BASE_URL}/optimize`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        if (!res.ok) {
            throw new Error(`HTTP error! status: ${res.status}`);
        }
        return await res.json();
    } catch (error) {
        console.error("Error optimizing routes:", error);
        throw error;
    }
}

export async function optimizeWithAnalytics(data, useRealRoutes = false, algorithm = "nearest-neighbor", apiKey = "") {
    try {
        const res = await fetch(`${BASE_URL}/optimize-analytics`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                orders: data.orders,
                shoppers: data.shoppers,
                useRealRoutes: useRealRoutes,
                algorithm: algorithm,
                apiKey: apiKey  // Pass API key to backend
            }),
        });
        if (!res.ok) {
            throw new Error(`HTTP error! status: ${res.status}`);
        }
        return await res.json();
    } catch (error) {
        console.error("Error optimizing with analytics:", error);
        throw error;
    }
}

export async function healthCheck() {
    try {
        const res = await fetch(`${BASE_URL}/health`);
        return await res.json();
    } catch (error) {
        console.error("Error checking health:", error);
        throw error;
    }
}

export async function runHybridOptimization({ orders, shoppers, options, onProgress, signal } = {}) {
    try {
        const response = await fetch(`${BASE_URL}/optimize-hybrid-stream`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                orders,
                shoppers,
                options,
            }),
            signal,
        });

        if (!response.ok) {
            throw new Error(`Hybrid solver failed with status ${response.status}`);
        }

        if (!response.body) {
            throw new Error("Hybrid solver response body was empty");
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = "";

        while (true) {
            const { value, done } = await reader.read();
            if (done) {
                break;
            }
            buffer += decoder.decode(value, { stream: true });
            let newlineIndex;
            while ((newlineIndex = buffer.indexOf("\n")) >= 0) {
                const line = buffer.slice(0, newlineIndex).trim();
                buffer = buffer.slice(newlineIndex + 1);
                if (!line) continue;
                const event = JSON.parse(line);
                if (event.type === "progress") {
                    onProgress && onProgress(event.data);
                } else if (event.type === "completed") {
                    return event.data;
                } else if (event.type === "error") {
                    throw new Error(event.error || "Hybrid solver error");
                }
            }
        }

        // Handle any remaining buffered content.
        const trailing = buffer.trim();
        if (trailing) {
            const event = JSON.parse(trailing);
            if (event.type === "completed") {
                return event.data;
            }
            if (event.type === "error") {
                throw new Error(event.error || "Hybrid solver error");
            }
        }

        throw new Error("Hybrid solver stream ended unexpectedly without a completion event");
    } catch (error) {
        console.error("Hybrid solver request failed:", error);
        throw error;
    }
}
