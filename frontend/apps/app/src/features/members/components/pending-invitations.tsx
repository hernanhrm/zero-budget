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
			className={`flex h-[52px] items-center justify-between px-6 ${!isLast ? "border-b border-border" : ""}`}
		>
			<div className="flex items-center gap-3">
				<div className="flex h-8 w-8 items-center justify-center bg-border">
					<span className="font-space-grotesk text-[11px] font-bold text-muted-foreground">
						{invitation.initials}
					</span>
				</div>
				<div className="flex flex-col gap-0.5">
					<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-foreground">
						{invitation.email}
					</span>
					<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
						{invitation.invitedAgo}
					</span>
				</div>
			</div>
			<div className="flex items-center gap-3">
				<div className="flex h-6 items-center justify-center border border-border px-3">
					<span className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-muted-foreground">
						{invitation.role}
					</span>
				</div>
				<button
					type="button"
					className="font-space-grotesk text-[11px] font-bold tracking-[1px] text-primary hover:opacity-80"
				>
					RESEND
				</button>
				<button
					type="button"
					className="text-muted-foreground hover:text-destructive"
				>
					<X className="size-3.5" />
				</button>
			</div>
		</div>
	)
}

export function PendingInvitations() {
	return (
		<div className="w-full border border-border">
			<div className="flex h-14 items-center justify-between bg-card px-6 border-b border-border">
				<div className="flex items-center gap-3">
					<div className="h-5 w-1 bg-primary" />
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
						PENDING INVITATIONS
					</span>
					<div className="flex h-5 w-6 items-center justify-center bg-primary">
						<span className="font-space-grotesk text-[11px] font-bold text-primary-foreground">
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