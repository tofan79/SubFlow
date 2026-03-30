<script lang="ts">
  import { Button, Card } from 'flowbite-svelte';
</script>

<main class="flex min-h-screen items-center justify-center p-8">
  <Card class="w-full max-w-md bg-[#1e1e24] border-gray-800 shadow-xl">
    <h1 class="mb-4 text-3xl font-bold tracking-tight text-white">
      Welcome to SubFlow
    </h1>
    <p class="mb-6 font-normal text-gray-400">
      Built with Wails v2, SvelteKit, Tailwind CSS, Flowbite Svelte, and the Geist font. Dark mode is configured by default.
    </p>
    <Button color="primary" class="w-full bg-primary-600 hover:bg-primary-700">Get Started</Button>
  </Card>
</main>
