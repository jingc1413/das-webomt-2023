
import { DasApiBase } from "./api-base";

export async function connectToWebSocket({ onMessage, onClose, onOpen, onError }) {
  try {
    const host = window.location.host;
    const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
    const ws = new WebSocket(`${proto}//${host}${DasApiBase}/ws`);
    ws.onmessage = onMessage;
    ws.onclose = onClose;
    ws.onopen = onOpen;
    ws.onerror = onError;
    return ws;
  } catch (e) {
    console.error(e)
    return undefined
  }
}