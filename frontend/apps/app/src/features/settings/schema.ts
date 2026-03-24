import { z } from "zod"

export const profileSchema = z.object({
	firstName: z.string().min(1, "First name is required"),
	lastName: z.string().min(1, "Last name is required"),
	timezone: z.string().min(1, "Timezone is required"),
})

export type ProfileFormValues = z.infer<typeof profileSchema>

export const changePasswordSchema = z
	.object({
		currentPassword: z.string().min(1, "Current password is required"),
		newPassword: z
			.string()
			.min(8, "Password must be at least 8 characters"),
		confirmPassword: z.string().min(1, "Please confirm your password"),
	})
	.refine((data) => data.newPassword === data.confirmPassword, {
		message: "Passwords do not match",
		path: ["confirmPassword"],
	})

export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>

export const deleteAccountSchema = z.object({
	confirmation: z
		.string()
		.refine((val) => val === "DELETE", "Type DELETE to confirm"),
})

export type DeleteAccountFormValues = z.infer<typeof deleteAccountSchema>
