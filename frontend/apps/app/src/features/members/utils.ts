import type { ApiInvitation, ApiMember, Member, PendingInvitation } from "./types"

export function formatDate(dateString: string): string {
	const date = new Date(dateString)
	return date
		.toLocaleDateString("en-US", {
			month: "short",
			day: "2-digit",
			year: "numeric",
		})
		.toUpperCase()
}

export function getInitials(name: string): string {
	const parts = name.split(" ")
	if (parts.length >= 2) {
		return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
	}
	return name.substring(0, 2).toUpperCase()
}

export function getTimeAgo(dateString: string): string {
	const now = new Date()
	const date = new Date(dateString)
	const diffMs = now.getTime() - date.getTime()
	const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

	if (diffDays === 0) return "INVITED TODAY"
	if (diffDays === 1) return "INVITED 1 DAY AGO"
	return `INVITED ${diffDays} DAYS AGO`
}

export function mapApiInvitation(invitation: ApiInvitation): PendingInvitation {
	const emailPrefix = invitation.email.split("@")[0]
	return {
		id: invitation.id,
		email: invitation.email.toUpperCase(),
		initials: getInitials(emailPrefix),
		role: invitation.role.toUpperCase(),
		invitedAgo: getTimeAgo(invitation.createdAt),
	}
}

export function mapApiMember(apiMember: ApiMember): Member {
	return {
		name: apiMember.user.name.toUpperCase(),
		email: apiMember.user.email.toUpperCase(),
		initials: getInitials(apiMember.user.name),
		role:
			(apiMember.role.toUpperCase() as "OWNER" | "EDITOR" | "VIEWER") ||
			"VIEWER",
		roleId: apiMember.role,
		joined: formatDate(apiMember.createdAt.toString()),
		isOwner: apiMember.role.toLowerCase() === "owner",
		userId: apiMember.userId,
	}
}
