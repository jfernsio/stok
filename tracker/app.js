let ws = null;
let currentSymbol = "AAPL";

const liveDiv = document.getElementById("liveCandle");
const historyTable = document.getElementById("historyTable");
const symbolSelect = document.getElementById("symbolSelect");

document.getElementById("connectBtn").onclick = () => {
  currentSymbol = symbolSelect.value;
  connectWS();
  loadHistory();
};

function connectWS() {
  if (ws) ws.close();

  ws = new WebSocket("ws://localhost:8080/ws");

  ws.onopen = () => {
    ws.send(currentSymbol);
    console.log("WS connected");
  };

  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    renderLive(msg);
  };

  ws.onclose = () => console.log("WS closed");
}

function renderLive(msg) {
  const c = msg.candle;

  liveDiv.innerHTML = `
    <div>Type:</div><div>${msg.type}</div>
    <div>Open:</div><div>${c.open.toFixed(2)}</div>
    <div>High:</div><div>${c.high.toFixed(2)}</div>
    <div>Low:</div><div>${c.low.toFixed(2)}</div>
    <div>Close:</div><div>${c.close.toFixed(2)}</div>
    <div>Volume:</div><div>${c.volume.toFixed(0)}</div>
  `;
}

async function loadHistory() {
  historyTable.innerHTML = "";

  const res = await fetch(
    `http://localhost:8080/stock-candles?symbol=${currentSymbol}`
  );
  const data = await res.json();

  data.slice(0, 10).forEach(c => {
    const row = document.createElement("tr");
    row.className = "border-t border-slate-700";

    row.innerHTML = `
      <td>${new Date(c.openTime).toLocaleTimeString()}</td>
      <td>${c.open.toFixed(2)}</td>
      <td>${c.high.toFixed(2)}</td>
      <td>${c.low.toFixed(2)}</td>
      <td>${c.close.toFixed(2)}</td>
      <td>${c.volume.toFixed(0)}</td>
    `;

    historyTable.appendChild(row);
  });
}
