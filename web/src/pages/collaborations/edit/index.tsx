import { RouteSectionProps, useParams } from '@solidjs/router';
import { fetchCollaboration } from '~/api';
import { setEditCollaboration } from '~/store';
import { createEffect, createResource, Show } from 'solid-js';

export default function EditCollaboration(props: RouteSectionProps) {
  const params = useParams();

  console.log('PARAMS:', params.id)

  if (!params.id) {
    setEditCollaboration({});

    return (
      <div>{props.children}</div>
    );
  } else {
    const [collaboration, _] = createResource(async () => {
      return await fetchCollaboration(Number(params.id));
    });

    createEffect(() => {
      setEditCollaboration(collaboration());
    })

    return (
      <Show when={!collaboration.loading}>
        <div>{props.children}</div>
      </Show>
    );
  }
}