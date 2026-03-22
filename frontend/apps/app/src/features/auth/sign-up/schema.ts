import { z } from "zod"

export const signUpSchema = z
	.object({
		firstName: z.string().min(1, "First name is required"),
		lastName: z.string().min(1, "Last name is required"),
		email: z
			.string()
			.min(1, "Email is required")
			.email("Invalid email address"),
		password: z.string().min(8, "Password must be at least 8 characters"),
		confirmPassword: z.string().min(1, "Please confirm your password"),
		terms: z.literal(true, {
			errorMap: () => ({ message: "You must accept the terms" }),
		}),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: "Passwords don't match",
		path: ["confirmPassword"],
	})

export type SignUpFormValues = z.infer<typeof signUpSchema>
