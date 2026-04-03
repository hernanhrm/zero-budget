import { FieldLabel } from "@workspace/ui/components/field"
import { cn } from "@workspace/ui/lib/utils"
import { cva, type VariantProps } from "class-variance-authority"
import type { ComponentProps } from "react"

const formFieldLabelVariants = cva(
	"font-space-grotesk font-bold tracking-[1px] text-muted-foreground",
	{
		variants: {
			size: {
				default: "text-[11px]",
				sm: "text-[10px]",
			},
		},
		defaultVariants: {
			size: "default",
		},
	},
)

export type FormFieldLabelProps = ComponentProps<typeof FieldLabel> &
	VariantProps<typeof formFieldLabelVariants>

export function FormFieldLabel({
	className,
	size,
	...props
}: FormFieldLabelProps) {
	return (
		<FieldLabel
			className={cn(formFieldLabelVariants({ size }), className)}
			{...props}
		/>
	)
}
