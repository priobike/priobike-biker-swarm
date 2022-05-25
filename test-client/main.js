import { Vector, distanceP2P } from "./utils.js";

const AUTH_ENDPOINT = "http://priobike.vkw.tu-dresden.de/production/session-wrapper/authentication";
const WS_ENDPOINT = "ws://priobike.vkw.tu-dresden.de/production/session-wrapper/websocket/sessions/";
const ROUTE_ENDPOINT = "http://priobike.vkw.tu-dresden.de/production/session-wrapper/getroute";

// IMPORTANT: DONT CHANGE THIS ID. Test clients are hidden in the session analyzer and removed from persistence after session closure.
const CLIENT_ID = "JS_TEST_CLIENT";

const INTERPOLATE_STEP_LENGTH = 4; // m
const POSITION_UPDATE_INTERVAL = 1000; // ms

// Wait a random amount of seconds between 0 and 60 before starting the ride
const waitTime = Math.floor(Math.random() * 60_000);
console.log(`Waiting ${waitTime} milliseconds before starting ride`);
setTimeout(async function() {
  console.log("Starting ride");
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

  // TS1: Shell-Tankstelle Theodor-Heuss-Platz über Edmund-Siemers-Allee nach An d. Verbindungsbahn
  // Signalgruppen (Stand: 18. Mai 2022): 271_31, 271_33, 279_24, 279_25, 413_2
  // https://www.google.de/maps/dir/'53.56090651069028,9.99085153989529'/'53.564343128008701,9.978573058782963'/@53.5625761,9.9813999,16z/data=!3m1!4b1!4m10!4m9!1m3!2m2!1d9.9908515!2d53.5609065!1m3!2m2!1d9.9785731!2d53.5643431!3e1
  const testRoute1 = [
    {
      'lon': 9.99085153989529, 'lat': 53.56090651069028
    },
    {
      'lon': 9.978573058782963, 'lat': 53.564343128008701
    },
  ];

  // TS1 (Rückzu): An d. Verbindungsbahn über Edmund-Siemers-Allee nach McDonalds Theodor-Heuss-Platz
  // Signalgruppen (Stand: 18. Mai 2022): 413_3, 2150_8, 271_1 (befindet sich ganz am Ende der Strecke und kann auch wegfallen)
  // https://www.google.de/maps/dir/'53.564123305799505,9.978725860677336'/'53.561170720637968,9.989768344244068'/@53.5626634,9.9797544,16z/data=!3m1!4b1!4m10!4m9!1m3!2m2!1d9.9787259!2d53.5641233!1m3!2m2!1d9.9897683!2d53.5611707!3e1
  const testRoute1Back = [
    {
      'lon': 9.978725860677336, 'lat': 53.564123305799505
    },
    {
      'lon': 9.989768344244068, 'lat': 53.561170720637968
    },
  ];

  // TS2: Busparkplatz Willy-Brandt-Straße über B4 nach Millerntorwache
  // Signalgruppen (Stand: 18. Mai 2022): 70_27, 76_18, 76_20, 256_27, 333_12, 92_37, 92_38, 193_30, 193_32
  // https://www.google.de/maps/dir/'53.547722154285324,10.004045134575035'/'53.550264133830126,9.971739418506827'/@53.5469702,9.9965353,16.18z/data=!4m20!4m19!1m13!2m2!1d10.0040451!2d53.5477222!3m4!1m2!1d9.9968626!2d53.5475663!3s0x47b18f1d4fa210bf:0x7164391b88468bbd!3m4!1m2!1d9.9882061!2d53.547483!3s0x47b18f04d1f5667d:0x1d6d663ff6ea0f94!1m3!2m2!1d9.9717394!2d53.5502641!3e1!5m2!1e4!1e3
  const testRoute2 = [
    {
      'lon': 10.004045134575035, 'lat': 53.547722154285324
    },
    // Eingefügter Wegpunkt für Signalgruppen 92_37 und 92_38 (Kreuzung Neanderstraße, Radweg wurde nicht genutzt)
    {
      'lon': 9.978636, 'lat': 53.549482
    },
    {
      'lon': 9.971739418506827, 'lat': 53.550264133830126
    },
  ];

  // TS2 (rückzu): Millerntorpl./Alter Elbpark über B4 nach Deichtor Office Zentrum
  // Signalgruppen (Stand: 18. Mai 2022): 92_43, 92_44, 310_3, 333_29, 333_30, 256_30, 532_33, 70_42, 70_43
  // https://www.google.de/maps/dir/+53.54990402934412,9.971606990198367/'53.547262160720436,10.004240381440082'/@53.5472429,10.0037521,109m/data=!3m1!1e3!4m10!4m9!1m3!2m2!1d9.971607!2d53.549904!1m3!2m2!1d10.0042404!2d53.5472622!3e1!5m2!1e4!1e3
  const testRoute2Back = [
    {
      'lon': 9.971606990198367, 'lat': 53.54990402934412
    },
    {
      'lon': 10.004240381440082, 'lat': 53.547262160720436
    },
  ];

  // TS3: Nordseite - Hauptbahnhof/Steintorwall über Lombardsbrücke, Stephansplatz nach Holstenwall 13
  // Signalgruppen (Stand: 25. Mai 2022): 
  // 77_2, 535_23, 535_24, 200_12, 200_31, 119_42, 119_33, 119_36, 119_39, 118_63, 118_64, 104_5, 257_42, 257_2
  // Detektiert 257_8 statt 257_2 fälschlicherweise
  // https://www.google.com/maps/dir/53.5511715,10.0062077/53.5575131,9.99471/53.5575762,9.9828379/53.55285,+9.976352/@53.5532471,9.9828094,15.59z/data=!4m11!4m10!1m0!1m0!1m0!1m5!1m1!1s0x0:0x3f367588de00e142!2m2!1d9.9765066!2d53.5528316!3e1
  const testRoute3 = [
    {
      'lon': 10.0062077, 'lat': 53.5511715
    },
    {
      'lon': 9.99471, 'lat': 53.5575131
    },
    {
      'lon': 9.9828379, 'lat': 53.5575762
    },
    {
      'lon': 9.976352, 'lat': 53.55285
    },
  ];

  // TS3 (rückzu): Südseite - Holstenwall 13 über Lombardsbrücke nach Hauptbahnhof/Steintorwall
  // Signalgruppen (Stand: 25. Mai 2022):
  // 257_3, 104_1, 118_14/118_16/118_13, (119_56), (119_52), 200_51, 535_35, 535_36, (77_32)
  // Detektiert 257_1 fälschlicherweise
  // https://www.google.com/maps/dir/53.55285,+9.976352/53.5579687,9.9859757/'53.551241482916915,10.005804047062561'/@53.5511614,10.00482,18.3z/data=!4m20!4m19!1m10!1m1!1s0x0:0x3f367588de00e142!2m2!1d9.976352!2d53.55285!3m4!1m2!1d9.9790178!2d53.5552125!3s0x47b18f16a0edd3cb:0x54538b05a0342812!1m0!1m5!1m1!1s0x0:0x77ea443ba90e82b2!2m2!1d10.005804!2d53.5512415!3e1
  const testRoute3Back = [
    {
      'lon': 9.976352, 'lat': 53.55285
    },
    {
      'lon': 9.9859757, 'lat': 53.5579687
    },
    {
      'lon': 10.005804047062561, 'lat': 53.551241482916915
    }
  ];

  // Select a random route
  const testRoutes = [testRoute1, testRoute1Back, testRoute2, testRoute2Back, testRoute3, testRoute3Back];
  const randomRoute = testRoutes[Math.floor(Math.random() * testRoutes.length)];

  const routingResponse = await fetch(ROUTE_ENDPOINT, {
    method: "POST",
    body: JSON.stringify({
      sessionId,
      waypoints: randomRoute
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

      const date = new Date();
      const iso8601 = date.toISOString();

      socket.send(
        JSON.stringify({
          jsonrpc: "2.0",
          method: "PositionUpdate",
          params: {
            lon: ridingPoints[t].lon,
            lat: ridingPoints[t].lat,
            speed: INTERPOLATE_STEP_LENGTH / (POSITION_UPDATE_INTERVAL / 1000),
            timestamp: iso8601,
            heading: 0.0, // 0.0 is a default value
            accuracy: 0.0, // 0.0 is a default value
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
}, waitTime);