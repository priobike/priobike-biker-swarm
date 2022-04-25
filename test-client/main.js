import { Vector, distanceP2P } from "./utils.js";

const AUTH_ENDPOINT = "http://priobike.vkw.tu-dresden.de/production/session-wrapper/authentication";
const WS_ENDPOINT = "ws://priobike.vkw.tu-dresden.de/production/session-wrapper/websocket/sessions/";
const ROUTE_ENDPOINT = "http://priobike.vkw.tu-dresden.de/production/session-wrapper/getroute";

const CLIENT_ID = "JS_TEST_CLIENT";

const INTERPOLATE_STEP_LENGTH = 4; // m
const POSITION_UPDATE_INTERVAL = 1000; // ms

// ======= AUTH REQUEST =======
console.log("======= AUTH REQUEST =======");

const authResponse = await fetch(AUTH_ENDPOINT, {
  method: "POST",
  body: JSON.stringify({
    clientId: CLIENT_ID,
  }),
  headers: {
    "Content-type": "application/json; charset=UTF-8",
  },
}).then((response) => response.json());

console.log(authResponse);

const sessionId = authResponse.sessionId;

// ======= ROUTING REQUEST =======
console.log("======= ROUTING REQUEST =======");

const routingResponse = await fetch(ROUTE_ENDPOINT, {
  method: "POST",
  body: JSON.stringify({
    sessionId,
    waypoints: [
      {
        'lon': 9.977496, 'lat': 53.56415
      },
      {
        'lon': 9.990059, 'lat': 53.560791
      },
    ]
  }),
  headers: {
    "Content-type": "application/json; charset=UTF-8",
  },
})
.catch((err) => {
  console.log(err);
})
.then((response) => {
  console.log(response);
  return response.json();
})

console.log(routingResponse);

console.log("======= INTERPOLATE ROUTE =======");

let ridingPoints = [];
for (let i = 0; i < routingResponse.route.length; i++) {
  const point = routingResponse.route[i];
  const nextPoint = routingResponse.route[i + 1];

  if (!nextPoint) {
    break;
  }

  const pointV = new Vector(point.lat, point.lon);
  const nextPointV = new Vector(nextPoint.lat, nextPoint.lon);
  const distance = distanceP2P(
    point.lat,
    point.lon,
    nextPoint.lat,
    nextPoint.lon
  );

  const count = Math.trunc(distance / INTERPOLATE_STEP_LENGTH);
  const fraction = 1 / count;

  if (distance > INTERPOLATE_STEP_LENGTH) {
    for (let x = 0; x <= count - 1; x++) {
      let interPoint = nextPointV
        .subtract(pointV)
        .multiply(x * fraction)
        .add(pointV);
        ridingPoints.push({ lat: interPoint.x, lon: interPoint.y });
    }
  } else {
    ridingPoints.push(point);
  }
}

/**
 * Jitter a point on the WGS84 projection, using meters as unit.
 */
function jitterPoint(vector, meters) {
  const lat = vector.x;
  const lon = vector.y;
  const R = 6378137.0; // Earth's mean radius in meters
  const dLat = meters / R;
  const dLon = meters / (R * Math.cos(Math.PI * lat / 180));
  const latO = lat + dLat * 180 / Math.PI;
  const lonO = lon + dLon * 180 / Math.PI;
  return new Vector(latO, lonO);
}

// Jitter all points on the WGS84 projection, using meters as unit.
ridingPoints = ridingPoints.map((point) => {
  // Use a GPS error of [-1,1] meters
  const meters = Math.trunc(Math.random() * 2) - 1;
  const jittered = jitterPoint(new Vector(point.lat, point.lon), meters);
  return { lat: jittered.x, lon: jittered.y };
});


// ======= UPDATE POSITION VIA JSON-RPC =======
console.log("======= UPDATE POSITION VIA JSON-RPC =======");

const socket = new WebSocket(WS_ENDPOINT + sessionId);

socket.onopen = function (e) {
  console.log("[open] Connection established");

  // Start Navigation
  socket.send(
    JSON.stringify({
      jsonrpc: "2.0",
      method: "Navigation",
      params: {
        sessionId: sessionId,
        active: true,
      },
      id: Math.trunc(Math.random() * 1000) + "",
    })
  );

  // Update position every two seconds
  let t = 0;

  const interval = setInterval(() => {
    console.log(
      `sending position (${ridingPoints[t].lat}, ${ridingPoints[t].lon})`
    );

    socket.send(
      JSON.stringify({
        jsonrpc: "2.0",
        method: "PositionUpdate",
        params: {
          lon: ridingPoints[t].lon,
          lat: ridingPoints[t].lat,
          speed: INTERPOLATE_STEP_LENGTH / (POSITION_UPDATE_INTERVAL / 1000),
        },
      })
    );

    // pick next point on route in the next loop
    t++;

    // stop updating and sending stop signal when route is done
    if (t >= ridingPoints.length) {
      console.log("riding done");
      socket.send(
        JSON.stringify({
          jsonrpc: "2.0",
          method: "Navigation",
          params: {
            sessionId: sessionId,
            active: false,
          },
          id: Math.trunc(Math.random() * 1000) + "",
        })
      );
      clearInterval(interval);

      // Shut down the script after 2 Seconds
      setTimeout(() => {
        socket.close();
        Deno.exit();
      }, 2000);
    }
  }, POSITION_UPDATE_INTERVAL);
};

// ======= RECEIVE RECOMMENDATIONS VIA JSON-RPC =======

socket.onmessage = (data) => {
  const msg = JSON.parse(data.data);

  if ((msg.method = "RecommendationUpdate")) {
    console.log("new Recommendation", msg.params);
  }
};

// ======= WEBSOCKET ERROR HANDLING =======

socket.onclose = function (event) {
  if (event.wasClean) {
    console.log(
      `[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`
    );
  } else {
    // e.g. server process killed or network down
    // event.code is usually 1006 in this case
    console.log("[close] Connection died");
  }
};

socket.onerror = function (error) {
  console.log(`[error] ${error.message}`);
};
