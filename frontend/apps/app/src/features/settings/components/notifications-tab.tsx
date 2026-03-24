import { Separator } from "@workspace/ui/components/separator"
import { Switch } from "@workspace/ui/components/switch"
import { useState } from "react"

interface NotificationSetting {
	id: string
	title: string
	description: string
	defaultEnabled: boolean
}

const notificationSettings: NotificationSetting[] = [
	{
		id: "budget-alerts",
		title: "BUDGET ALERTS",
		description: "GET NOTIFIED WHEN YOU EXCEED BUDGET LIMITS",
		defaultEnabled: true,
	},
	{
		id: "transaction-notifications",
		title: "TRANSACTION NOTIFICATIONS",
		description: "RECEIVE ALERTS FOR NEW TRANSACTIONS",
		defaultEnabled: true,
	},
	{
		id: "monthly-report",
		title: "MONTHLY REPORT",
		description: "GET A MONTHLY SUMMARY OF YOUR FINANCES",
		defaultEnabled: false,
	},
]

export function NotificationsTab() {
	const [settings, setSettings] = useState(() =>
		Object.fromEntries(
			notificationSettings.map((s) => [s.id, s.defaultEnabled]),
		),
	)

	const handleToggle = (id: string, checked: boolean) => {
		setSettings((prev) => ({ ...prev, [id]: checked }))
	}

	return (
		<div className="flex max-w-2xl flex-col">
			{notificationSettings.map((setting, index) => (
				<div key={setting.id}>
					<div className="flex h-16 items-center justify-between">
						<div className="flex flex-col gap-1">
							<span className="font-space-grotesk text-[13px] font-bold tracking-[1px] text-foreground">
								{setting.title}
							</span>
							<span className="font-ibm-plex-mono text-[10px] tracking-[1px] text-muted-foreground">
								{setting.description}
							</span>
						</div>
						<Switch
							checked={settings[setting.id] ?? false}
							onCheckedChange={(checked) =>
								handleToggle(setting.id, checked)
							}
						/>
					</div>
					{index < notificationSettings.length - 1 && <Separator />}
				</div>
			))}
		</div>
	)
}
