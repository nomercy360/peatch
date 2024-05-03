import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal, Match, onCleanup, Switch } from 'solid-js';
import { editUser, setEditUser, store } from '~/store';
import {
  API_BASE_URL,
  CDN_URL,
  fetchPresignedUrl,
  updateUser,
  uploadToS3,
} from '~/api';
import { usePopup } from '~/hooks/usePopup';

export default function SelectBadges() {
  const mainButton = useMainButton();
  const [imgFile, setImgFile] = createSignal<File | null>(null);
  const [_, setImgUploadProgress] = createSignal(0);

  const navigate = useNavigate();
  const { showAlert } = usePopup();

  const imgFromCDN = store.user.avatar_url
    ? CDN_URL + '/' + store.user.avatar_url
    : '';

  const [previewUrl, setPreviewUrl] = createSignal(imgFromCDN || '');

  const handleFileChange = (event: any) => {
    const file = event.target.files[0];
    if (file) {
      const maxSize = 1024 * 1024 * 5; // 7MB

      if (file.size > maxSize) {
        showAlert('Try to select a smaller file');
        return;
      }

      setImgFile(file);
      setPreviewUrl('');

      const reader = new FileReader();
      reader.onload = e => {
        setPreviewUrl(e.target?.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const generateRandomAvatar = () => {
    const url = `${API_BASE_URL}/avatar`;

    const resp = fetch(url);

    resp.then(response => {
      response.blob().then(blob => {
        const file = new File([blob], 'avatar.svg', {
          type: 'image/svg+xml',
        });
        setImgFile(file);
        setPreviewUrl('');
        setPreviewUrl(URL.createObjectURL(file));
      });
    });
  };

  const saveUser = async () => {
    if (imgFile() && imgFile() !== null) {
      mainButton.showProgress(true);
      try {
        const { path, url } = await fetchPresignedUrl(imgFile()!.name);
        await uploadToS3(
          url,
          imgFile()!,
          e => {
            setImgUploadProgress(Math.round((e.loaded / e.total) * 100));
          },
          () => {
            setImgUploadProgress(0);
          },
        );
        setEditUser('avatar_url', path);
      } catch (e) {
        console.error(e);
      }
    }
    await updateUser(editUser);

    mainButton.hideProgress();
    navigate('/users/' + store.user.id + '?refetch=true');
  };

  mainButton.onClick(saveUser);

  createEffect(() => {
    if (editUser.avatar_url || imgFile()) {
      mainButton.enable('Save');
    } else {
      mainButton.disable('Save');
    }
  });

  onCleanup(() => {
    mainButton.offClick(saveUser);
  });

  return (
    <FormLayout
      title="Upload your photo"
      description="Select one with good lighting and minimal background details"
      screen={5}
      totalScreens={6}
    >
      <div class="mt-5 flex h-full items-center justify-center">
        <div class="flex flex-col items-center justify-center gap-2">
          <Switch>
            <Match when={previewUrl()}>
              <ImageBox imgURL={previewUrl()} onFileChange={handleFileChange} />
            </Match>
            <Match when={!previewUrl()}>
              <UploadBox onFileChange={handleFileChange} />
            </Match>
          </Switch>
          <button class="h-10 text-link" onClick={generateRandomAvatar}>
            Generate a random avatar
          </button>
        </div>
      </div>
    </FormLayout>
  );
}

function ImageBox({
                    imgURL,
                    onFileChange,
                  }: {
  imgURL: string;
  onFileChange: any;
}) {
  return (
    <div class="mt-5 flex h-full items-center justify-center">
      <div class="relative flex size-56 flex-col items-center justify-center gap-2">
        <img
          src={imgURL}
          alt="Uploaded image preview"
          class="size-56 rounded-xl object-cover"
        />
        <input
          class="absolute size-full cursor-pointer rounded-xl opacity-0"
          type="file"
          accept="image/*"
          onChange={onFileChange}
        />
      </div>
    </div>
  );
}

function UploadBox({ onFileChange }: { onFileChange: any }) {
  return (
    <>
      <div class="relative flex size-56 flex-col items-center justify-center rounded-xl bg-main">
        <input
          class="absolute size-full opacity-0"
          type="file"
          accept="image/*"
          onChange={onFileChange}
        />
        <span class="material-symbols-rounded text-secondary pointer-events-none z-10 text-[45px]">
          camera_alt
        </span>
      </div>
    </>
  );
}
