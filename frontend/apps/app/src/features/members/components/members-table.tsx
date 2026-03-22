import { ChevronDown, Trash2 } from "lucide-react"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@workspace/ui/components/select"

interface Member {
	name: string
	email: string
	initials: string
	role: "OWNER" | "EDITOR" | "VIEWER"
	roleId: string
	joined: string
	isOwner: boolean
	userId: string
}

interface MembersTableProps {
	members: Member[]
	isLoading: boolean
	error: string | null
}

function formatDate(dateString: string): string {
	const date = new Date(dateString)
	return date.toLocaleDateString("en-US", {
		month: "short",
		day: "2-digit",
		year: "numeric",
	}).toUpperCase()
}

function getInitials(name: string): string {
	const parts = name.split(" ")
	if (parts.length >= 2) {
		return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
	}
	return name.substring(0, 2).toUpperCase()
}

function mapApiMember(apiMember: {
	userId: string
	user: { name: string; email: string }
	role: string
	createdAt: string
}): Member {
	return {
		name: apiMember.user.name.toUpperCase(),
		email: apiMember.user.email.toUpperCase(),
		initials: getInitials(apiMember.user.name),
		role: (apiMember.role.toUpperCase() as "OWNER" | "EDITOR" | "VIEWER") || "VIEWER",
		roleId: apiMember.role,
		joined: formatDate(apiMember.createdAt.toString()),
		isOwner: apiMember.role.toLowerCase() === "owner",
		userId: apiMember.userId,
	}
}

function OwnerBadge({ role }: { role: string }) {
	return (
		<div className="flex h-6 items-center justify-start gap-1.5 pr-3">
			<div className="h-1.5 w-1.5 rounded-full bg-primary" />
			<span className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-primary">
				{role}
			</span>
		</div>
	)
}

function RoleSelect({
	role,
	userId,
}: {
	role: "OWNER" | "EDITOR" | "VIEWER"
	userId: string
}) {
	const [value, setValue] = useState(role)

	const handleValueChange = (newValue: string) => {
		setValue(newValue as "OWNER" | "EDITOR" | "VIEWER")
	}

	return (
		<Select value={value} onValueChange={handleValueChange}>
			<SelectTrigger className="flex h-6 w-fit items-center justify-between gap-1.5 rounded-none border border-border bg-transparent px-3 py-0 text-[10px] font-space-grotesk font-bold tracking-[1px] text-foreground shadow-none hover:border-muted-foreground [&>svg:last-child]:hidden">
				<SelectValue />
				<ChevronDown className="size-2.5 text-muted-foreground" />
			</SelectTrigger>
			<SelectContent align="start" className="min-w-[120px] rounded-none border-border bg-card">
				<SelectItem value="OWNER" className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-[#F5F5F0] focus:bg-[#2D2D2D] focus:text-[#FFD600]">OWNER</SelectItem>
				<SelectItem value="EDITOR" className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-[#F5F5F0] focus:bg-[#2D2D2D] focus:text-[#FFD600]">EDITOR</SelectItem>
				<SelectItem value="VIEWER" className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-[#F5F5F0] focus:bg-[#2D2D2D] focus:text-[#FFD600]">VIEWER</SelectItem>
			</SelectContent>
		</Select>
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
			className={`flex h-16 items-center px-6 ${!isLast ? "border-b border-border" : ""}`}
		>
			<div className="flex flex-1 items-center gap-3">
				<div
					className={`flex h-9 w-9 items-center justify-center ${member.isOwner ? "bg-primary" : "bg-[#2D2D2D]"}`}
				>
					<span
						className={`font-space-grotesk text-xs font-bold ${member.isOwner ? "text-[#1A1A1A]" : "text-[#6B6B6B]"}`}
					>
						{member.initials}
					</span>
				</div>
				<div className="flex flex-col gap-0.5">
					<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-[#F5F5F0]">
						{member.name}
					</span>
					<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-[#3D3D3D]">
						{member.email}
					</span>
				</div>
			</div>
			
			<div className="flex w-40 items-center justify-start">
				{member.isOwner ? (
					<OwnerBadge role={member.role} />
				) : (
					<RoleSelect role={member.role} userId={member.userId} />
				)}
			</div>

			<span className="w-[140px] font-ibm-plex-mono text-xs tracking-[1px] text-[#6B6B6B]">
				{member.joined}
			</span>
			<div className="flex w-20 items-center justify-center">
				{member.isOwner ? (
					<span className="font-space-grotesk text-sm font-bold text-[#3D3D3D]">
						—
					</span>
				) : (
					<button
						type="button"
						className="text-[#FF6B35] hover:opacity-80"
					>
						<Trash2 className="size-4" />
					</button>
				)}
			</div>
		</div>
	)
}

function LoadingState() {
	return (
		<div className="flex h-16 items-center px-6 border-b border-border">
			<div className="flex flex-1 items-center gap-3">
				<div className="h-9 w-9 animate-pulse rounded bg-[#2D2D2D]" />
				<div className="flex flex-col gap-0.5">
					<div className="h-4 w-24 animate-pulse rounded bg-[#2D2D2D]" />
					<div className="h-3 w-32 animate-pulse rounded bg-[#2D2D2D]" />
				</div>
			</div>
			<div className="h-6 w-24 animate-pulse rounded bg-[#2D2D2D]" />
			<div className="ml-6 h-4 w-24 animate-pulse rounded bg-[#2D2D2D]" />
		</div>
	)
}

export function MembersTable({ members: membersData, isLoading, error }: MembersTableProps) {
	const members: Member[] = membersData.map(mapApiMember)

	return (
		<div className="w-full border border-border">
			<div className="flex h-14 items-center justify-between bg-card px-6 border-b border-border">
				<div className="flex items-center gap-3">
					<div className="h-5 w-1 bg-primary" />
					<span className="font-space-grotesk text-sm font-bold tracking-[1px] text-foreground">
						ACTIVE MEMBERS
					</span>
					<span className="font-ibm-plex-mono text-xs text-muted-foreground">
						{members.length}
					</span>
				</div>
			</div>
			<div className="flex h-10 items-center px-6 border-b border-border">
				<span className="flex-1 font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					MEMBER
				</span>
				<span className="w-40 font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					ROLE
				</span>
				<span className="w-[140px] font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					JOINED
				</span>
				<span className="w-20 text-center font-space-grotesk text-[11px] font-bold tracking-[1px] text-muted-foreground">
					ACTIONS
				</span>
			</div>
			{error && (
				<div className="flex h-16 items-center justify-center px-6">
					<span className="font-ibm-plex-mono text-xs text-[#FF6B35]">
						{error}
					</span>
				</div>
			)}
			{isLoading ? (
				<>
					<LoadingState />
					<LoadingState />
					<LoadingState />
				</>
			) : (
				members.map((member, index) => (
					<MemberRow
						key={member.userId}
						member={member}
						isLast={index === members.length - 1}
					/>
				))
			)}
		</div>
	)
}
