export interface AccountRow {
	id: string
	name: string
	mask: string
	type: string
	institution: string
	balance: string
	balanceIsPrimary: boolean
}

export const MOCK_ACCOUNT_ROWS: AccountRow[] = [
	{
		id: "1",
		name: "PRIMARY CHECKING",
		mask: "****4821",
		type: "CHECKING",
		institution: "CHASE BANK",
		balance: "$6,122.18",
		balanceIsPrimary: false,
	},
	{
		id: "2",
		name: "BILLS & EXPENSES",
		mask: "****7293",
		type: "CHECKING",
		institution: "CHASE BANK",
		balance: "$2,112.00",
		balanceIsPrimary: false,
	},
	{
		id: "3",
		name: "EMERGENCY FUND",
		mask: "****1056",
		type: "SAVINGS",
		institution: "ALLY BANK",
		balance: "$4,613.34",
		balanceIsPrimary: true,
	},
]
