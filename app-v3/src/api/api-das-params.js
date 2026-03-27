
import * as base from "./api-base";

export async function getParameterValues({ sub, params, oids, values = [] }) {
  const body = [];
  oids.forEach(oid => {
    const param = params.getParam(oid);
    if (param && param.Readable) {
      const value = values.find(v => v.oid === oid);
      if (value != undefined) {
        body.push({ ID: oid, Value: value.value })
      } else {
        body.push({ ID: oid })
      }
    }
  })
  return new Promise((resolve, reject) => {
    base.httpPost(base.DasApiBase + "/devices/" + sub + "/parameters/get", body)
      .then(data => {
        const out = data.map(item => {
          return { oid: item.ID, value: item.Value, code: item.Code };
        });
        return resolve(out);
      })
      .catch(e => {
        return reject(e);
      })
  })
}

export async function setParameterValues({ sub, params, values }) {
  const body = [];
  values.forEach(v => {
    const param = params.getParam(v.oid);
    if (param && param.Writable && v.value !== undefined) {
      body.push({ ID: v.oid, Value: v.value })
    }
  })
  return new Promise((resolve, reject) => {
    base.httpPost(base.DasApiBase + "/devices/" + sub + "/parameters/set", body)
      .then(data => {
        const out = data.map(item => {
          return { oid: item.ID, value: item.Value, code: item.Code };
        });
        return resolve(out);
      })
      .catch(e => {
        return reject(e);
      })
  })
}