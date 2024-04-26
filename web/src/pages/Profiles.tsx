import { useButtons } from '../hooks/useBackButton';
import { createEffect, onCleanup } from 'solid-js';

export default function Profiles() {
  const { setBackVisibility, onBackClick, offBackClick } = useButtons()

  const back = () => {
    history.back();
  };

  createEffect(() => {
    setBackVisibility(true);
    onBackClick(back);
  });

  onCleanup(() => {
    setBackVisibility(false);
    offBackClick(back);
  });

  return (
    <div>
      <button onClick={back}>Back</button>
      <h1>Profiles</h1>
    </div>
  );
}