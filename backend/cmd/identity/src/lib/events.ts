const getGoApiUrl = () => process.env.GO_API_URL || "http://localhost:8080";
const getInternalApiKey = () => process.env.INTERNAL_API_KEY;

async function publishEvent(event: string, payload: Record<string, unknown>) {
  const goApiUrl = getGoApiUrl();
  const internalApiKey = getInternalApiKey();

  if (!internalApiKey) {
    console.warn(
      "[identity] INTERNAL_API_KEY not set, skipping event publish",
    );
    return;
  }

  const eventPayload = { event, payload };
  console.log("[identity] publishing event to Go API", eventPayload);

  try {
    const res = await fetch(`${goApiUrl}/v1/events/publish`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-API-Key": internalApiKey,
      },
      body: JSON.stringify(eventPayload),
    });

    console.log("[identity] Go API response", {
      status: res.status,
      statusText: res.statusText,
    });

    if (!res.ok) {
      const body = await res.text();
      console.error("[identity] Go API error response body:", body);
    }
  } catch (err) {
    console.error(`[identity] Failed to publish ${event} event:`, err);
  }
}

export async function publishPasswordReset(user: {
  id: string;
  email: string;
  name: string;
}, url: string) {
  await publishEvent("user.password_reset", {
    userId: user.id,
    email: user.email,
    name: user.name,
    resetUrl: url,
  });
}

export async function publishVerificationEmail(user: {
  id: string;
  email: string;
  name: string;
}, url: string) {
  await publishEvent("user.verification_email", {
    userId: user.id,
    email: user.email,
    name: user.name,
    verificationUrl: url,
  });
}

export async function publishUserSignedUp(user: {
  id: string;
  email: string;
  name: string;
}) {
  await publishEvent("user.signed_up", {
    userId: user.id,
    email: user.email,
    name: user.name,
  });
}
