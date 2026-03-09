import { useEffect } from "react";
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import { ROUTE_COLORS } from "../lib/utils";

const BIRMINGHAM = [33.5186, -86.8104];

const shopperIcon = new L.DivIcon({
  html: `<div style="
    width: 28px; height: 28px;
    background: #00C389;
    border: 2px solid #09090b;
    border-radius: 50%;
    box-shadow: 0 0 12px rgba(0,195,137,0.4);
    display: flex; align-items: center; justify-content: center;
  ">
    <svg width="14" height="14" viewBox="0 0 24 24" fill="white" xmlns="http://www.w3.org/2000/svg">
      <path d="M18.92 6.01C18.72 5.42 18.16 5 17.5 5h-11c-.66 0-1.21.42-1.42 1.01L3 12v8c0 .55.45 1 1 1h1c.55 0 1-.45 1-1v-1h12v1c0 .55.45 1 1 1h1c.55 0 1-.45 1-1v-8l-2.08-5.99zM6.5 16c-.83 0-1.5-.67-1.5-1.5S5.67 13 6.5 13s1.5.67 1.5 1.5S7.33 16 6.5 16zm11 0c-.83 0-1.5-.67-1.5-1.5s.67-1.5 1.5-1.5 1.5.67 1.5 1.5-.67 1.5-1.5 1.5zM5 11l1.5-4.5h11L19 11H5z"/>
    </svg>
  </div>`,
  className: "",
  iconSize: [28, 28],
  iconAnchor: [14, 14],
  popupAnchor: [0, -14],
});

const orderIcon = new L.DivIcon({
  html: `<div style="
    width: 22px; height: 22px;
    background: #f59e0b;
    border: 2px solid #09090b;
    border-radius: 50%;
    box-shadow: 0 0 8px rgba(245,158,11,0.3);
    display: flex; align-items: center; justify-content: center;
  ">
    <svg width="12" height="12" viewBox="0 0 24 24" fill="white" xmlns="http://www.w3.org/2000/svg">
      <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V5h14v14z"/>
      <path d="M12 6l-4 4h3v4h2v-4h3z"/>
    </svg>
  </div>`,
  className: "",
  iconSize: [22, 22],
  iconAnchor: [11, 11],
  popupAnchor: [0, -11],
});

function MapBounds({ orders, shoppers }) {
  const map = useMap();

  useEffect(() => {
    const points = [];
    shoppers.forEach((s) => points.push([s.lat, s.lng]));
    orders.forEach((o) => points.push([o.lat, o.lng]));

    if (points.length > 0) {
      map.fitBounds(points, { padding: [40, 40], maxZoom: 14 });
    }
  }, [orders, shoppers, map]);

  return null;
}

function MapView({ orders = [], shoppers = [], assignments = [], routeGeometries = [] }) {
  const shopperOrderMap = {};
  assignments.forEach((a) => {
    shopperOrderMap[a.shopperId] = a;
  });

  const orderMap = {};
  orders.forEach((o) => {
    orderMap[o.id] = o;
  });

  return (
    <MapContainer
      center={BIRMINGHAM}
      zoom={12}
      className="h-full w-full"
      zoomControl={true}
    >
      <TileLayer
        attribution='&copy; <a href="https://carto.com">CARTO</a>'
        url="https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"
      />

      <MapBounds orders={orders} shoppers={shoppers} />

      {/* Route polylines */}
      {routeGeometries.length > 0
        ? routeGeometries.map((geo, index) => (
            <Polyline
              key={`geo-${geo.shopperId}`}
              positions={geo.points}
              pathOptions={{
                color: ROUTE_COLORS[index % ROUTE_COLORS.length],
                weight: 3,
                opacity: 0.7,
                dashArray: null,
              }}
            />
          ))
        : assignments.map((assignment, index) => {
            const shopper = shoppers.find((s) => s.id === assignment.shopperId);
            if (!shopper) return null;

            const waypoints = [
              [shopper.lat, shopper.lng],
              ...assignment.route
                .map((id) => orderMap[id])
                .filter(Boolean)
                .map((o) => [o.lat, o.lng]),
            ];

            return (
              <Polyline
                key={`route-${assignment.shopperId}`}
                positions={waypoints}
                pathOptions={{
                  color: ROUTE_COLORS[index % ROUTE_COLORS.length],
                  weight: 2.5,
                  opacity: 0.5,
                  dashArray: "8 6",
                }}
              />
            );
          })}

      {/* Shopper markers */}
      {shoppers.map((shopper) => {
        const assignment = shopperOrderMap[shopper.id];
        return (
          <Marker
            key={`shopper-${shopper.id}`}
            position={[shopper.lat, shopper.lng]}
            icon={shopperIcon}
          >
            <Popup>
              <div style={{ minWidth: "160px" }}>
                <p style={{ fontWeight: 600, fontSize: "13px", color: "#00C389", marginBottom: "6px" }}>
                  Shopper
                </p>
                <p style={{ fontSize: "11px", color: "#a1a1aa", fontFamily: "monospace" }}>
                  {shopper.id}
                </p>
                <p style={{ fontSize: "11px", color: "#a1a1aa", marginTop: "4px" }}>
                  Capacity: {shopper.capacity}
                </p>
                {assignment && (
                  <p style={{ fontSize: "11px", color: "#00C389", marginTop: "4px" }}>
                    {assignment.route.length} stops · {assignment.totalDistance} km
                  </p>
                )}
              </div>
            </Popup>
          </Marker>
        );
      })}

      {/* Order markers */}
      {orders.map((order) => (
        <Marker
          key={`order-${order.id}`}
          position={[order.lat, order.lng]}
          icon={orderIcon}
        >
          <Popup>
            <div style={{ minWidth: "140px" }}>
              <p style={{ fontWeight: 600, fontSize: "13px", color: "#f59e0b", marginBottom: "6px" }}>
                Order
              </p>
              <p style={{ fontSize: "11px", color: "#a1a1aa", fontFamily: "monospace" }}>
                {order.id}
              </p>
              <p style={{ fontSize: "11px", color: "#a1a1aa", marginTop: "4px" }}>
                {order.itemCount} items · {order.deliveryWindow}
              </p>
            </div>
          </Popup>
        </Marker>
      ))}
    </MapContainer>
  );
}

export default MapView;
