import ActionDonePopup from '~/components/ActionDonePopup';
import { useNavigation } from '~/hooks/useNavigation';
import { useLocation } from '@solidjs/router';

export default function CollaborationPublished() {
  const { navigateBack } = useNavigation();
  const location = useLocation().state;

  return (
    <ActionDonePopup
      action="Collaboration published!"
      description="We have shared your collaboration with the community"
      callToAction="There are 12 people you might be interested to collaborate with"
    />
  );
}
