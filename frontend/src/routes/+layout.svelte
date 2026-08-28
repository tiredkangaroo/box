<script lang="ts">
  import "./layout.css";
  import type { Mailbox } from "$lib/types";
  import favicon from "$lib/assets/favicon.svg";

  let { children } = $props();
  let mailboxes = $state(undefined);

  async function getMailboxes() {
    const resp = await fetch("/api/mailboxes");
    return await resp.json();
  }
  getMailboxes().then((data) => {
    mailboxes = data;
  });
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<nav>
  <div>
    <h1 class="text-center">box</h1>
    <h2>Mailboxes</h2>
    {#if mailboxes != undefined}
      <div id="nav-mailboxes-list">
        {#each mailboxes as Array<Mailbox> as mailbox}
          <button
            class="nav-mailbox"
            onclick={(e) => {
              (e.target as HTMLButtonElement).classList.toggle("active");
            }}>{mailbox.username}</button
          >
        {/each}
        <button
          class="nav-mailbox"
          onclick={(e) => {
            (e.target as HTMLButtonElement).classList.toggle("active");
          }}>meow</button
        >
      </div>
    {/if}
  </div>
  <a role="button" href="/mailboxes" id="manage-mailboxes-button">Manage Mailboxes</a>
</nav>
{@render children()}
