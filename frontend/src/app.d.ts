declare namespace svelteHTML {
  interface HTMLAttributes<T> {
    'on:shortcut'?: (event: CustomEvent<any>) => void;
  }
}