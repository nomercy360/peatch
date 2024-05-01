import { editCollaboration, setEditCollaboration } from '~/store';
import { RouteSectionProps, useNavigate } from '@solidjs/router';

export default function EditCollaboration(props: RouteSectionProps) {
  const navigate = useNavigate();

  setEditCollaboration({});

  if (!editCollaboration.title || !editCollaboration.description) {
    navigate('/collaborations/edit');
  }

  return <div>{props.children}</div>;
}
