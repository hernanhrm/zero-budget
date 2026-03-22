import { X } from "lucide-react"

interface PendingInvitation {
	email: string
	initials: string
	role: string
	invitedAgo: string
}

const pendingInvitations: PendingInvitation[] = [
	{
		email: "SARAH.MILLER@EMAIL.COM",
		initials: "SM",
		role: "EDITOR",
		invitedAgo: "INVITED 2 DAYS AGO",
	},
	{
		email: "MIKE.REYES@EMAIL.COM",
		initials: "MR",
		role: "VIEWER",
		invitedAgo: "INVITED 5 DAYS AGO",
	},
]

function PendingInvitationRow({
	invitation,
	isLast,
}: {
	invitation: PendingInvitation
	isLast: boolean
}) {
	return (
		<div
			className={`flex h-[52px] items-center justify-between px-6 ${!isLast ? "border-b border-[var(--zb-border)]" : ""}`}
		>
			<div className="flex items-center gap-3">
				<div className="flex h-8 w-8 items-center justify-center bg-[var(--zb-border)]">
					<span className="font-space-grotesk text-[11px] font-bold text-[var(--zb-text-secondary)]">
						{invitation.initials}
					</span>
				</div>
				<div className="flex flex-col gap-0.5">
					<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-[var(--zb-text-primary)]">
						{invitation.email}
					</span>
					<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-[var(--zb-text-muted)]">
						{invitation.invitedAgo}
					</span>
				</div>
			</div>
			<div className="flex items-center gap-3">
				<div className="flex h-6 items-center justify-center border border-[var(--zb-border)] px-3">
					<span className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-[var(--zb-text-secondary)]">
						{invitation.role}
					</span>
				</div>
				<button
					type="button"
					className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-[var(--zb-accent)] hover:opacity-80"
				>
					RESEND
				</button>
				<button
					type="button"
					className="text-[var(--zb-text-secondary)] hover:text-[var(--zb-danger)]"
				>
					<X className="size-3.5" />
				</button>
			</div>
		</div>
	)
}

export function PendingInvitations() {
	return (
		<div className="w-full border border-[var(--zb-border)]">
			<div className="flex h-14 items-center justify-between bg-[var(--zb-bg-elevated)] px-6 border-b border-[var(--zb-border)]">
				<div className="flex items-center gap-3">
					<div className="h-5 w-1 bg-[var(--zb-accent)]" />
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-[var(--zb-text-primary)]">
						PENDING INVITATIONS
					</span>
					<div className="flex h-5 w-6 items-center justify-center bg-[var(--zb-accent)]">
						<span className="font-space-grotesk text-[11px] font-bold text-[var(--zb-bg)]">
							{pendingInvitations.length}
						</span>
					</div>
				</div>
			</div>
			{pendingInvitations.map((invitation, index) => (
				<PendingInvitationRow
					key={invitation.email}
					invitation={invitation}
					isLast={index === pendingInvitations.length - 1}
				/>
			))}
		</div>
	)
}
