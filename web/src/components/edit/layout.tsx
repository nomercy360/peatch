import { ProgressBar } from '~/components/edit/progress';

export function FormLayout(props: {
  children: any;
  title: string;
  description: string;
  screen: number;
  totalScreens: number;
}) {
  return (
    <div class="flex h-screen flex-col items-center justify-start p-3.5">
      <ProgressBar screen={props.screen} totalScreens={props.totalScreens} />
      <p class="mt-2 text-3xl">{props.title}</p>
      <p class="mt-1 text-sm text-gray">{props.description}</p>
      {props.children}
    </div>
  );
}
