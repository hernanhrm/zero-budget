import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@workspace/ui/components/select"
import { ChevronDown } from "lucide-react"
import { useState } from "react"

interface RoleSelectProps {
	role: "OWNER" | "ADMIN" | "MEMBER"
	userId: string
}

export function RoleSelect({ role, userId: _userId }: RoleSelectProps) {
	const [value, setValue] = useState(role)

	const handleValueChange = (newValue: string) => {
		setValue(newValue as "OWNER" | "ADMIN" | "MEMBER")
	}

	return (
		<Select value={value} onValueChange={handleValueChange}>
			<SelectTrigger className="flex h-6 w-fit items-center justify-between gap-1.5 rounded-none border border-border bg-transparent px-3 py-0 text-[10px] font-space-grotesk font-bold tracking-[1px] text-foreground shadow-none hover:border-muted-foreground [&>svg:last-child]:hidden">
				<SelectValue />
				<ChevronDown className="size-2.5 text-muted-foreground" />
			</SelectTrigger>
			<SelectContent
				align="start"
				className="min-w-[120px] rounded-none border-border bg-card"
			>
				<SelectItem
					value="ADMIN"
					className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-foreground focus:bg-skeleton focus:text-accent"
				>
					ADMIN
				</SelectItem>
				<SelectItem
					value="MEMBER"
					className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-foreground focus:bg-skeleton focus:text-accent"
				>
					MEMBER
				</SelectItem>
			</SelectContent>
		</Select>
	)
}
