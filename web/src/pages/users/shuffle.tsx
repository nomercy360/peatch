import { createInfiniteQuery, keepPreviousData } from '@tanstack/solid-query';
import { CDN_URL, fetchMatchingProfiles } from '~/lib/api';
import { createEffect, createSignal, For, Show } from 'solid-js';
import Badge from '~/components/Badge';
import { useMainButton } from '~/lib/useMainButton';

export default function ShufflePage() {
  const mainButton = useMainButton();

  const profileQ = createInfiniteQuery(() => ({
    queryKey: ['matchingProfiles'],
    queryFn: ({ pageParam }) => fetchMatchingProfiles(pageParam),
    getNextPageParam: (lastPage, allPages) => {
      if (lastPage.length === 0) return null;
      return allPages.length + 1;
    },
    initialPageParam: 1,
    placeholderData: keepPreviousData,
    refetchOnMount: true,
  }));

  const [currentProfile, setCurrentProfile] = createSignal(0);

  // const handleScrollToNext = () => {
  // 	window.scrollTo({
  // 		top: window.innerHeight - 80,
  // 		behavior: 'smooth',
  // 	})
  // }

  const [profileRefs, setProfileRefs] = createSignal<HTMLElement[]>([]);

  const handleNextProfile = async () => {
    // Scroll to the next profile
    const nextProfileElement = profileRefs()[currentProfile()];
    if (nextProfileElement) {
      nextProfileElement.scrollIntoView({ behavior: 'smooth' });
    }

    const nextProfile = currentProfile() + 1;
    setCurrentProfile(nextProfile);

    if (nextProfile % 5 === 4) {
      await profileQ.fetchNextPage();
    }
  };

  // one page = 5 profiles. If we on a 4th profile of the page, we need to load next page.
  // we need to implment scroll to next also. Also now we override profiles on each page change. need to store them all

  createEffect(() => {
    if (currentProfile() > 0) {
      mainButton.enable('Message in Telegram').onClick(() => {
        console.log('Message in Telegram');
      });
    }
  });

  return (
    <div class="flex min-h-screen flex-col items-center justify-center bg-secondary">
      <div class="flex h-[calc(100vh-80px)] flex-col items-center justify-center px-10 text-center">
        <Sparkles />
        <p class="mt-2 text-3xl text-main">
          A feed with collaborations and people, selected just for you
        </p>
        <p class="mt-2 text-secondary">
          Once we will find something or someone, that might be interesting for
          you, we will add it here.
        </p>
        <button
          class="mt-36 h-10 text-link"
          onClick={() => handleNextProfile()}
        >
          Start scrolling 􀆈
        </button>
      </div>
      <div class="px-2">
        <For each={profileQ.data?.pages}>
          {group => (
            <For each={group}>
              {profile => (
                <div
                  ref={el => {
                    setProfileRefs([...profileRefs(), el]);
                  }}
                  class="mb-4 flex min-h-screen flex-col items-center justify-start rounded-xl bg-main pb-4 text-center"
                >
                  <div class="relative w-full">
                    <img
                      src={CDN_URL + '/' + profile.avatar_url}
                      alt={profile.username}
                      class="aspect-square w-full shrink-0 rounded-xl object-cover"
                    />
                    <div class="absolute bottom-0 left-0 h-40 w-full bg-gradient-to-t from-main to-transparent" />
                  </div>
                  <p class="text-3xl font-black text-pink">
                    {profile.last_name
                      ? profile.first_name + ' ' + profile.last_name + ':'
                      : profile.first_name + ':'}
                  </p>
                  <p class="text-3xl font-black text-main">{profile.title}</p>
                  <div class="w-full px-4 text-center">
                    <p class="text-hint">{profile.description}</p>
                  </div>
                  <div class="mt-4 flex w-full flex-col items-start justify-start px-4">
                    <Show when={profile.badges && profile?.badges.length > 0}>
                      <p class="text-sm font-semibold text-secondary">
                        Interests in common
                      </p>
                      <div class="mt-2 flex flex-row flex-wrap items-center justify-center gap-2">
                        <For each={profile.badges}>
                          {badge => (
                            <Badge
                              icon={badge.icon!}
                              name={badge.text!}
                              color={badge.color!}
                            />
                          )}
                        </For>
                      </div>
                    </Show>
                    <Show
                      when={
                        profile.opportunities &&
                        profile?.opportunities.length > 0
                      }
                    >
                      <p class="mt-4 text-sm font-semibold text-secondary">
                        You're both open to
                      </p>
                      <div class="mt-2 flex flex-row flex-wrap items-center justify-center gap-2">
                        <For each={profile.opportunities}>
                          {opportunity => (
                            <Badge
                              icon={opportunity.icon!}
                              name={opportunity.text!}
                              color={opportunity.color!}
                            />
                          )}
                        </For>
                      </div>
                    </Show>
                  </div>
                  <button
                    class="mt-4 h-10 text-link"
                    onClick={handleNextProfile}
                  >
                    Next profile
                  </button>
                </div>
              )}
            </For>
          )}
        </For>
      </div>
    </div>
  );
}

