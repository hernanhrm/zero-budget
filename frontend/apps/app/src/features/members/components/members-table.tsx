import { ChevronDown, Trash2 } from "lucide-react"

interface Member {
	name: string
	email: string
	initials: string
	role: "OWNER" | "EDITOR" | "VIEWER"
	joined: string
	isOwner: boolean
}

const members: Member[] = [
	{
		name: "JOHN DOE",
		email: "JOHN.DOE@EMAIL.COM",
		initials: "JD",
		role: "OWNER",
		joined: "JAN 15, 2026",
		isOwner: true,
	},
	{
		name: "JANE SMITH",
		email: "JANE.SMITH@EMAIL.COM",
		initials: "JS",
		role: "EDITOR",
		joined: "FEB 03, 2026",
		isOwner: false,
	},
	{
		name: "ALEX WONG",
		email: "ALEX.WONG@EMAIL.COM",
		initials: "AW",
		role: "VIEWER",
		joined: "MAR 10, 2026",
		isOwner: false,
	},
]

function RoleBadge({ role, isOwner }: { role: string; isOwner: boolean }) {
	if (isOwner) {
		return (
			<div className="flex h-6 w-40 items-center justify-center gap-1.5 px-3">
				<div className="h-1.5 w-1.5 rounded-full bg-[var(--zb-accent)]" />
				<span className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-[var(--zb-accent)]">
					{role}
				</span>
			</div>
		)
	}

	return (
		<button
			type="button"
			className="flex h-7 w-40 items-center justify-center gap-1.5 border border-[var(--zb-border)] px-3 hover:border-[var(--zb-text-secondary)]"
		>
			<span className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-[var(--zb-text-primary)]">
				{role}
			</span>
			<ChevronDown className="size-2.5 text-[var(--zb-text-secondary)]" />
		</button>
	)
}

function MemberRow({
	member,
	isLast,
}: {
	member: Member
	isLast: boolean
}) {
	return (
		<div
			className={`flex h-16 items-center px-6 ${!isLast ? "border-b border-[var(--zb-border)]" : ""}`}
		>
			<div className="flex flex-1 items-center gap-3">
				<div
					className={`flex h-9 w-9 items-center justify-center ${member.isOwner ? "bg-[var(--zb-accent)]" : "bg-[var(--zb-border)]"}`}
				>
					<span
						className={`font-space-grotesk text-xs font-bold ${member.isOwner ? "text-[var(--zb-bg)]" : "text-[var(--zb-text-secondary)]"}`}
					>
						{member.initials}
					</span>
				</div>
				<div className="flex flex-col gap-0.5">
					<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-[var(--zb-text-primary)]">
						{member.name}
					</span>
					<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-[var(--zb-text-muted)]">
						{member.email}
					</span>
				</div>
			</div>
			<RoleBadge role={member.role} isOwner={member.isOwner} />
			<span className="w-[140px] font-ibm-plex-mono text-xs tracking-[1px] text-[var(--zb-text-secondary)]">
				{member.joined}
			</span>
			<div className="flex w-20 items-center justify-center">
				{member.isOwner ? (
					<span className="font-space-grotesk text-sm font-bold text-[var(--zb-text-muted)]">
						—
					</span>
				) : (
					<button
						type="button"
						className="text-[var(--zb-danger)] hover:opacity-80"
					>
						<Trash2 className="size-4" />
					</button>
				)}
			</div>
		</div>
	)
}

export function MembersTable() {
	return (
		<div className="w-full border border-[var(--zb-border)]">
			<div className="flex h-14 items-center justify-between bg-[var(--zb-bg-elevated)] px-6 border-b border-[var(--zb-border)]">
				<div className="flex items-center gap-3">
					<div className="h-5 w-1 bg-[var(--zb-accent)]" />
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-[var(--zb-text-primary)]">
						ACTIVE MEMBERS
					</span>
					<span className="font-ibm-plex-mono text-xs text-[var(--zb-text-secondary)]">
						{members.length}
					</span>
				</div>
			</div>
			<div className="flex h-10 items-center px-6 border-b border-[var(--zb-border)]">
				<span className="flex-1 font-space-grotesk text-[11px] font-bold tracking-[1px] text-[var(--zb-text-secondary)]">
					MEMBER
				</span>
				<span className="w-40 font-space-grotesk text-[11px] font-bold tracking-[1px] text-[var(--zb-text-secondary)]">
					ROLE
				</span>
				<span className="w-[140px] font-space-grotesk text-[11px] font-bold tracking-[1px] text-[var(--zb-text-secondary)]">
					JOINED
				</span>
				<span className="w-20 text-center font-space-grotesk text-[11px] font-bold tracking-[1px] text-[var(--zb-text-secondary)]">
					ACTIONS
				</span>
			</div>
			{members.map((member, index) => (
				<MemberRow
					key={member.email}
					member={member}
					isLast={index === members.length - 1}
				/>
			))}
		</div>
	)
}
