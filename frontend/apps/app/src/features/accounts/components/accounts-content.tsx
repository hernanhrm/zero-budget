import { useSuspenseQuery } from "@tanstack/react-query"
import type { Account } from "@workspace/api"
import {
	getGetV1AccountsQueryKey,
	getV1Accounts,
} from "@workspace/api/hooks/accounts/accounts"
import { toast } from "@workspace/ui/lib/toast"
import { useEffect, useMemo, useRef } from "react"
import { mapAccountsToRows } from "../account-row"
import { AccountsTable } from "./accounts-table"
import { MetricCards } from "./metric-cards"

/**
 * Go API wraps JSON as `{ "data": <payload> }` (see backend/pkg/httpresponse).
 * Orval types assume the raw body is the list; unwrap the envelope when present.
 */
function normalizeAccountsPayload(parsedBody: unknown): Account[] {
	if (parsedBody == null) {
		return []
	}
	if (Array.isArray(parsedBody)) {
		return parsedBody as Account[]
	}
	if (
		typeof parsedBody === "object" &&
		parsedBody !== null &&
		"data" in parsedBody
	) {
		const inner = (parsedBody as { data: unknown }).data
		if (inner == null || !Array.isArray(inner)) {
			return []
		}
		return inner as Account[]
	}
	return []
}

export function AccountsContent() {
	const lastAccountsErrorToast = useRef<string | null>(null)

	const { data: res } = useSuspenseQuery({
		queryKey: getGetV1AccountsQueryKey(),
		queryFn: ({ signal }) =>
			getV1Accounts(undefined, { signal, credentials: "include" }),
	})

	const accounts = useMemo((): Account[] => {
		if (!res || res.status !== 200) {
			return []
		}
		return normalizeAccountsPayload(res.data)
	}, [res])

	const rows = useMemo(() => mapAccountsToRows(accounts), [accounts])

	const listError = useMemo((): string | null => {
		if (res && res.status !== 200) {
			return `Could not load accounts (${res.status}).`
		}
		return null
	}, [res])

	useEffect(() => {
		if (!listError) {
			lastAccountsErrorToast.current = null
			return
		}
		if (lastAccountsErrorToast.current === listError) {
			return
		}
		lastAccountsErrorToast.current = listError
		toast.error(listError)
	}, [listError])

	return (
		<>
			<MetricCards accounts={accounts} />

			<div className="min-w-0 overflow-x-auto">
				<AccountsTable rows={rows} />
			</div>
		</>
	)
}
