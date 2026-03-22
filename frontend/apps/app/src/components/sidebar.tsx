import { Link, useMatchRoute } from "@tanstack/react-router"

const navItems = [
	{ label: "BUDGET", to: "/" },
	{ label: "TRANSACTIONS", to: "/transactions" },
	{ label: "ACCOUNTS", to: "/accounts" },
	{ label: "MEMBERS", to: "/members" },
	{ label: "SETTINGS", to: "/settings" },
] as const

function NavItem({ label, to }: { label: string; to: string }) {
	const matchRoute = useMatchRoute()
	const isActive = matchRoute({ to, fuzzy: true })

	return (
		<Link
			to={to}
			className={`flex h-12 items-center gap-4 px-6 ${
				isActive
					? "border-l-[3px] border-[var(--zb-accent)]"
					: ""
			}`}
		>
			<div
				className={`h-2 w-2 ${isActive ? "bg-[var(--zb-accent)]" : "bg-[var(--zb-text-muted)]"}`}
			/>
			<span
				className={`font-space-grotesk text-xs tracking-[1.5px] ${
					isActive
						? "font-bold text-[var(--zb-text-primary)]"
						: "font-medium text-[var(--zb-text-secondary)]"
				}`}
			>
				{label}
			</span>
		</Link>
	)
}

export function Sidebar() {
	return (
		<aside className="flex h-full w-[260px] shrink-0 flex-col justify-between border-r border-[var(--zb-border)]">
			<div className="flex flex-col">
				<div className="flex items-center gap-3 border-b border-[var(--zb-border)] px-6 py-6">
					<div className="flex h-9 w-9 items-center justify-center bg-[var(--zb-accent)]">
						<span className="font-space-grotesk text-lg font-bold text-[var(--zb-bg)]">
							B
						</span>
					</div>
					<span className="font-space-grotesk text-base font-bold tracking-[2px] text-[var(--zb-accent)]">
						BUDGETFORGE
					</span>
				</div>
				<nav className="flex flex-col gap-0.5 py-6">
					{navItems.map((item) => (
						<NavItem
							key={item.label}
							label={item.label}
							to={item.to}
						/>
					))}
				</nav>
			</div>
			<div className="border-t border-[var(--zb-border)] px-6 pb-6">
				<div className="flex h-16 items-center gap-3 pt-4">
					<div className="h-2.5 w-2.5 bg-[var(--zb-accent)]" />
					<span className="font-space-grotesk text-xs font-bold tracking-[1px] text-[var(--zb-text-primary)]">
						JOHN.DOE
					</span>
				</div>
			</div>
		</aside>
	)
}
