import { ProgressBar } from '~/components/edit/progress'

export function FormLayout(props: {
  children: any
  title: string
  description: string
  screen: number
  totalScreens: number
}) {
  return (
    <div class="flex h-fit min-h-screen flex-col items-center justify-start p-3.5">
      <ProgressBar screen={props.screen} totalScreens={props.totalScreens} />
      <p class="mt-2 text-center text-3xl font-semibold">{props.title}</p>
      <p class="mt-1 text-center text-sm text-secondary-foreground">
        {props.description}
      </p>
      {props.children}
    </div>
  )
}
