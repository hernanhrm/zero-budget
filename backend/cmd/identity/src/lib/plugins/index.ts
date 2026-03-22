import { openAPI } from "better-auth/plugins";
import { dash } from "@better-auth/infra";
import { organizationPlugin } from "./organization.js";
import { twoFactorPlugin } from "./two-factor.js";
import { apiKeyPlugin } from "./api-key.js";

export const plugins = [
  openAPI(),
  dash(),
  organizationPlugin,
  twoFactorPlugin,
  apiKeyPlugin,
];
