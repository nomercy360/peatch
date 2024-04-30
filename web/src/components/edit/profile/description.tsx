import { FormLayout } from '../layout';
import TextArea from '../../TextArea';

export default function FillDescription(props: {
  description: string;
  setDescription: (description: string) => void;
}) {
  const resizer = (e: any) => {
    e.target.style.height = 'auto';
    e.target.style.height = e.target.scrollHeight + 'px';

    props.setDescription(e.target.value);
  };

  return (
    <FormLayout
      title="Introduce yourself"
      description="Tell others about your backround, achievments and goals"
    >
      <TextArea value={props.description} setValue={props.setDescription} />
    </FormLayout>
  );
}
