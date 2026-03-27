
export function sleep(delay) {
  return new Promise((resolve) => setTimeout(resolve, delay))
}

export function bytesToSize(bytes) {
  const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
  if (bytes === 0) {
    return "0 Bytes";
  }
  const regex = /(.+)\.00$/;
  const regex2 = /(.+)\.0$/;
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i))
    .toFixed(2)
    .replace(regex, "$1")
    .replace(regex2, "$1")} ${sizes[i]}`;
}

export function getPageVersion() {
  return sessionStorage.getItem('page-version') ?? 'elem';
}

export function setPageVersion(version = 'elem') {
  return sessionStorage.setItem('page-version', version);
}