import { useTranslations } from '~/lib/locale-context'

export default function NotFound() {
  const { t } = useTranslations()

  return (
    <section class="p-8">
      <h1 class="text-lg font-bold">{t('pages.notFound.title')}</h1>
      <p class="mt-4">It's gone 😞</p>
    </section>
  )
}
