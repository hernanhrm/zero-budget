import { cn } from "@workspace/ui/lib/utils"

const Table = ({ className, ...props }: React.ComponentProps<"table">) => (
	<table
		data-slot="table"
		className={cn("w-full caption-bottom text-sm", className)}
		{...props}
	/>
)
Table.displayName = "Table"

const TableHeader = ({
	className,
	...props
}: React.ComponentProps<"thead">) => (
	<thead data-slot="table-header" className={cn("[&_tr]:border-b", className)} {...props} />
)
TableHeader.displayName = "TableHeader"

const TableBody = ({ className, ...props }: React.ComponentProps<"tbody">) => (
	<tbody
		data-slot="table-body"
		className={cn("[&_tr:last-child]:border-0", className)}
		{...props}
	/>
)
TableBody.displayName = "TableBody"

const TableFooter = ({ className, ...props }: React.ComponentProps<"tfoot">) => (
	<tfoot
		data-slot="table-footer"
		className={cn(
			"border-t bg-muted/50 font-medium [&_tr]:last:border-b-0",
			className,
		)}
		{...props}
	/>
)
TableFooter.displayName = "TableFooter"

const TableRow = ({ className, ...props }: React.ComponentProps<"tr">) => (
	<tr
		data-slot="table-row"
		className={cn(
			"border-b border-border transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted",
			className,
		)}
		{...props}
	/>
)
TableRow.displayName = "TableRow"

const TableHead = ({ className, ...props }: React.ComponentProps<"th">) => (
	<th
		data-slot="table-head"
		className={cn(
			"h-10 px-2 text-left align-middle font-medium text-muted-foreground [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]",
			className,
		)}
		{...props}
	/>
)
TableHead.displayName = "TableHead"

const TableCell = ({ className, ...props }: React.ComponentProps<"td">) => (
	<td
		data-slot="table-cell"
		className={cn(
			"p-2 align-middle [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]",
			className,
		)}
		{...props}
	/>
)
TableCell.displayName = "TableCell"

const TableCaption = ({ className, ...props }: React.ComponentProps<"caption">) => (
	<caption
		data-slot="table-caption"
		className={cn("mt-4 text-sm text-muted-foreground", className)}
		{...props}
	/>
)
TableCaption.displayName = "TableCaption"

export {
	Table,
	TableHeader,
	TableBody,
	TableFooter,
	TableHead,
	TableRow,
	TableCell,
	TableCaption,
}
