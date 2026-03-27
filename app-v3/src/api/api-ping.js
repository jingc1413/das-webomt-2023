import * as base from "./api-base";



export async function createPingJob(info) {
  return base.httpPost(base.ApiBase + '/diag/ping/jobs', info);
}

// open ws  /api/diag/ping/jobs/{token}/ws
export function connectPingJob(token, { onMessage, onClose, onOpen, onError }) {
  try {
    let host = window.location.host;
    let proto = window.location.protocol === "https:" ? "wss:" : "ws:";
    let ws = new WebSocket(`${proto}//${host}${base.ApiBase}/diag/ping/jobs/${token}/ws`);
    ws.onmessage = onMessage;
    ws.onclose = onClose;
    ws.onopen = onOpen;
    ws.onerror = onError;
    return ws;
  } catch (e) {
    return undefined
  }
}

export async function runPingJob(token) {
  return base.httpPost(base.ApiBase + `/diag/ping/jobs/${token}/run`, null);
}

export async function stopPingJob(token) {
  return base.httpPost(base.ApiBase + `/diag/ping/jobs/${token}/cancel`, null);
}
