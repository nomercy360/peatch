import { createSignal, Match, Switch } from 'solid-js';
import { FormLayout } from '../../pages/users/edit';

export default function ImageUpload(props: {
  imgURL: string;
  imgFile: File | null;
  setImgFile: (file: File) => void;
}) {
  const [previewUrl, setPreviewUrl] = createSignal(props.imgURL);

  const handleFileChange = (event: any) => {
    const file = event.target.files[0];
    if (file) {
      props.setImgFile(file);
      setPreviewUrl('');

      const reader = new FileReader();
      reader.onload = e => {
        setPreviewUrl(e.target?.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const generateRandomAvatar = () => {
    const url = 'https://source.boringavatars.com/beam';

    const resp = fetch(url);

    resp.then(response => {
      response.blob().then(blob => {
        const file = new File([blob], 'avatar.svg', {
          type: 'image/svg+xml',
        });
        props.setImgFile(file);
        setPreviewUrl(url);
      });
    });
  };

  return (
    <FormLayout
      title="Upload your photo"
      description="Select one with good lighting and minimal background details"
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
          <button class="h-10 text-peatch-blue" onClick={generateRandomAvatar}>
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
      <div class="relative flex size-56 flex-col items-center justify-center rounded-xl bg-peatch-bg">
        <input
          class="absolute size-full opacity-0"
          type="file"
          accept="image/*"
          onChange={onFileChange}
        />
        <span class="material-symbols-rounded pointer-events-none z-10 text-[45px] text-peatch-light-black">
          camera_alt
        </span>
      </div>
    </>
  );
}
