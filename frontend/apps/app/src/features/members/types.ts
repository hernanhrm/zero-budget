export interface Member {
	id: string
	name: string
	email: string
	initials: string
	role: "OWNER" | "ADMIN" | "MEMBER"
	roleId: string
	joined: string
	isOwner: boolean
	userId: string
}

export interface MembersTableProps {
	members: Member[]
	isLoading: boolean
	error: string | null
	currentUserId: string
	onSuccess: () => void
}

export interface PendingInvitation {
	id: string
	email: string
	initials: string
	role: string
	invitedAgo: string
}

export interface ApiInvitation {
	id: string
	email: string
	organizationId: string
	role: string
	status: string
	inviterId: string
	createdAt: string
	expiresAt: string
}

export interface ApiMember {
	id: string
	userId: string
	user: { name: string; email: string }
	role: string
	createdAt: string | Date
}
