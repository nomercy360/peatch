import { createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { Button } from '~/components/ui/button'
import { TextField, TextFieldInput } from '~/components/ui/text-field'
import { login } from '~/lib/api'

export default function LoginPage() {
	const [token, setToken] = createSignal('')
	const [isLoading, setIsLoading] = createSignal(false)
	const [error, setError] = createSignal('')
	const navigate = useNavigate()

	const handleSubmit = async (e: Event) => {
		e.preventDefault()
		setIsLoading(true)
		setError('')

		try {
			const isValid = await login(token())
			if (isValid) {
				navigate('/users')
			} else {
				setError('Invalid token')
			}
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Login failed')
		} finally {
			setIsLoading(false)
		}
	}

	return (
		<div class="flex min-h-screen items-center justify-center bg-secondary">
			<div class="w-full max-w-md">
				<div class="rounded-lg shadow-md p-8">
					<h1 class="text-2xl font-bold text-center mb-6">Admin Login</h1>

					<form onSubmit={handleSubmit} class="space-y-4">
						<div>
							<label class="block text-sm font-medium text-muted-foreground mb-2">
								API Token
							</label>
							<TextField value={token()} onChange={setToken}>
								<TextFieldInput
									type="password"
									placeholder="Enter your API token"
									required
								/>
							</TextField>
						</div>

						{error() && (
							<div class="text-sm text-destructive-foreground">{error()}</div>
						)}

						<Button
							type="submit"
							class="w-full"
							disabled={isLoading()}
						>
							{isLoading() ? 'Logging in...' : 'Login'}
						</Button>
					</form>
				</div>
			</div>
		</div>
	)
}
