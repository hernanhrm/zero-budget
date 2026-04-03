import { openAPI } from "better-auth/plugins";
import { organizationPlugin } from "./organization.js";
import { twoFactorPlugin } from "./two-factor.js";
import { apiKeyPlugin } from "./api-key.js";

export const plugins = [
  openAPI(),
  organizationPlugin,
  twoFactorPlugin,
  apiKeyPlugin,
];
