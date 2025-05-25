import { Link } from '~/components/link'
import { useTranslations } from '~/lib/locale-context'

export default function FillProfilePopup(props: { onClose: () => void }) {
  const { t } = useTranslations()

  return (
    <div class="relative rounded-lg bg-secondary p-2 text-center">
      <button
        class="absolute right-2 top-2 size-5 rounded-full bg-background"
        onClick={() => props.onClose()}
      >
        <span class="material-symbols-rounded text-secondary-foreground">
          close
        </span>
      </button>
      <div class="text-green flex items-center justify-center gap-1 text-xl font-bold">
        <span class="material-symbols-rounded text-blue-400">people</span>
        {t('pages.users.fillProfilePopup.title')}
      </div>
      <p class="mt-1 text-sm text-secondary-foreground">
        {t('pages.users.fillProfilePopup.description')}
      </p>
      <Link
        class="mt-2 flex h-8 w-full items-center justify-center rounded-lg bg-primary text-sm font-semibold"
        href="/users/edit"
      >
        {t('pages.users.fillProfilePopup.action')}
      </Link>
    </div>
  )
}
