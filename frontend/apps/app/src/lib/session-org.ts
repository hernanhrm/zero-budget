// biome-ignore lint/style/useImportType: `typeof authClient.getSession` requires the runtime client value.
import { authClient } from "#/lib/auth-client"

export type AppSession = Awaited<ReturnType<typeof authClient.getSession>>

/** Active organization id from Better Auth session payload (organization plugin). */
export function getActiveOrganizationId(
	session: AppSession | null,
): string | undefined {
	const data = session?.data as
		| { session?: { activeOrganizationId?: string | null } }
		| undefined
	return data?.session?.activeOrganizationId ?? undefined
}
