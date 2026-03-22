interface OwnerBadgeProps {
	role: string
}

export function OwnerBadge({ role }: OwnerBadgeProps) {
	return (
		<div className="flex h-6 items-center justify-start gap-1.5 pr-3">
			<div className="h-1.5 w-1.5 rounded-full bg-primary" />
			<span className="font-space-grotesk text-[10px] font-bold tracking-[1px] text-primary">
				{role}
			</span>
		</div>
	)
}
