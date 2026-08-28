<script lang="ts">
  import "./layout.css";
  import favicon from "$lib/assets/favicon.svg";

  let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<nav>
  <h1>box</h1>
  <a role="button" href="/mailboxes">Manage Mailboxes</a>
</nav>
{@render children()}
