
import * as base from "./api-base";
import * as iam from "./api-iam";
import * as user from "./api-user";
import * as dasWebSocket from "./api-das-ws";
import * as dasModel from "./api-das-model";
import * as dasDevice from "./api-das-device";
import * as dasParams from "./api-das-params";
import * as dasUpgrade from "./api-das-upgrade";
import * as dasFile from "./api-das-file";
import * as localDevice from "./api-local";
import * as ping from "./api-ping";

export default {
  ...base,
  ...iam,
  ...user,
  ...dasWebSocket,
  ...dasModel,
  ...dasDevice,
  ...dasParams,
  ...dasFile,
  ...dasUpgrade,
  ...localDevice,
  ...ping
} ;