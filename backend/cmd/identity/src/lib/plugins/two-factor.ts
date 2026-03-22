import { twoFactor } from "better-auth/plugins";

export const twoFactorPlugin = twoFactor({
  issuer: "Zero Budget",
  schema: {
    twoFactor: {
      fields: {
        userId: "user_id",
      },
    },
    user: {
      fields: {
        twoFactorEnabled: "two_factor_enabled",
      },
    },
  },
});
