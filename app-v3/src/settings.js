export default {
  title: "app-v3",
  /**
   * @type {boolean} true | false
   * @description
   */
  nodeTest: process.env.NODE_ENV === "development",

  /**
   * @type {boolean} true | false
   * @description Whether fix the header
   */
  fixedHeader: false,

  /**
   * @type {boolean} true | false
   * @description Whether show the logo in sidebar
   */
  sidebarLogo: false,
  /**
   * @type {}
   * @description  device info
   */
  appInfo: {
    version: "",
    build: "",
  },
  deviceInfo: {
    g_login_user: "admin", //username
    g_is_au: 1, //Whether it is the primary AU
    tokenRole: "admin", //User Role Read-Only Read-Write root
    deviceType: 240,
    name: "webomt", //Device model name, login page displayed
    deviceID: 0,
    repeaterType: null,
    deviceNmae: "",
    master_device_name: "",
    site_name: "",
    currentType: "Primary A3",
    currentVersion: "0.12",
    // currentType: "N3-RU",
    // currentVersion: "0.11",
  },
};
