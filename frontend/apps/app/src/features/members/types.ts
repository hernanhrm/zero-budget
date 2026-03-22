export interface Member {
	name: string
	email: string
	initials: string
	role: "OWNER" | "EDITOR" | "VIEWER"
	roleId: string
	joined: string
	isOwner: boolean
	userId: string
}

export interface MembersTableProps {
	members: Member[]
	isLoading: boolean
	error: string | null
}

export interface PendingInvitation {
	email: string
	initials: string
	role: string
	invitedAgo: string
}

export interface ApiMember {
	userId: string
	user: { name: string; email: string }
	role: string
	createdAt: string
}
