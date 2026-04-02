export function shortcuts(node: HTMLElement) {
  function handleKeydown(event: KeyboardEvent) {
    if (
      event.target instanceof HTMLInputElement ||
      event.target instanceof HTMLTextAreaElement ||
      event.target instanceof HTMLSelectElement
    ) {
      if (!event.ctrlKey && !event.metaKey && !event.altKey) {
        return;
      }
    }

    node.dispatchEvent(
      new CustomEvent('shortcut', {
        detail: {
          key: event.key,
          code: event.code,
          ctrl: event.ctrlKey,
          shift: event.shiftKey,
          alt: event.altKey,
          meta: event.metaKey,
          originalEvent: event
        }
      })
    );
  }

  window.addEventListener('keydown', handleKeydown);

  return {
    destroy() {
      window.removeEventListener('keydown', handleKeydown);
    }
  };
}