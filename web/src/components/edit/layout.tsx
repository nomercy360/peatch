export function FormLayout(props: {
  children: any;
  title: string;
  description: string;
}) {
  return (
    <>
      <p class="mt-2 text-3xl">{props.title}</p>
      <p class="mt-1 text-sm text-gray">{props.description}</p>
      {props.children}
    </>
  );
}