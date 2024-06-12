import { RouteSectionProps, useNavigate, useParams } from '@solidjs/router'
import { fetchPost } from '~/lib/api'
import { setEditPost, setEditPostId } from '~/store'
import { createEffect, createResource, Show } from 'solid-js'

export default function EditCollaboration(props: RouteSectionProps) {
	const params = useParams()

	if (!params.id) {
		setEditPost({
			city: '',
			country: '',
			country_code: '',
			description: '',
			image_url: '',
		})

		const navigate = useNavigate()
		navigate('/posts/edit')
		return <div>{props.children}</div>
	} else {
		const [post, _] = createResource(async () => {
			return await fetchPost(Number(params.id))
		})

		createEffect(() => {
			if (!post.loading) {
				setEditPost(post())
				setEditPostId(post().id)
			}
		})

		return (
			<Show when={!post.loading}>
				<div>{props.children}</div>
			</Show>
		)
	}
}
