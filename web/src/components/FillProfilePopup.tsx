export default function FillProfilePopup({ onClose }: { onClose: () => void }) {
  return (
    <div class="fixed bottom-0 right-0 flex w-full flex-col items-center justify-center rounded-t-3xl bg-peatch-blue px-4 py-3">
      <button
        class="flex w-full flex-row items-start justify-between gap-4 text-start"
        onClick={onClose}
      >
        <span class="text-[24px] font-bold leading-tight text-white">
          Set up your profile to collaborate with others.
        </span>
        <span class="flex size-6 items-center justify-center rounded-full bg-white/10">
          <span class="material-symbols-rounded text-[24px] text-white">
            close
          </span>
        </span>
      </button>
      <p class="mt-2 text-xl text-white">
        It only takes 5 minutes, but it can significantly improve your
        networking. According to our data, every third user finds someone to
        collaborate with within the first three days.
      </p>
      <a
        class="mt-4 flex h-12 w-full items-center justify-center rounded-2xl bg-white text-center text-peatch-blue"
        href="/users/[edit]"
      >
        Set up profile
      </a>
    </div>
  );
}