const Sparkles = () => {
  return (
    <svg
      width="23"
      height="27"
      viewBox="0 0 23 27"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M13.1172 26.832C12.9688 26.832 12.8398 26.7812 12.7305 26.6797C12.6289 26.5859 12.5664 26.457 12.543 26.293C12.3633 25.0117 12.168 23.9336 11.957 23.0586C11.7539 22.1914 11.4805 21.4766 11.1367 20.9141C10.793 20.3516 10.3398 19.9023 9.77734 19.5664C9.22266 19.2305 8.51172 18.9609 7.64453 18.7578C6.78516 18.5469 5.72266 18.3594 4.45703 18.1953C4.28516 18.1797 4.14844 18.1172 4.04688 18.0078C3.94531 17.8984 3.89453 17.7656 3.89453 17.6094C3.89453 17.4609 3.94531 17.332 4.04688 17.2227C4.14844 17.1133 4.28516 17.0508 4.45703 17.0352C5.72266 16.8945 6.78906 16.7266 7.65625 16.5312C8.52344 16.3359 9.23828 16.0664 9.80078 15.7227C10.3633 15.3789 10.8164 14.9219 11.1602 14.3516C11.5039 13.7812 11.7773 13.0586 11.9805 12.1836C12.1914 11.3008 12.3789 10.2148 12.543 8.92578C12.5664 8.76953 12.6289 8.64453 12.7305 8.55078C12.8398 8.44922 12.9688 8.39844 13.1172 8.39844C13.2734 8.39844 13.4023 8.44922 13.5039 8.55078C13.6055 8.64453 13.6719 8.76953 13.7031 8.92578C13.8672 10.2148 14.0508 11.3008 14.2539 12.1836C14.4648 13.0586 14.7383 13.7812 15.0742 14.3516C15.418 14.9219 15.8711 15.3789 16.4336 15.7227C16.9961 16.0664 17.7109 16.3359 18.5781 16.5312C19.4453 16.7266 20.5156 16.8945 21.7891 17.0352C21.9531 17.0508 22.0859 17.1133 22.1875 17.2227C22.2891 17.332 22.3398 17.4609 22.3398 17.6094C22.3398 17.7656 22.2891 17.8984 22.1875 18.0078C22.0859 18.1172 21.9531 18.1797 21.7891 18.1953C20.5156 18.3359 19.4453 18.5039 18.5781 18.6992C17.7109 18.8945 16.9961 19.1641 16.4336 19.5078C15.8711 19.8438 15.418 20.2969 15.0742 20.8672C14.7383 21.4375 14.4648 22.1641 14.2539 23.0469C14.0508 23.9297 13.8672 25.0117 13.7031 26.293C13.6719 26.457 13.6055 26.5859 13.5039 26.6797C13.4023 26.7812 13.2734 26.832 13.1172 26.832ZM5.17188 13.8711C4.9375 13.8711 4.80469 13.7422 4.77344 13.4844C4.67969 12.7031 4.57812 12.0898 4.46875 11.6445C4.36719 11.1992 4.19922 10.8672 3.96484 10.6484C3.73828 10.4219 3.39453 10.25 2.93359 10.1328C2.48047 10.0156 1.85547 9.89062 1.05859 9.75781C0.792969 9.71875 0.660156 9.58594 0.660156 9.35938C0.660156 9.14062 0.777344 9.01172 1.01172 8.97266C1.81641 8.81641 2.44922 8.67969 2.91016 8.5625C3.37109 8.4375 3.71875 8.26562 3.95312 8.04688C4.1875 7.82812 4.35938 7.50391 4.46875 7.07422C4.57812 6.63672 4.67969 6.02734 4.77344 5.24609C4.80469 4.98828 4.9375 4.85938 5.17188 4.85938C5.40625 4.85938 5.53906 4.98438 5.57031 5.23438C5.67188 6.02344 5.77344 6.64453 5.875 7.09766C5.98438 7.55078 6.15234 7.89453 6.37891 8.12891C6.61328 8.35547 6.96094 8.52344 7.42188 8.63281C7.88281 8.74219 8.51953 8.85547 9.33203 8.97266C9.43359 8.98047 9.51562 9.01953 9.57812 9.08984C9.64844 9.16016 9.68359 9.25 9.68359 9.35938C9.68359 9.57812 9.56641 9.71094 9.33203 9.75781C8.51953 9.91406 7.88281 10.0547 7.42188 10.1797C6.96875 10.2969 6.625 10.4688 6.39062 10.6953C6.16406 10.9141 5.99609 11.2422 5.88672 11.6797C5.77734 12.1172 5.67188 12.7266 5.57031 13.5078C5.55469 13.6094 5.51172 13.6953 5.44141 13.7656C5.37109 13.8359 5.28125 13.8711 5.17188 13.8711ZM10.8438 5.80859C10.6953 5.80859 10.6094 5.73047 10.5859 5.57422C10.4922 5.09766 10.4102 4.72266 10.3398 4.44922C10.2773 4.17578 10.1758 3.96875 10.0352 3.82812C9.90234 3.67969 9.69531 3.5625 9.41406 3.47656C9.13281 3.39063 8.73828 3.30469 8.23047 3.21875C8.07422 3.1875 7.99609 3.09766 7.99609 2.94922C7.99609 2.80859 8.07422 2.72266 8.23047 2.69141C8.73828 2.59766 9.13281 2.51172 9.41406 2.43359C9.69531 2.34766 9.90234 2.23438 10.0352 2.09375C10.1758 1.94531 10.2773 1.73437 10.3398 1.46094C10.4102 1.1875 10.4922 0.8125 10.5859 0.335938C10.6094 0.179688 10.6953 0.101562 10.8438 0.101562C10.9844 0.101562 11.0703 0.179688 11.1016 0.335938C11.1875 0.8125 11.2656 1.1875 11.3359 1.46094C11.4062 1.73437 11.5078 1.94531 11.6406 2.09375C11.7812 2.23438 11.9922 2.34766 12.2734 2.43359C12.5547 2.51172 12.9492 2.59766 13.457 2.69141C13.6133 2.72266 13.6914 2.80859 13.6914 2.94922C13.6914 3.09766 13.6133 3.1875 13.457 3.21875C12.9492 3.30469 12.5547 3.39063 12.2734 3.47656C11.9922 3.5625 11.7812 3.67969 11.6406 3.82812C11.5078 3.96875 11.4062 4.17578 11.3359 4.44922C11.2656 4.72266 11.1875 5.09766 11.1016 5.57422C11.0703 5.73047 10.9844 5.80859 10.8438 5.80859Z"
        fill="#FDAC2C"
      />
    </svg>
  );
};
