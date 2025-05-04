import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { createSignal, Match, onCleanup, onMount, Switch } from 'solid-js'
import { editUser, store } from '~/store'
import {
	API_BASE_URL,
	updateUser,
	uploadUserAvatar,
} from '~/lib/api'
import { usePopup } from '~/lib/usePopup'
import { useTranslations } from '~/lib/locale-context'
import { queryClient } from '~/App'

export default function ImageUpload() {
	const mainButton = useMainButton()
	const { t } = useTranslations()
	const [imgFile, setImgFile] = createSignal<File | null>(null)

	const navigate = useNavigate()
	const { showAlert } = usePopup()

	const imgFromCDN = store.user.avatar_url
		? `https://assets.peatch.io/cdn-cgi/image/width=400/${store.user.avatar_url}`
		: ''

	const [previewUrl, setPreviewUrl] = createSignal(imgFromCDN || '')

	const handleFileChange = (event: any) => {
		const file = event.target.files[0]
		if (file) {
			const maxSize = 1024 * 1024 * 5 // 7MB

			if (file.size > maxSize) {
				showAlert('Try to select a smaller file')
				return
			}

			setImgFile(file)
			setPreviewUrl('')

			const reader = new FileReader()
			reader.onload = e => {
				setPreviewUrl(e.target?.result as string)
			}
			reader.readAsDataURL(file)
		}
	}

	const generateRandomAvatar = () => {
		const url = `${API_BASE_URL}/avatar`

		const resp = fetch(url)

		resp.then(response => {
			response.blob().then(blob => {
				const file = new File([blob], 'avatar.svg', {
					type: 'image/svg+xml',
				})
				setImgFile(file)
				setPreviewUrl('')
				setPreviewUrl(URL.createObjectURL(file))
			})
		})
	}

	const saveUser = async () => {
		try {
			const file = imgFile()
			if (file) {
				mainButton.showProgress(true)
				await uploadUserAvatar(file)
			}

			await updateUser({
				location_id: editUser.location.id,
				first_name: editUser.first_name,
				last_name: editUser.last_name,
				title: editUser.title,
				description: editUser.description,
				badge_ids: editUser.badge_ids,
				opportunity_ids: editUser.opportunity_ids,
			})
		} catch (e) {
			console.error(e)
		} finally {
			mainButton.hideProgress()
		}

		queryClient.invalidateQueries({ queryKey: ['profiles', store.user.id] })
		window.Telegram.WebApp.disableClosingConfirmation()
		navigate(`/users/${store.user.id}`)
	}

	onMount(() => {
		mainButton.onClick(saveUser)
		mainButton.enable(t('common.buttons.save'))
	})

	onCleanup(() => {
		mainButton.offClick(saveUser)
	})

	return (
		<FormLayout
			title={t('pages.users.edit.image.title')}
			description={t('pages.users.edit.image.description')}
			screen={5}
			totalScreens={6}
		>
			<div class="mt-5 flex h-full items-center justify-center">
				<div class="flex flex-col items-center justify-center gap-2">
					<Switch>
						<Match when={previewUrl()}>
							<ImageBox imgURL={previewUrl()} onFileChange={handleFileChange} />
						</Match>
						<Match when={!previewUrl()}>
							<UploadBox onFileChange={handleFileChange} />
						</Match>
					</Switch>
					<button class="text-link h-10" onClick={generateRandomAvatar}>
						{t('common.buttons.generateRandomAvatar')}
					</button>
				</div>
			</div>
		</FormLayout>
	)
}

type ImageBoxProps = {
	imgURL: string
	onFileChange: (event: any) => void
}

function ImageBox(props: ImageBoxProps) {
	return (
		<div class="mt-5 flex h-full items-center justify-center">
			<div class="relative flex size-56 flex-col items-center justify-center gap-2">
				<img
					src={props.imgURL}
					alt="Uploaded image preview"
					class="size-56 rounded-xl object-cover"
				/>
				<input
					class="absolute size-full cursor-pointer rounded-xl opacity-0"
					type="file"
					accept="image/*"
					onChange={e => props.onFileChange(e)}
				/>
			</div>
		</div>
	)
}

type UploadBoxProps = {
	onFileChange: (event: any) => void
}

function UploadBox(props: UploadBoxProps) {
	return (
		<>
			<div class="relative flex size-56 flex-col items-center justify-center rounded-xl">
				<input
					class="absolute size-full opacity-0"
					type="file"
					accept="image/*"
					onChange={e => props.onFileChange(e)}
				/>
				<span class="material-symbols-rounded pointer-events-none z-10 text-[45px] text-secondary">
					camera_alt
				</span>
			</div>
		</>
	)
}
