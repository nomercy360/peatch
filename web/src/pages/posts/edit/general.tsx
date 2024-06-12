import { FormLayout } from '~/components/edit/layout'
import { editPost, editPostId, setEditPost, setEditUser } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import {
	createEffect,
	createSignal,
	Match,
	onCleanup,
	onMount,
	Switch,
} from 'solid-js'
import { useNavigate, useParams } from '@solidjs/router'
import { Link } from '~/components/Link'
import {
	createPost,
	fetchPresignedUrl,
	updatePost,
	uploadToS3,
} from '~/lib/api'
import { queryClient } from '~/App'

export const [imgFile, setImgFile] = createSignal<File | null>(null)
export const [previewUrl, setPreviewUrl] = createSignal('')

export default function GeneralInfo() {
	const mainButton = useMainButton()

	const navigate = useNavigate()

	const idPath = useParams().id ? '/' + useParams().id : ''

	const create = async () => {
		const created = await createPost(editPost)
		navigate('/posts/' + created.id)
	}

	const edit = async () => {
		await updatePost(editPostId(), editPost)
		await queryClient.invalidateQueries({
			queryKey: ['posts', String(editPostId())],
		})
		navigate('/posts/' + editPostId())
	}

	const createOrEditPost = async () => {
		if (imgFile() && imgFile() !== null) {
			mainButton.disable('Save').showProgress(true)
			try {
				const { path, url } = await fetchPresignedUrl(imgFile()!.name)
				await uploadToS3(
					url,
					imgFile()!,
					e => {},
					() => {},
				)
				setEditPost('image_url', path)
			} catch (e) {
				console.error(e)
			} finally {
				mainButton.enable('Save').showProgress(false)
				setImgFile(null)
				setPreviewUrl('')
			}
		}

		if (editPostId()) {
			await edit()
		} else {
			await create()
		}
	}

	onMount(() => {
		mainButton.enable('Save')
		mainButton.onClick(createOrEditPost)
	})

	createEffect(() => {
		if (editPost.title && editPost.description) {
			mainButton.enable('Save')
		} else {
			mainButton.disable('Save')
		}
	})

	onCleanup(() => {
		mainButton.offClick(createOrEditPost)
	})

	const handleFileChange = (event: any) => {
		const file = event.target.files[0]
		if (file) {
			const maxSize = 1024 * 1024 * 5 // 7MB

			if (file.size > maxSize) {
				window.Telegram.WebApp.showAlert('Try to select a smaller file')
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

	const resolveImage = () => {
		return previewUrl() || editPost.image_url || null
	}

	return (
		<FormLayout
			title="Whats happening?!"
			description="You can tell about your project or share your thoughts here."
			screen={1}
			totalScreens={1}
		>
			<div class="mt-5 flex w-full flex-col items-center justify-start gap-3">
				<input
					maxLength={70}
					class="h-10 w-full rounded-lg bg-main px-2.5 text-main placeholder:text-hint"
					placeholder="Title"
					value={editPost.title}
					onInput={e => setEditPost('title', e.currentTarget.value)}
				/>
				<textarea
					class="size-full h-24 resize-none rounded-lg bg-main p-2.5 text-main placeholder:text-hint"
					placeholder="Description"
					value={editPost.description}
					onInput={e => setEditPost('description', e.currentTarget.value)}
					autocomplete="off"
					autocapitalize="off"
					spellcheck={false}
					maxLength={200}
				/>
				<Switch>
					<Match when={resolveImage()}>
						<div class="relative flex aspect-video w-full flex-col items-center justify-center gap-2">
							<img
								src={resolveImage()!}
								alt="Uploaded image preview"
								class="aspect-[4/3] rounded-xl object-cover"
							/>
							<button
								class="absolute right-2.5 top-2.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-main"
								onClick={() => {
									setPreviewUrl('')
									setImgFile(null)
									setEditPost('image_url', null)
								}}
							>
								<span class="material-symbols-rounded text-[20px] text-button">
									close
								</span>
							</button>
						</div>
					</Match>
					<Match when={!previewUrl()}>
						<label class="flex h-10 w-full cursor-pointer items-center justify-between rounded-lg bg-main px-2.5 text-hint">
							<input
								class="sr-only h-10 w-40 rounded-lg bg-main px-2.5 text-main placeholder:text-hint"
								type="file"
								accept="image/*"
								onChange={handleFileChange}
							/>
							+ Add Photo
						</label>
					</Match>
				</Switch>
				<Link
					href={`/posts/edit${idPath}/location`}
					state={{ back: true }}
					class="flex h-10 w-full cursor-pointer flex-row items-center justify-between rounded-lg px-2.5 text-hint"
				>
					<Switch>
						<Match when={editPost.city && editPost.country}>
							<span>
								{editPost.city}, {editPost.country}
							</span>
							<span class="material-symbols-rounded text-[20px] text-hint">
								edit
							</span>
						</Match>
						<Match when={!editPost.city}>
							Add Location
							<span class="material-symbols-rounded text-hint">
								chevron_right
							</span>
						</Match>
					</Switch>
				</Link>
			</div>
		</FormLayout>
	)
}

const CheckBoxInput = (props: {
	text: string
	checked: boolean
	setChecked: (value: boolean) => void
}) => {
	return (
		<label class="group flex h-10 w-full cursor-pointer items-center justify-between">
			<p class="text-secondary">{props.text}</p>
			<input
				type="checkbox"
				class="sr-only"
				checked={props.checked}
				onChange={() => props.setChecked(!props.checked)}
			/>
			<span class="flex size-7 items-center justify-center rounded-lg border">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="3"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="size-5 scale-0 text-accent opacity-0 group-has-[:checked]:scale-100 group-has-[:checked]:opacity-100"
				>
					<path d="M5 12l5 5l10 -10" />
				</svg>
			</span>
		</label>
	)
}
