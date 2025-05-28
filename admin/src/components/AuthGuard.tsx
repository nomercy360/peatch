import { createEffect, ParentProps, Show } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { getToken } from "~/lib/api";

export function AuthGuard(props: ParentProps) {
  const navigate = useNavigate();

  createEffect(() => {
    const token = getToken();
    if (!token) {
      navigate("/login");
    }
  });

  return (
    <Show when={getToken()} fallback={null}>
      {props.children}
    </Show>
  );
}