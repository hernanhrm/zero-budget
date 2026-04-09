import { ModulePageHeader } from "@workspace/ui/components/module-page-header"
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from "@workspace/ui/components/tabs"
import { getActiveOrganizationId } from "#/lib/session-org"
import { Route } from "#/routes/__root"
import { CurrencyTab } from "./components/currency-tab"
import { DangerZoneTab } from "./components/danger-zone-tab"
import { NotificationsTab } from "./components/notifications-tab"
import { ProfileTab } from "./components/profile-tab"
import { SecurityTab } from "./components/security-tab"

export function SettingsPage() {
	const { session } = Route.useRouteContext()
	const organizationId = getActiveOrganizationId(session)

	return (
		<div className="flex h-full flex-col overflow-auto p-10">
			<Tabs defaultValue="profile" className="flex-col gap-0">
				<ModulePageHeader
					title="SETTINGS"
					description="CONFIGURE YOUR ACCOUNT AND PREFERENCES"
				/>

				<TabsList
					variant="line"
					className="mt-8 h-11 w-full justify-start gap-0 border-b border-border p-0"
				>
					<TabsTrigger
						value="profile"
						className="h-full !flex-initial !py-0 px-5 font-space-grotesk text-xs font-bold tracking-[1px] after:bg-primary"
					>
						PROFILE
					</TabsTrigger>
					<TabsTrigger
						value="security"
						className="h-full !flex-initial !py-0 px-5 font-space-grotesk text-xs font-bold tracking-[1px] after:bg-primary"
					>
						SECURITY
					</TabsTrigger>
					<TabsTrigger
						value="currency"
						className="h-full !flex-initial !py-0 px-5 font-space-grotesk text-xs font-bold tracking-[1px] after:bg-primary"
					>
						CURRENCY
					</TabsTrigger>
					<TabsTrigger
						value="notifications"
						className="h-full !flex-initial !py-0 px-5 font-space-grotesk text-xs font-bold tracking-[1px] after:bg-primary"
					>
						NOTIFICATIONS
					</TabsTrigger>
					<TabsTrigger
						value="danger-zone"
						className="h-full !flex-initial !py-0 px-5 font-space-grotesk text-xs font-bold tracking-[1px] text-[#FF6B35] after:bg-[#FF6B35] data-active:text-[#FF6B35]"
					>
						DANGER ZONE
					</TabsTrigger>
				</TabsList>

				<TabsContent value="profile" className="mt-8">
					<ProfileTab />
				</TabsContent>
				<TabsContent value="security" className="mt-8">
					<SecurityTab />
				</TabsContent>
				<TabsContent value="currency" className="mt-8">
					<CurrencyTab organizationId={organizationId} />
				</TabsContent>
				<TabsContent value="notifications" className="mt-8">
					<NotificationsTab />
				</TabsContent>
				<TabsContent value="danger-zone" className="mt-8">
					<DangerZoneTab />
				</TabsContent>
			</Tabs>
		</div>
	)
}
