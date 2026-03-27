import { openApiRequest as request } from "@/utils/request";
const pako = require('pako');

export const ProxyBase = "/proxy";
export const ApiBase = "/api";
export const DasApiBase = ApiBase + "/das";

export async function httpGet(u, body, options) {
  return request(u, {
    method: 'GET',
    data: body,
    ...(options || {}),
  })
}

// export async function httpGetGzip(u) {
//     const options = { responseType: 'arraybuffer' } // Specify the response type as array buffer to get the raw binary data
//     const resp = await request(u, {
//         method: 'GET',
//         ...(options || {}),
//     })
//     const data = JSON.parse(pako.inflate(resp, { to: 'string' }));
//     return data
// }

export async function httpPost(u, body, options) {
  return request(u, {
    method: 'POST',
    data: body,
    ...(options || {}),
  })
}

export async function httpDelete(u, body, options) {
  return request(u, {
    method: 'DELETE',
    data: body,
    ...(options || {}),
  })
}